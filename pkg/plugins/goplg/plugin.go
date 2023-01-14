package goplg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/haapjari/glass/pkg/models"
	"github.com/haapjari/glass/pkg/utils"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

var (
	GITHUB_API_TOKEN                string = fmt.Sprintf("%v", utils.GetGithubToken())
	GITHUB_USERNAME                 string = fmt.Sprintf("%v", utils.GetGithubUsername())
	SOURCEGRAPH_GRAPHQL_API_BASEURL string = utils.GetSourceGraphGQLApi()
	GITHUB_GRAPHQL_API_BASEURL      string = utils.GetGitHubGQLApi()
	REPOSITORY_API_BASEURL          string = "http://" + utils.GetBaseUrl() + "/" + "repository"
)

type GoPlugin struct {
	GitHubApiToken string
	GitHubUsername string
	HttpClient     *http.Client
	Parser         *Parser
	DatabaseClient *gorm.DB
	GitHubClient   *http.Client
	MaxRoutines    int
	LibraryCache   map[string]int
}

func NewGoPlugin(DatabaseClient *gorm.DB) *GoPlugin {
	g := new(GoPlugin)
	var err error

	g.HttpClient = &http.Client{Timeout: time.Minute * 10}
	g.GitHubClient = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_API_TOKEN},
	))
	g.DatabaseClient = DatabaseClient
	g.LibraryCache = make(map[string]int)
	g.Parser = NewParser()
	g.MaxRoutines, err = strconv.Atoi(utils.GetMaxGoRoutines())
	utils.CheckErr(err)

	return g
}

// Fetch Repositories and Enrich the Repositories with Metadata.
func (g *GoPlugin) GenerateRepositoryData(c int) {
	g.fetchRepositories(c)
	g.processRepositories()
	g.processLibraries()
}

// ************************************************************* //
// ************************************************************* //
// ************************************************************* //

// Fetches initial metadata of the repositories. Crafts a SourceGraph GraphQL request, and
// parses the repository location to the database table.
func (g *GoPlugin) fetchRepositories(count int) {
	queryStr := `{
		search(query: "lang:go + AND select:repo AND repohasfile:go.mod AND count:` + strconv.Itoa(count) + `", version:V2) { results {
				repositories {
					name
				}
			}
		}
	}`

	rawReqBody := map[string]string{
		"query": queryStr,
	}

	jsonReqBody, err := json.Marshal(rawReqBody)
	utils.CheckErr(err)

	bytesReqBody := bytes.NewBuffer(jsonReqBody)

	request, err := http.NewRequest("POST", SOURCEGRAPH_GRAPHQL_API_BASEURL, bytesReqBody)
	request.Header.Set("Content-Type", "application/json")
	utils.CheckErr(err)

	res, err := g.HttpClient.Do(request)
	utils.CheckErr(err)

	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	utils.CheckErr(err)

	var response models.SourceGraphResponseStruct
	json.Unmarshal([]byte(responseBody), &response)

	g.processSourceGraphResponse(len(response.Data.Search.Results.Repositories), response.Data.Search.Results.Repositories)
	g.enrichWithPrimaryRepositoryData()
	g.pruneDuplicates()

	// TODO: Alot of requests seem to result primary language repositories, which arent Go.
	// Those have to be pruned out.
}

// TODO: Doc
func (g *GoPlugin) processRepositories() {
	repositories := g.getAllRepositories()

	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		if err := os.Mkdir("tmp", 0777); err != nil {
			utils.CheckErr(err)
		}
	}

	// append the https:// and .git prefix and postfix the RepositoryUrl variables.
	for i := 0; i < len(repositories); i++ {
		repositories[i].RepositoryUrl = "https://" + repositories[i].RepositoryUrl + ".git"
	}

	for _, repo := range repositories {
		if repo.OriginalCodebaseSize == "" {
			err := utils.Command("git", "clone", "--depth", "1", repo.RepositoryUrl, "tmp"+"/"+repo.RepositoryName)
			if err != nil {
				fmt.Printf("Error while cloning repository %s: %s, skipping...\n", repo.RepositoryUrl, err)
				continue
			}

			lin, err := g.calculateCodeLines("tmp/" + repo.RepositoryName)
			if err != nil {
				fmt.Print(err.Error())
				continue
			}

			g.updateCoreSize(repo.RepositoryName, lin)
		}
		g.pruneTemporaryFolder()
	}
}

