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
	GITHUB_API_TOKEN                string = fmt.Sprintf("%v", utils.GetGithubApiToken())
	GITHUB_USERNAME                 string = fmt.Sprintf("%v", utils.GetGithubUsername())
	SOURCEGRAPH_GRAPHQL_API_BASEURL string = utils.GetSourceGraphGraphQlApiBaseurl()
	GITHUB_GRAPHQL_API_BASEURL      string = utils.GetGithubGraphQlApiBaseurl()
	REPOSITORY_API_BASEURL          string = utils.GetRepositoryApiBaseUrl()
)

type GoPlugin struct {
	GitHubApiToken string
	GitHubUsername string
	HttpClient     *http.Client
	Parser         *Parser
	DatabaseClient *gorm.DB
	GitHubClient   *http.Client
	MaxThreads     int
	BatchSize      int
}

func NewGoPlugin(DatabaseClient *gorm.DB) *GoPlugin {
	g := new(GoPlugin)

	g.HttpClient = &http.Client{}
	g.HttpClient.Timeout = time.Minute * 10 // TODO: Environment Variable

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_API_TOKEN},
	)

	g.GitHubClient = oauth2.NewClient(context.Background(), tokenSource)
	g.DatabaseClient = DatabaseClient
	g.MaxThreads = 20
	g.BatchSize = 10

	g.Parser = NewParser()

	return g
}

// Fetch Repositories and Enrich the Repositories with Metadata.
func (g *GoPlugin) GetRepositoryMetadata(c int) {
	g.fetchRepositories(c)
	g.deleteDuplicateRepositories()
	g.enrichWithMetadata()

	// TODO: Alot of requests seem to result primary language repositories, which arent Go.
	// Those have to be pruned out.

	// TODO: Optimizations.
	// There can be goroutine optimizations done in this function.
	g.calcRepoSize()

	// TODO: Optimizations.
	g.calcReposLibSizes()
}

// Delete duplicate repositories.
func (g *GoPlugin) deleteDuplicateRepositories() {
	repositories := g.getAllRepositories()
	duplicateRepositories := findDuplicateRepositoryEntries(repositories.RepositoryData)
	amount := len(duplicateRepositories)

	for i := 0; i < amount; i++ {
		// copy the model, which is going to be deleted
		var r models.Repository

		name := duplicateRepositories[i].RepositoryName

		// Find matching repository from the database.
		if err := g.DatabaseClient.Where("repository_name = ?", name).First(&r).Error; err != nil {
			utils.CheckErr(err)
		}

		// delete from database
		g.DatabaseClient.Delete(&r)
	}
}

// Enriches the metadata with "Original Codebase Size" variables.
func (g *GoPlugin) calcRepoSize() {
	// fetch all the repositories from the database.
	repositories := g.getAllRepositories()

	// Check if the "tmp" directory exists.
	if _, err := os.Stat("tmp"); os.IsNotExist(err) {
		// Create a temporary directory to clone the repositories into.
		if err := os.Mkdir("tmp", 0777); err != nil {
			utils.CheckErr(err)
		}
	}

	// append the https:// and .git prefix and postfix the RepositoryUrl variables.
	for i := 0; i < len(repositories.RepositoryData); i++ {
		repositories.RepositoryData[i].RepositoryUrl = "https://" + repositories.RepositoryData[i].RepositoryUrl + ".git"
	}

	// iterate through all repositories
	for _, repo := range repositories.RepositoryData {
		// If the OriginalCodebaseSize variable is empty, analyze the repository.
		// Otherwise skip the repository, in order to avoid double analysis.
		if repo.OriginalCodebaseSize == "" {
			// Clone the repository into a temporary directory.
			// Attempt to clone "master" branch.
			output, err := runCommand("git", "clone", "--depth", "1", repo.RepositoryUrl, "tmp"+"/"+repo.RepositoryName)
			if err != "" {
				fmt.Println(err)
			}
			fmt.Println(output)

			// Run "gocloc" and calculate the amount of lines.
			lines := runGocloc("tmp/" + repo.RepositoryName)

			// Update the database.
			g.updatePrimaryCodeLinesToDatabase(repo.RepositoryName, lines)

			// Delete the repository.
			output, err = runCommand("rm", "-rf", "tmp"+"/"+repo.RepositoryName)
			fmt.Println(output)
			if err != "" {
				fmt.Println(err)
			}
		}
	}
}

