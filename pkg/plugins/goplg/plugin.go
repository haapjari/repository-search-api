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

	// TODO: These doesn't have "support for persistent analysis"
	// g.pruneDuplicates()
	// g.pruneLanguages()
	// g.processRepositories()
	// g.processLibraries()
}

// ************************************************************* //
// ************************************************************* //
// ************************************************************* //

// Fetches initial metadata of the repositories. Crafts a SourceGraph GraphQL request, and
// parses the repository location to the database table.
func (g *GoPlugin) fetchRepositories(count int) {
	log.Println("Fetching Repositories.")

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

	// *** //

	length := len(response.Data.Search.Results.Repositories)
	repos := response.Data.Search.Results.Repositories

	// var wg sync.WaitGroup
	// s := make(chan int, g.MaxRoutines)

	// Write the repositories to the database, avoid writing duplicates.
	for i := 0; i < length; i++ {
		//		s <- 1
		//wg.Add(1)
		//go func(i int) {
		r := models.Repository{RepositoryName: repos[i].Name, RepositoryUrl: repos[i].Name, OpenIssueCount: "", ClosedIssueCount: "", OriginalCodebaseSize: "", LibraryCodebaseSize: "", RepositoryType: "", PrimaryLanguage: ""}
		if !g.checkIfRepositoryExists(r) {
			g.DatabaseClient.Create(&r)
		}
		//defer func() { <-s }()
		//	}(i)
		//wg.Done()
	}

	//	wg.Wait()

	//	for !(len(s) == 0) {
	//
	// continue
	// }

	// *** //

	databaseSnapshot := g.getAllRepositories()

	// Reads the repositories -tables values to memory, crafts a GitHub GraphQL requests of the
	// repositories, and appends the database entries with Open Issue Count, Closed Issue Count,
	// Commit Count, Original Codebase Size, Repository Type, Primary Language, Stargazers Count,
	// Creation Date, License.
	for i := 0; i < len(databaseSnapshot); i++ {
		if !g.hasBeenEnriched(databaseSnapshot[i]) {

			fmt.Println(databaseSnapshot[i].RepositoryUrl)
			fmt.Println("Repo has not been enriched.")

			repositoryOwner, repositoryName := g.Parser.ParseRepository(databaseSnapshot[i].RepositoryUrl)

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
			}`, repositoryOwner, repositoryName)

			rawRequestBody := map[string]string{
				"query": queryStr,
			}

			requestBody, err := json.Marshal(rawRequestBody)
			utils.CheckErr(err)

			b := bytes.NewBuffer(requestBody)

			githubRequest, err := http.NewRequest("POST", GITHUB_GRAPHQL_API_BASEURL, b)
			utils.CheckErr(err)

			githubRequest.Header.Set("Accept", "application/vnd.github.v3+json")

			githubResponse, err := g.GitHubClient.Do(githubRequest)
			utils.CheckErr(err)

			defer githubResponse.Body.Close()

			githubResponseBody, err := ioutil.ReadAll(githubResponse.Body)
			utils.CheckErr(err)

			var githubResponseStruct models.GitHubResponseStruct
			json.Unmarshal([]byte(githubResponseBody), &githubResponseStruct)

			// Start Updating Values.

			var existingResponseStruct models.Repository

			if err := g.DatabaseClient.Where("repository_name = ?", databaseSnapshot[i].RepositoryName).First(&existingResponseStruct).Error; err != nil {
				utils.CheckErr(err)
			}

			existingResponseStruct.RepositoryName = repositoryName
			existingResponseStruct.RepositoryUrl = databaseSnapshot[i].RepositoryUrl
			existingResponseStruct.OpenIssueCount = strconv.Itoa(githubResponseStruct.Data.Repository.OpenIssues.TotalCount)
			existingResponseStruct.ClosedIssueCount = strconv.Itoa(githubResponseStruct.Data.Repository.ClosedIssues.TotalCount)
			existingResponseStruct.CommitCount = strconv.Itoa(githubResponseStruct.Data.Repository.DefaultBranchRef.Target.History.TotalCount)
			existingResponseStruct.RepositoryType = "primary"
			existingResponseStruct.PrimaryLanguage = githubResponseStruct.Data.Repository.PrimaryLanguage.Name
			existingResponseStruct.CreationDate = githubResponseStruct.Data.Repository.CreatedAt
			existingResponseStruct.StargazerCount = strconv.Itoa(githubResponseStruct.Data.Repository.StargazerCount)
			existingResponseStruct.LicenseInfo = githubResponseStruct.Data.Repository.LicenseInfo.Key
			existingResponseStruct.LatestRelease = githubResponseStruct.Data.Repository.LatestRelease.PublishedAt

			g.DatabaseClient.Model(&existingResponseStruct).Updates(existingResponseStruct)
		}
	}
}

// Process repositories, fetch metadata and calculate how many lines of code there are
// in the repository.
func (g *GoPlugin) processRepositories() {
	r := g.getAllRepositories()

	// var wg sync.WaitGroup
	// s := make(chan int, g.MaxRoutines)

	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		if err := os.Mkdir("tmp", 0777); err != nil {
			utils.CheckErr(err)
		}
	}

	// Append the https:// and .git prefix and postfix the RepositoryUrl variables.
	for i := 0; i < len(r); i++ {
		name := r[i].RepositoryName

		// Fetching a real time version of the repository, from the database, and checking against fields
		// of that, instead of checking against fields of copy, which have been modified in this function.
		var databaseVersion models.Repository
		if err := g.DatabaseClient.Where("repository_name = ?", name).First(&databaseVersion).Error; err != nil {
			log.Print("repository " + r[i].RepositoryName + " doesn't exist in the database")
		}

		// Check, wether this has been already analyzed, and skip unnecessary analysis.
		if databaseVersion.OriginalCodebaseSize == "" {
			log.Print("Processing Repository: " + r[i].RepositoryName)

			r[i].RepositoryUrl = "https://" + r[i].RepositoryUrl + ".git"

			err := utils.Command("git", "clone", "--depth", "1", r[i].RepositoryUrl, "tmp"+"/"+r[i].RepositoryName)
			if err != nil {
				fmt.Printf("Error while cloning repository %s: %s, skipping...\n", r[i].RepositoryUrl, err)
			}

			lin, err := g.calculateCodeLines("tmp" + "/" + r[i].RepositoryName)
			if err != nil {
				fmt.Print(err.Error())
			}

			g.updateCoreSize(r[i].RepositoryName, lin)
			g.pruneTemporaryFolder()
		}
	}
}

// Loop through repositories, generate the dependency map from the go.mod files of the
// repositories, download the dependencies to the local disk, calculate their sizes and
// save the values to the database.
func (g *GoPlugin) processLibraries() {
	log.Println("Processing Libraries.")

	r := g.getAllRepositories()
	libs := g.generateDependenciesMap(r)
	var wg sync.WaitGroup
	var m sync.RWMutex
	s := make(chan int, g.MaxRoutines)

	os.Setenv("GOPATH", utils.GetProcessDirPath())

	utils.CopyFile("go.mod", "go.mod.bak")
	utils.CopyFile("go.sum", "go.sum.bak")

	// Loop through repositories.
	for i := 0; i < len(r); i++ {
		name := r[i].RepositoryName
		l := 0
		s <- 1

		// Loop through the libraries, which are saved to the map, where dependencies
		// are accessible by repository name. Download them to the local disk, calculate
		// their sizes and append to the 'l' -variable.
		for j := 0; j < len(libs[name]); j++ {

			wg.Add(1)
			go func(j int) {
				// When the concurrently running goroutines is zero, we'll close the sem
				// channel, in order to avoid a memory leak, because sem is initialized
				// for every repository.
				defer func() {
					wg.Done()
					<-s
				}()

				m.RLock()
				value, ok := g.LibraryCache[libs[name][j]]
				m.RUnlock()
				if ok {
					m.Lock()
					l += value
					m.Unlock()
					return
				} else {
					// This is not the most elegant way, and it's using more computation without mutexes.
					// multiple goroutines crash, and try to download same library, and this leads to errors
					// but now the program has functionality to recover from that. This increases performance
					// alot, but it's using waste computation. This might just be a risk, that has to be accepted.

					// m.Lock()
					err := utils.Command("go", "get", "-d", "-v", convertToDownloadableFormat(libs[name][j]))
					if err != nil {
						fmt.Printf("Error while processing library %s: %s, skipping...\n", libs[name][j], err)
					}
					// m.Unlock()

					lin, err := g.calculateCodeLines(utils.GetProcessDirPath() + "/" + "pkg/mod" + "/" + parseLibraryUrl(libs[name][j]))
					if err != nil {
						fmt.Println("Error, while calculating code lines:", err.Error())
					}

					m.Lock()
					g.LibraryCache[libs[name][j]] = lin
					l += lin
					m.Unlock()
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

	// When the Channel Length is not 0, there is still running Threads.
	for !(len(s) == 0) {
		continue
	}

	close(s)

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

func (g *GoPlugin) pruneLanguages() {
	r := g.getAllRepositories()
	for i := 0; i < len(r); i++ {

		if r[i].PrimaryLanguage != "Go" {
			repo := r[i]
			g.DatabaseClient.Delete(repo)
		}

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
