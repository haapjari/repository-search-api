package goplg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	GitHubClient   *http.Client
	Parser         *Parser
	DatabaseClient *gorm.DB
}

func NewGoPlugin(DatabaseClient *gorm.DB) *GoPlugin {
	g := new(GoPlugin)

	g.HttpClient = &http.Client{}
	g.HttpClient.Timeout = time.Minute * 10 // TODO: Is this enough (?)

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_API_TOKEN},
	)

	g.GitHubClient = oauth2.NewClient(context.Background(), tokenSource)
	g.DatabaseClient = DatabaseClient

	g.Parser = NewParser()

	return g
}

func (g *GoPlugin) GetRepositoryMetadata(count int) {
	// g.getBaseMetadataFromSourceGraph(count)
	g.enrichMetadataWithGitHubApi() // TODO
}

// PRIVATE METHODS

func (g *GoPlugin) getBaseMetadataFromSourceGraph(count int) {
	queryStr := "{search(query:\"lang:go +  AND select:repo AND repohasfile:go.mod AND count:" + strconv.Itoa(count) + "\", version:V2){results{repositories{name}}}}"

	rawReqBody := map[string]string{
		"query": queryStr,
	}

	jsonReqBody, err := json.Marshal(rawReqBody)
	utils.LogErr(err)

	bytesReqBody := bytes.NewBuffer(jsonReqBody)

	request, err := http.NewRequest("POST", SOURCEGRAPH_GRAPHQL_API_BASEURL, bytesReqBody)
	request.Header.Set("Content-Type", "application/json")
	utils.LogErr(err)

	res, err := g.HttpClient.Do(request)
	utils.LogErr(err)

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	utils.LogErr(err)

	resMap, err := g.Parser.ParseSourceGraphResponse(string(resBody))
	utils.LogErr(err)

	if resMap != nil {
		for i := 0; i < count; i++ {
			if resMap["repositories"].([]interface{})[i].(map[string]interface{}) != nil {
				for _, value := range resMap["repositories"].([]interface{})[i].(map[string]interface{}) {
					r := models.Repository{RepositoryName: fmt.Sprintf("%v", value), RepositoryUrl: fmt.Sprintf("%v", value), OpenIssueCount: "", ClosedIssueCount: "", OriginalCodebaseSize: "", LibraryCodebaseSize: "", RepositoryType: "", PrimaryLanguage: ""}
					g.DatabaseClient.Create(&r)
				}
			}
		}
	}
}

func (g *GoPlugin) enrichMetadataWithGitHubApi() {
	// Read all the Repositories to Memory
	// Parse URL of the first Repository
	// Parse Owner and Name of the Repository

	r := g.getAllRepositories()
	c := len(r.RepositoryData)

	var wg sync.WaitGroup

	// TODO: is this goroutine necessary (?)
	for i := 0; i < c; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			owner := g.Parser.ParseRepositoryOwner(r.RepositoryData[i].RepositoryUrl)
			name := g.Parser.ParseRepositoryName(r.RepositoryData[i].RepositoryUrl)

			fmt.Println(owner)
			fmt.Println(name)
		}(i)
	}

	wg.Wait()

	// owner := "sulu"
	// name := "sulu"
	// queryStr := "{repository(owner: \"" + owner + "\", name: \"" + name + "\") {defaultBranchRef {target {... on Commit {history {totalCount}}}}openIssues: issues(states:OPEN) {totalCount}closedIssues: issues(states:CLOSED) {totalCount}languages {totalSize}}}"

	// rawGithubRequestBody := map[string]string{
	// 	"query": queryStr,
	// }

	// jsonGithubRequestBody, err := json.Marshal(rawGithubRequestBody)
	// utils.LogErr(err)

	// bytesReqBody := bytes.NewBuffer(jsonGithubRequestBody)

	// githubRequest, err := http.NewRequest("POST", GITHUB_GRAPHQL_API_BASEURL, bytesReqBody)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// githubRequest.Header.Set("Accept", "application/vnd.github.v3+json")

	// githubResponse, err := g.GitHubClient.Do(githubRequest)
	// utils.LogErr(err)

	// defer githubResponse.Body.Close()

	// githubResponseBody, err := ioutil.ReadAll(githubResponse.Body)
	// utils.LogErr(err)

	// // Parse Response Body to JSON
	// var jsonGithubResponse GitHubResponse
	// json.Unmarshal([]byte(githubResponseBody), &jsonGithubResponse)

	// Save Values to Database Entries
}

func (g *GoPlugin) getAllRepositories() models.RepositoryResponse {
	getRepositoriesRequest, err := http.NewRequest("GET", REPOSITORY_API_BASEURL, nil)
	utils.LogErr(err)

	getRepositoriesRequest.Header.Set("Content-Type", "application/json")

	getRepositoriesResponse, err := g.HttpClient.Do(getRepositoriesRequest)
	utils.LogErr(err)

	defer getRepositoriesResponse.Body.Close()

	getRepositoriesResponseBody, err := ioutil.ReadAll(getRepositoriesResponse.Body)
	utils.LogErr(err)

	var repositories models.RepositoryResponse
	json.Unmarshal([]byte(getRepositoriesResponseBody), &repositories)

	return repositories
}