func (g *GoPlugin) updatePrimaryCodeLinesToDatabase(name string, lines int) {
	// Copy the repository struct to a new variable.
	var repositoryStruct models.Repository

	// Find matching repository from the database.
	if err := g.DatabaseClient.Where("repository_name = ?", name).First(&repositoryStruct).Error; err != nil {
		utils.CheckErr(err)
	}

	// Update the OriginalCodebaseSize variable, with calculated value.
	repositoryStruct.OriginalCodebaseSize = strconv.Itoa(lines)

	// Update the database.
	g.DatabaseClient.Model(&repositoryStruct).Updates(repositoryStruct)
}

func (g *GoPlugin) updateLibraryCodeLinesToDatabase(name string, lines int) {
	// Copy the repository struct to a new variable.
	var repositoryStruct models.Repository

	// Find matching repository from the database.
	if err := g.DatabaseClient.Where("repository_name = ?", name).First(&repositoryStruct).Error; err != nil {
		utils.CheckErr(err)
	}

	// Update the OriginalCodebaseSize variable, with calculated value.
	repositoryStruct.LibraryCodebaseSize = strconv.Itoa(lines)

	// Update the database.
	g.DatabaseClient.Model(&repositoryStruct).Updates(repositoryStruct)
}

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

	// Parse Body to JSON
	jsonReqBody, err := json.Marshal(rawReqBody)
	utils.CheckErr(err)

	bytesReqBody := bytes.NewBuffer(jsonReqBody)

	// Craft a request
	request, err := http.NewRequest("POST", SOURCEGRAPH_GRAPHQL_API_BASEURL, bytesReqBody)
	request.Header.Set("Content-Type", "application/json")
	utils.CheckErr(err)

	// Execute request
	res, err := g.HttpClient.Do(request)
	utils.CheckErr(err)

	defer res.Body.Close()

	// Read all bytes from the response
	sourceGraphResponseBody, err := ioutil.ReadAll(res.Body)
	utils.CheckErr(err)

	// Parse bytes JSON.
	var jsonSourceGraphResponse SourceGraphResponse
	json.Unmarshal([]byte(sourceGraphResponseBody), &jsonSourceGraphResponse)

	// Write the response to Database.
	g.writeSourceGraphResponseToDatabase(len(jsonSourceGraphResponse.Data.Search.Results.Repositories), jsonSourceGraphResponse.Data.Search.Results.Repositories)
}