// Loop through repositories, generate the dependency map from the go.mod files of the
// repositories, download the dependencies to the local disk, calculate their sizes and
// save the values to the database.
func (g *GoPlugin) processLibraries() {
	r := g.getAllRepositories()
	libs := g.generateDependenciesMap(r)

	os.Setenv("GOPATH", utils.GetProcessDirPath())

	utils.CopyFile("go.mod", "go.mod.bak")
	utils.CopyFile("go.sum", "go.sum.bak")

	// Loop through repositories.
	for i := 0; i < len(r); i++ {
		name := r[i].RepositoryName
		l := 0

		var wg sync.WaitGroup
		var mux sync.RWMutex
		var sem = make(chan struct{}, g.MaxRoutines)
		var semCount int64

		// Loop through the libraries, which are saved to the map, where dependencies
		// are accessible by repository name. Download them to the local disk, calculate
		// their sizes and append to the 'l' -variable.
		for j := 0; j < len(libs[name]); j++ {
			semCount++
			sem <- struct{}{}
			wg.Add(1)
			go func(j int) {
				// When the concurrently running goroutines is zero, we'll close the sem
				// channel, in order to avoid a memory leak, because sem is initialized
				// for every repository.
				defer func() {
					<-sem
					// Using the atomic package here is important, as it's being accessed
					// by multiple goroutines and we need to ensure that the value of semCount
					// is read and written atomically, otherwise, we may end up with a race condition.
					atomic.AddInt64(&semCount, -1)
					if atomic.LoadInt64(&semCount) == 0 {
						close(sem)
					}
					wg.Done()
				}()

				mux.RLock()
				value, ok := g.LibraryCache[libs[name][j]]
				mux.RUnlock()
				if ok {
					mux.Lock()
					l += value
					mux.Unlock()
					return
				} else {
					// This has to be protected with mutex, because "go get" is modifying same "go.mod" -file.
					mux.Lock()
					err := utils.Command("go", "get", "-d", "-v", convertToDownloadableFormat(libs[name][j]))
					if err != nil {
						fmt.Printf("Error while processing library %s: %s, skipping...\n", libs[name][j], err)
						return
					}
					mux.Unlock()

					lin, err := g.calculateCodeLines(utils.GetProcessDirPath() + "/" + "pkg/mod" + "/" + parseLibraryUrl(libs[name][j]))
					if err != nil {
						fmt.Println(err.Error())
						return
					}

					mux.Lock()
					g.LibraryCache[libs[name][j]] = lin
					l += lin
					mux.Unlock()
				}
			}(j)
		}
		wg.Wait()

		g.pruneTemporaryFolder()

		utils.RemoveFiles("go.mod", "go.sum")
		utils.CopyFile("go.mod.bak", "go.mod")
		utils.CopyFile("go.sum.bak", "go.sum")

		var repositoryStruct models.Repository

		if err := g.DatabaseClient.Where("repository_name = ?", name).First(&repositoryStruct).Error; err != nil {
			utils.CheckErr(err)
		}

		repositoryStruct.LibraryCodebaseSize = strconv.Itoa(l)
		g.DatabaseClient.Model(&repositoryStruct).Updates(repositoryStruct)
	}

	os.Setenv("GOPATH", utils.GetDefaultGoPath())

	utils.RemoveFiles("go.mod", "go.sum")
	utils.CopyFile("go.mod.bak", "go.mod")
	utils.CopyFile("go.sum.bak", "go.sum")
	utils.RemoveFiles("go.mod.bak", "go.sum.bak")
}

