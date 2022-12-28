package goplg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/haapjari/glass/pkg/models"
	"github.com/haapjari/glass/pkg/utils"
	"github.com/hhatto/gocloc"
	JSONParser "github.com/tidwall/gjson"
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

	g.Parser = NewParser()

	return g
}

func (g *GoPlugin) GetRepositoryMetadata(count int) {
	// g.fetchRepositories(count)
	// g.deleteDuplicateRepositories // TODO
	// g.enrichRepositoriesWithPrimaryData()

	g.calculateSizeOfPrimaryRepositories()

	// g.enrichRepositoriesWithLibraryData("")
}

// Enriches the metadata with "Original Codebase Size" variables.
func (g *GoPlugin) calculateSizeOfPrimaryRepositories() {
	// fetch all the repositories from the database.
	repositories := g.getAllRepositories()

	// calculate the amount of repositories, and save it to variable.
	len := len(repositories.RepositoryData)

	// append the https:// and .git prefix and postfix the RepositoryUrl variables.
	for i := 0; i < len; i++ {
		repositories.RepositoryData[i].RepositoryUrl = "https://" + repositories.RepositoryData[i].RepositoryUrl + ".git"
	}

	// TODO: Run this in a loop.

	var (
		outputString string
		errorString  string
	)

	outputString, errorString = runCommand("git", "clone", repositories.RepositoryData[0].RepositoryUrl)

	fmt.Println(outputString)
	fmt.Println(errorString)

	lines := runGoCloc(repositories.RepositoryData[0].RepositoryName)

	fmt.Println(lines)

	// TODO: Update the database with the new data.

	outputString, errorString = runCommand("rm", "-rf", repositories.RepositoryData[0].RepositoryName)

	fmt.Println(outputString)
	fmt.Println(errorString)
}

// Wrapper for "exec/os" command execution.
// Copied from blog: https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
func runCommand(name string, arg ...string) (string, string) {
	cmd := exec.Command(name, arg...)

	var stdout, stderr []byte
	var errStdout, errStderr error

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	err := cmd.Start()
	utils.CheckErr(err)

	// WaitGroup ensures, cmd.Wait() is called, after we finish reading from stdin and stdout.
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
		wg.Done()
	}()

	stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}

	outStr, errStr := string(stdout), string(stderr)

	return outStr, errStr
}

// Helper function for running commands with "os/exec".
// Copied from blog: https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

// Fetches initial metadata of the repositories. Crafts a SourceGraph GraphQL request, and
// parses the repository location to the database table.
func (g *GoPlugin) fetchRepositories(count int) {
	// TODO: Export this as a environment or configuration value.
	queryStr := "{search(query:\"lang:go +  AND select:repo AND repohasfile:go.mod AND count:" + strconv.Itoa(count) + "\", version:V2){results{repositories{name}}}}"

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

// Calculates the lines of code using https://github.com/hhatto/gocloc
// in the path provided and return the value.
func runGoCloc(path string) int {
	languages := gocloc.NewDefinedLanguages()
	options := gocloc.NewClocOptions()

	paths := []string{
		path,
	}

	processor := gocloc.NewProcessor(languages, options)

	result, err := processor.Analyze(paths)
	utils.CheckErr(err)

	return int(result.Total.Code)
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

			queryStr := "{repository(owner: \"" + owner + "\", name: \"" + name + "\") {defaultBranchRef {target {... on Commit {history {totalCount}}}}openIssues: issues(states:OPEN) {totalCount}closedIssues: issues(states:CLOSED) {totalCount}languages {totalSize}stargazerCount licenseInfo {key}createdAt latestRelease{publishedAt} primaryLanguage{name}}}"

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

// TODO
func (g *GoPlugin) enrichWithLibraryData(repositoryUrl string) {
	// Query String
	// TODO: Replace URL from queryString with repositoryUrl
	queryString := "{repository(name: \"github.com/kubernetes/kubernetes\") {defaultBranch {target {commit {blob(path: \"go.mod\") {content}}}}}}"

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
	outerModFile := JSONParser.Get(string(sourceGraphResponseBody), "data.repository.defaultBranch.target.commit.blob.content")

	// TODO: Using flags is not the best way to do this, rethink this.
	// Flag Variable, which keeps track wether the project has inner go.mod files.
	innerModfiles := false

	// replace -keyword means, that the project contains inner modfiles, boolean flag will be turned to true.
	if strings.Count(outerModFile.String(), "replace") > 0 {
		innerModfiles = true
	}

	// Parsing of the inner modfiles.

	var (
		innerModFilesString string
		innerModFiles       []string
	)

	// If the project has inner modfiles, parse the locations to the variable innerModFilesSlice.
	if innerModfiles {
		innerModFiles = strings.Split(outerModFile.String(), "replace")
		innerModFilesString = innerModFiles[1]

		// remove trailing '(' and ')' runes from the string
		innerModFilesString = utils.RemoveCharFromString(innerModFilesString, '(')
		innerModFilesString = utils.RemoveCharFromString(innerModFilesString, ')')

		// remove trailing whitespace,
		innerModFilesString = strings.TrimSpace(innerModFilesString)

		// split the strings from newline characters
		innerModFiles = strings.Split(innerModFilesString, "\n")

		// remove trailing whitespaces from the beginning and end of the elements in the string
		for i, line := range innerModFiles {
			innerModFiles[i] = strings.TrimSpace(line)
		}

		// remove characters until => occurence.
		for i, line := range innerModFiles {
			str := strings.Index(line, "=>")
			innerModFiles[i] = line[str+2:]
		}

		// TODO: Refactor prefix and postfix, with actual repository names.
		// add the "https://raw.githubusercontent.com/kubernetes/kubernetes/master" prefix and
		// "/go.mod" postfix to the string.
		prefix := "https://raw.githubusercontent.com/kubernetes/kubernetes/master"
		postfix := "go.mod"

		for i, line := range innerModFiles {
			// removing the dot from the beginning of the file
			line = strings.Replace(line, ".", "", 1)

			// remove whitespace inbetween the string
			line = strings.Replace(line, " ", "", -1)

			// concat the prefix and postfix
			line = prefix + line + "/" + postfix

			// modify the original slice
			innerModFiles[i] = line
		}
	}

	// TODO: Create variable, that holds the amount of code lines of a repository.
	// TODO: Create a cache, that holds the code lines of analyzed libraries.
	// TODO: Parse the library names of the outer modfile.
	// TODO: Parse the library names of the inner modfiles.
	// TODO: Implement functionality to calculate code lines of the libraries.

	fmt.Println(outerModFile.String())
}