// Reads the repositories -tables values to memory, crafts a GitHub GraphQL requests of the
// repositories, and appends the database entries with Open Issue Count, Closed Issue Count,
// Commit Count, Original Codebase Size, Repository Type, Primary Language, Stargazers Count,
// Creation Date, License.
func (g *GoPlugin) enrichWithMetadata() {
	r := g.getAllRepositories()
	c := len(r.RepositoryData)

	var wg sync.WaitGroup

	// Semaphore is a safeguard to goroutines, to allow only "MaxThreads" run at the same time.
	semaphore := make(chan int, g.MaxThreads)

	for i := 0; i < c; i++ {
		semaphore <- 1
		wg.Add(1)

		go func(i int) {
			// Parse Owner and Name values from the Repository, which are used in the GraphQL query.
			owner, name := g.Parser.ParseRepository(r.RepositoryData[i].RepositoryUrl)

			// Query String
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

			rawGithubRequestBody := map[string]string{
				"query": queryStr,
			}

			// Parse body to JSON.
			jsonGithubRequestBody, err := json.Marshal(rawGithubRequestBody)
			utils.CheckErr(err)

			bytesReqBody := bytes.NewBuffer(jsonGithubRequestBody)

			// Craft a request.
			githubRequest, err := http.NewRequest("POST", GITHUB_GRAPHQL_API_BASEURL, bytesReqBody)
			if err != nil {
				log.Fatalln(err)
			}

			githubRequest.Header.Set("Accept", "application/vnd.github.v3+json")

			// Execute a request with Oauth2 client.
			githubResponse, err := g.GitHubClient.Do(githubRequest)
			utils.CheckErr(err)

			defer githubResponse.Body.Close()

			// Read the response bytes to a variable.
			githubResponseBody, err := ioutil.ReadAll(githubResponse.Body)
			utils.CheckErr(err)

			// Parse bytes to JSON.
			var jsonGithubResponse GitHubResponse
			json.Unmarshal([]byte(githubResponseBody), &jsonGithubResponse)

			var existingRepositoryStruct models.Repository

			// Search for existing model, which matches the id and copy the values to the "existingRepositoryStruct" variable.
			if err := g.DatabaseClient.Where("id = ?", r.RepositoryData[i].Id).First(&existingRepositoryStruct).Error; err != nil {
				utils.CheckErr(err)
			}

			// Create new struct, with updated values.
			var newRepositoryStruct models.Repository

			newRepositoryStruct.RepositoryName = name
			newRepositoryStruct.RepositoryUrl = r.RepositoryData[i].RepositoryUrl
			newRepositoryStruct.OpenIssueCount = strconv.Itoa(jsonGithubResponse.Data.Repository.OpenIssues.TotalCount)
			newRepositoryStruct.ClosedIssueCount = strconv.Itoa(jsonGithubResponse.Data.Repository.ClosedIssues.TotalCount)
			newRepositoryStruct.CommitCount = strconv.Itoa(jsonGithubResponse.Data.Repository.DefaultBranchRef.Target.History.TotalCount)
			newRepositoryStruct.RepositoryType = "primary"
			newRepositoryStruct.PrimaryLanguage = jsonGithubResponse.Data.Repository.PrimaryLanguage.Name
			newRepositoryStruct.CreationDate = jsonGithubResponse.Data.Repository.CreatedAt
			newRepositoryStruct.StargazerCount = strconv.Itoa(jsonGithubResponse.Data.Repository.StargazerCount)
			newRepositoryStruct.LicenseInfo = jsonGithubResponse.Data.Repository.LicenseInfo.Key
			newRepositoryStruct.LatestRelease = jsonGithubResponse.Data.Repository.LatestRelease.PublishedAt

			// Update the existing model, with values from the new struct.
			g.DatabaseClient.Model(&existingRepositoryStruct).Updates(newRepositoryStruct)

			defer func() { <-semaphore }()
		}(i)
		wg.Done()
	}

	wg.Wait()

	// When the Channel Length is not 0, there is still running Threads.
	for !(len(semaphore) == 0) {
		continue
	}
}