// ************************************************************* //
// ************************************************************* //
// ************************************************************* //

func (g *GoPlugin) pruneDuplicates() {
	repositories := g.getAllRepositories()
	duplicateRepositories := findDuplicates(repositories)
	amount := len(duplicateRepositories)

	for i := 0; i < amount; i++ {
		var r models.Repository
		name := duplicateRepositories[i].RepositoryName

		if err := g.DatabaseClient.Where("repository_name = ?", name).First(&r).Error; err != nil {
			utils.CheckErr(err)
		}

		g.DatabaseClient.Delete(&r)
	}
}

// Reads the repositories -tables values to memory, crafts a GitHub GraphQL requests of the
// repositories, and appends the database entries with Open Issue Count, Closed Issue Count,
// Commit Count, Original Codebase Size, Repository Type, Primary Language, Stargazers Count,
// Creation Date, License.
func (g *GoPlugin) enrichWithPrimaryRepositoryData() {
	r := g.getAllRepositories()
	c := len(r)
	s := make(chan int, g.MaxRoutines)
	var wg sync.WaitGroup

	for i := 0; i < c; i++ {
		s <- 1
		wg.Add(1)

		go func(i int) {
			owner, name := g.Parser.ParseRepository(r[i].RepositoryUrl)

			queryStr := fmt.Sprintf(`{
					repository(owner: "%s", name: "%s") {
						defaultBranchRef {
							target {
								... on Commit {
								history {
									totalCount
								}
							}
						}
					}	
					openIssues: issues(states:OPEN) {
						totalCount
					}
					closedIssues: issues(states:CLOSED) {
						totalCount
					}
					languages {
						totalSize
					}
					stargazerCount
					licenseInfo {
						key
					}
					createdAt
					latestRelease{
						publishedAt
					}
					primaryLanguage{
						name
					}
				}
			}`, owner, name)

			rawRequestBody := map[string]string{
				"query": queryStr,
			}

			requestBody, err := json.Marshal(rawRequestBody)
			utils.CheckErr(err)

			b := bytes.NewBuffer(requestBody)

			githubRequest, err := http.NewRequest("POST", GITHUB_GRAPHQL_API_BASEURL, b)
			if err != nil {
				log.Fatalln(err)
			}

			githubRequest.Header.Set("Accept", "application/vnd.github.v3+json")

			githubResponse, err := g.GitHubClient.Do(githubRequest)
			utils.CheckErr(err)

			defer githubResponse.Body.Close()

			githubResponseBody, err := ioutil.ReadAll(githubResponse.Body)
			utils.CheckErr(err)

			var newResponseStruct models.GitHubResponseStruct
			json.Unmarshal([]byte(githubResponseBody), &newResponseStruct)

			var existingResponseStruct models.Repository

			if err := g.DatabaseClient.Where("id = ?", r[i].Id).First(&existingResponseStruct).Error; err != nil {
				utils.CheckErr(err)
			}

			var newRepositoryStruct models.Repository

			newRepositoryStruct.RepositoryName = name
			newRepositoryStruct.RepositoryUrl = r[i].RepositoryUrl
			newRepositoryStruct.OpenIssueCount = strconv.Itoa(newResponseStruct.Data.Repository.OpenIssues.TotalCount)
			newRepositoryStruct.ClosedIssueCount = strconv.Itoa(newResponseStruct.Data.Repository.ClosedIssues.TotalCount)
			newRepositoryStruct.CommitCount = strconv.Itoa(newResponseStruct.Data.Repository.DefaultBranchRef.Target.History.TotalCount)
			newRepositoryStruct.RepositoryType = "primary"
			newRepositoryStruct.PrimaryLanguage = newResponseStruct.Data.Repository.PrimaryLanguage.Name
			newRepositoryStruct.CreationDate = newResponseStruct.Data.Repository.CreatedAt
			newRepositoryStruct.StargazerCount = strconv.Itoa(newResponseStruct.Data.Repository.StargazerCount)
			newRepositoryStruct.LicenseInfo = newResponseStruct.Data.Repository.LicenseInfo.Key
			newRepositoryStruct.LatestRelease = newResponseStruct.Data.Repository.LatestRelease.PublishedAt

			g.DatabaseClient.Model(&existingResponseStruct).Updates(newRepositoryStruct)

			defer func() { <-s }()
		}(i)
		wg.Done()
	}

	wg.Wait()

	// When the Channel Length is not 0, there is still running Threads.
	for !(len(s) == 0) {
		continue
	}
}

