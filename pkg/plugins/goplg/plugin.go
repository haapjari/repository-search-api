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

	// g.calculateSizeOfPrimaryRepositories()

	g.enrichWithLibraryData("")
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

	var (
		output string
		err    string
	)

	// Calculate the amount of code lines for each repository.
	for i := 0; i < len; i++ {
		// If the OriginalCodebaseSize variable is empty, analyze the repository.
		// Otherwise skip the repository, in order to avoid double analysis.
		if repositories.RepositoryData[i].OriginalCodebaseSize == "" {
			// Save the repository url and name to variables, for easier access and less repetition.
			url := repositories.RepositoryData[i].RepositoryUrl
			name := repositories.RepositoryData[i].RepositoryName

			// Clone the repository.
			output, err = runCommand("git", "clone", url)
			fmt.Println(output)
			if err != "" {
				fmt.Println(err)
			}

			// Run "gocloc" and calculate the amount of lines.
			lines := runGoCloc(name)

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

			// Delete the repository.
			output, err = runCommand("rm", "-rf", name)
			fmt.Println(output)
			if err != "" {
				fmt.Println(err)
			}
		}
	}
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
func (g *GoPlugin) enrichWithLibraryData(url string) {
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

	// Flag variable, which keeps track wether the project has inner go.mod files.
	isThereInnerModFiles := false

	// replace -keyword means, that the project contains inner modfiles, boolean flag will be turned to true.
	if strings.Count(outerModFile.String(), "replace") > 0 {
		isThereInnerModFiles = true
	}

	// Parsing of the inner modfiles.

	var (
		innerModFilesString string
		innerModFiles       []string
	)

	// If the project has inner modfiles, parse the locations to the variable innerModFilesSlice.
	if isThereInnerModFiles {
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

	var (
		libraries []string
		// totalCodeLines           int
	)

	// Libraries are saved in a slice of "libraries"
	libraries = parseLibrariesFromModFile(outerModFile.String())

	// Printing
	for _, lib := range libraries {
		fmt.Println(lib)
	}

	// TODO: This will happen in a loop.
	// Perform a GET request, to get the content of the inner modfile.

	// We can discard the replace contents of the inner modfiles, because they already
	// taken in account.

	modFileContent := performGetRequest(innerModFiles[0])

	modFileContentSlice := parseLibrariesFromModFile(modFileContent)

	fmt.Println("---")

	// TODO: For some reason the parseLibrariesFromModFile is not working properly.
	// It's not taking in account the second parenthesis in modFileContent, why not (?)

	// Printing
	for _, lib := range modFileContentSlice {
		fmt.Println(lib)
	}

	fmt.Println("---")

	fmt.Println(modFileContent)

	// libraries = append(libraries, parseLibrariesFromModFile(modFileContent)...)

	// TODO: Parse the library names of the inner modfiles.
	// TODO: Create a cache, that holds the code lines of analyzed libraries.
	// TODO: Implement functionality to calculate code lines of the libraries.
	// TODO: Delete duplicates from the libraries slice.
}

// Performs a GET request to the specified URL.
func performGetRequest(url string) string {
	// Make a GET request to the specified URL
	resp, err := http.Get(url)
	utils.CheckErr(err)

	defer resp.Body.Close()

	// Read the response body into a variable
	body, err := ioutil.ReadAll(resp.Body)
	utils.CheckErr(err)

	return (string(body))
}

// Parses the library names from the go.mod file and save them to the libraries slice.
func parseLibrariesFromModFile(modFile string) []string {
	var (
		outerModFileString       string
		outerModFileSlice        []string
		outerModFileParsingSlice []string
		libraries                []string
	)

	// Count how many times the substring "require" occurs in the outer go.mod file, this
	// gives me the number of how many times the file needs to be splitted in order to be
	// correctly parsed.
	requireCount := strings.Count(modFile, "require")

	// Splitting the outer go.mod file to slice, using "require" as a separator.
	outerModFileSlice = strings.Split(modFile, "require")

	// Now we have count of "requireCount" amount of strings saved in a slice.
	// We can discard the 0th element, because it doesn't contain the library
	// metadata.
	for i := 1; i <= requireCount; i++ {
		// Taking the right side of the splitted string, and
		// re-splitting from the first "(" rune.
		outerModFileParsingSlice = strings.Split(outerModFileSlice[i], "(")

		// Splitting from the first ")" rune.
		outerModFileParsingSlice = strings.Split(outerModFileParsingSlice[i], ")")

		// Saving the contents inside the parenthesis of the outer go.mod file
		// to a separate string variable.
		outerModFileString = outerModFileParsingSlice[0]

		// After splitting, there are trailing whitespace, so trimming that out.
		outerModFileString = strings.TrimSpace(outerModFileString)

		// Split the string from newline characters.
		outerModFileParsingSlice = strings.Split(outerModFileString, "\n")

		// Remove trailing whitespaces from the beginning and end of the elements in the slice.
		for i, line := range outerModFileParsingSlice {
			outerModFileParsingSlice[i] = strings.TrimSpace(line)
		}

		libraries = append(libraries, outerModFileParsingSlice...)
	}

	// if the element contains substring "=>", its inside a replace parenthesis, which means
	// its not a library, but a location for inner go.mod file, and it will be discarded.
	// Create a new slice to hold the elements we want to keep
	filteredLibraries := make([]string, 0)

	// Iterate over a copy of the slice
	for _, lib := range libraries {
		if !strings.Contains(lib, "=>") {
			filteredLibraries = append(filteredLibraries, lib)
		}
	}

	// Update the original slice with the filtered slice.
	libraries = filteredLibraries

	return libraries
}