// Function gets a list of repositories and returns a map of repository names and their dependencies (parsed from go.mod file).
func (g *GoPlugin) createRepositoryDependenciesMap(repos []models.Repository) map[string][]string {
	repoCount := len(repos)

	// Map of Repository Name (as key) and go.mod -file's dependencies.
	libs := make(map[string][]string)

	var libMutex sync.Mutex

	// ---
	var wg sync.WaitGroup

	semaphore := make(chan struct{}, 20)

	for i := 0; i < repoCount; i++ {
		// Acquire a token from the semaphore
		wg.Add(1)
		semaphore <- struct{}{}

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

			// Construct the Query
			rawRequestBody := map[string]string{
				"query": queryString,
			}

			// Parse Body from Map to JSON
			jsonRequestBody, err := json.Marshal(rawRequestBody)
			utils.CheckErr(err)

			// Convert the Body from JSON to Bytes
			requestBodyInBytes := bytes.NewBuffer(jsonRequestBody)

			// Craft a Request
			request, err := http.NewRequest("POST", SOURCEGRAPH_GRAPHQL_API_BASEURL, requestBodyInBytes)
			request.Header.Set("Content-Type", "application/json")
			utils.CheckErr(err)

			// Execute Request
			res, err := g.HttpClient.Do(request)
			utils.CheckErr(err)

			// Close the Body, after surrounding function returns.
			defer res.Body.Close()

			// Read all bytes from the response. (Empties the res.Body)
			sourceGraphResponseBody, err := ioutil.ReadAll(res.Body)
			utils.CheckErr(err)

			// Parse JSON with "https://github.com/buger/jsonparser"
			outerModFile := extractDefaultBranchCommitBlobContent(sourceGraphResponseBody)

			// Parse the libraries from the go.mod file and inner go.mod files of a project and save them to variables.
			var (
				libraries     []string
				innerModFiles []string
			)

			// If the go.mod file has "replace" - keyword, it has inner go.mod files, parse them to a list.
			if checkInnerModFiles(outerModFile) {
				// Parse the ending from URL.
				owner, repo, err := parseRepositoryName(repoUrl)
				utils.CheckErr(err)

				innerModFiles = parseInnerModFiles(outerModFile, owner+"/"+repo)
			}

			// Protect the libraries slice with a mutex.
			libMutex.Lock()

			// Parse the name of libraries from modfile to a slice.
			libraries = parseLibrariesFromModFile(outerModFile)

			// If the go.mod file has "replace" - keyword, it has inner go.mod files,
			// append libraries from inner go.mod files to the libraries slice.
			if checkInnerModFiles(outerModFile) {
				// Parse the library names of the inner go.mod files, and append them to the libraries slice.
				for i := 0; i < len(innerModFiles); i++ {
					// Perform a GET request, to get the content of the inner modfile.
					// Append the libraries from the inner modfile to the libraries slice.
					libraries = append(libraries, parseLibrariesFromModFile(performGetRequest(innerModFiles[i]))...)
				}
			}

			// Remove duplicates from the libraries slice.
			libraries = removeDuplicates(libraries)

			// Append all the values to the map.
			libs[repoName] = append(libs[repoName], libraries...)

			// Release the mutex.
			libMutex.Unlock()

			// Release the token
			<-semaphore

			// Tell the WaitGroup that we're done
			wg.Done()
		}(i)
	}

	// Wait for all the goroutines to finish
	wg.Wait()

	return libs
}

// Function takes repos and libs and calculates the amount of library code lines for each repository, and writes that to db.
// Requires the libraries to be downloaded in the file system.
func (g *GoPlugin) calculateLibraryCodeLines(repos []models.Repository, libs map[string][]string) {
	repoCount := len(repos)
	var wg sync.WaitGroup

	// ---
	cache := make(map[string]int)
	var cacheLock sync.Mutex

	// TODO
	// Loop through repositories and libraries, and calculate the amount library code lines.
	for i := 0; i < repoCount; i++ {
		repoName := repos[i].RepositoryName
		libCount := len(libs[repoName])
		semaphore := make(chan struct{}, 20)
		totalLibraryCodeLines := 0

		// Lock Cache for Race Conditions.
		cacheLock.Lock()

		if lines, ok := cache[repoName]; ok {
			totalLibraryCodeLines = lines // if the results are in the cache, use them
		} else {
			// if the results are not in the cache, calculate them and add them to the cache
			for j := 0; j < libCount; j++ {
				wg.Add(1)

				go func(i int) {
					semaphore <- struct{}{} // reserve
					libPath := utils.GetTempGoPath() + "/" + "pkg/mod" + "/" + parseGoLibraryUrl(libs[repoName][j])
					lines := runGocloc(libPath)
					totalLibraryCodeLines += lines
					<-semaphore // release

					wg.Done()
				}(i)

				wg.Wait()
			}
			// add the results to the cache
			cache[repoName] = totalLibraryCodeLines
		}
		cacheLock.Unlock()

		g.updateLibraryCodeLinesToDatabase(repoName, totalLibraryCodeLines)
	}
}