// Function gets a list of repositories and returns a map of repository names and their dependencies (parsed from go.mod file).
func (g *GoPlugin) generateDependenciesMap(repos []models.Repository) map[string][]string {
	c := len(repos)
	libs := make(map[string][]string)
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, g.MaxRoutines)

	for i := 0; i < c; i++ {
		// Acquire a token from the semaphore
		wg.Add(1)
		sem <- struct{}{}

		// Launch a goroutine
		go func(i int) {
			repoUrl := repos[i].RepositoryUrl
			repoName := repos[i].RepositoryName

			// Query String
			queryString := fmt.Sprintf(`{
			repository(name: "%s") {
				defaultBranch {
					target {
						commit {
							blob(path: "go.mod") {
								content
							}
						}
					}
				}
			}
		}`, repoUrl)

			rawRequestBody := map[string]string{
				"query": queryString,
			}

			jsonRequestBody, err := json.Marshal(rawRequestBody)
			utils.CheckErr(err)

			requestBodyInBytes := bytes.NewBuffer(jsonRequestBody)

			request, err := http.NewRequest("POST", SOURCEGRAPH_GRAPHQL_API_BASEURL, requestBodyInBytes)
			request.Header.Set("Content-Type", "application/json")
			utils.CheckErr(err)

			res, err := g.HttpClient.Do(request)
			utils.CheckErr(err)

			defer res.Body.Close()

			sourceGraphResponseBody, err := ioutil.ReadAll(res.Body)
			utils.CheckErr(err)

			outerModFile := extractDefaultBranchCommitBlobContent(sourceGraphResponseBody)

			var (
				libraries     []string
				innerModFiles []string
			)

			// If the go.mod file has "replace" - keyword, it has inner go.mod files, parse them to a list.
			if g.Parser.CheckForInnerModFiles(outerModFile) {
				// Parse the ending from URL.
				owner, repo, err := parseName(repoUrl)
				utils.CheckErr(err)

				innerModFiles = g.Parser.ParseInnerModFiles(outerModFile, owner+"/"+repo)
			}

			// Protect the libraries slice with a mutex.
			mu.Lock()

			// Parse the name of libraries from modfile to a slice.
			libraries = g.Parser.ParseDependenciesFromModFile(outerModFile)

			// If the go.mod file has "replace" - keyword, it has inner go.mod files,
			// append libraries from inner go.mod files to the libraries slice.
			if g.Parser.CheckForInnerModFiles(outerModFile) {
				// Parse the library names of the inner go.mod files, and append them to the libraries slice.
				for i := 0; i < len(innerModFiles); i++ {
					// Perform a GET request, to get the content of the inner modfile.
					// Append the libraries from the inner modfile to the libraries slice.
					libraries = append(libraries, g.Parser.ParseDependenciesFromModFile(utils.GetRequest(innerModFiles[i]))...)
				}
			}

			// Remove duplicates from the libraries slice.
			libraries = utils.RemoveDuplicates(libraries)

			// Append all the values to the map.
			libs[repoName] = append(libs[repoName], libraries...)

			mu.Unlock()
			<-sem
			wg.Done()
		}(i)
	}

	wg.Wait()

	return libs
}