// Loop through the repositories, and download the libraries to the local machine.
// TODO: All the repositories are downloaded modified to the same go.mod file - need to address this.
func (g *GoPlugin) downloadGoLibraries(repos []models.Repository, libs map[string][]string) {
	repoCount := len(repos)

	// Reinitialize - allow 20 concurrent goroutines.
	// semaphore := make(chan struct{}, 20)
	// TODO: Enable
	// var wg sync.WaitGroup

	// var goModLock sync.Mutex

	// Read GOPATH variables from the environment.
	tempGoPath := utils.GetTempGoPath()
	goPath := utils.GetGoPath()

	// Change GOPATH to point to temporary directory.
	os.Setenv("GOPATH", tempGoPath)

	// Copy and backup go.mod and go.sum files.
	// This is due to the fact that the go.mod file is modified when downloading libraries,
	// and we don't want to modify the original go.mod file.
	// utils.CopyFile("go.mod", "go.mod.bak")
	// utils.CopyFile("go.sum", "go.sum.bak")

	for i := 0; i < repoCount; i++ {
		repoName := repos[i].RepositoryName
		libCount := len(libs[repoName])

		fmt.Println("Case of repo: ", repoName)
		fmt.Println("We should iterate over: ", libCount, "libraries.")

		// wg.Done()
		// 		}(i)

		// TODO: Implement the analysis in batches.

		// Loop throught the libraries of the repository.
		// TODO: The last batch might not be of size "BatchSize", how to tackle this edge case?
		for z := 0; z < libCount; z++ {
			// Going through the values in batches of "BatchSize".
			if z != 0 && (z+1)%g.BatchSize == 0 {
				// Process the batch of libraries
				for j := z - (g.BatchSize - 1); j <= z; j++ {
					libUrl := parseUrlToDownloadFormat(libs[repoName][j])
					fmt.Println("Index: ", j, "Lib: ", libUrl)
					// TODO: Download and analyze the library
				}
				// TODO: Delete the processed libraries from the disk

				fmt.Println("---")
			}
		}

		// Process the remaining libraries that didn't fit into a full batch
		if libCount%g.BatchSize != 0 {
			for j := libCount - (libCount % g.BatchSize); j < libCount; j++ {
				libUrl := parseUrlToDownloadFormat(libs[repoName][j])
				fmt.Println("Index: ", j, "Lib: ", libUrl)
				// TODO: Download and analyze the library
			}
			// TODO: Delete the processed libraries from the disk
		}
	}

	// TODO: Enable
	// wg.Wait()

	// Change GOPATH to point back to the original directory.
	os.Setenv("GOPATH", goPath)

	// Reset go.mod and go.sum files.
	// 	utils.RemoveFile("go.mod")
	// utils.RemoveFile("go.sum")
	// utils.CopyFile("go.mod.bak", "go.mod")
	// utils.CopyFile("go.sum.bak", "go.sum")
	// utils.RemoveFile("go.mod.bak")
	// utils.RemoveFile("go.sum.bak")

	// Prune the tmp/ folder, if we arent in development mode.
	// 	if !(utils.GetLocalenv() == "development") {
	// os.RemoveAll(utils.GetTempGoPath())
	//	}
}

// TODO
// Enrich the values in the repositories -table with the codebase sizes of the libraries, and append them to the database.
// Before running the gocloc, the vendor means, that the local path is different.
func (g *GoPlugin) calcReposLibSizes() {
	repos := g.getAllRepositories()

	// Map of Repository Name (as key) and go.mod -file's dependencies.
	libs := g.createRepositoryDependenciesMap(repos.RepositoryData)

	// TODO
	g.downloadGoLibraries(repos.RepositoryData, libs)

	// TODO
	// g.calculateLibraryCodeLines(repos.RepositoryData, libs)
}
