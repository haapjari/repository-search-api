package goplg

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/haapjari/glass/pkg/models"
	"github.com/haapjari/glass/pkg/utils"
)

func (g *GoPlugin) writeSourceGraphResponseToDatabase(length int, repositories []SourceGraphRepositoriesStruct) {
	var wg sync.WaitGroup

	// Semaphore is a safeguard to goroutines, to allow only "MaxThreads" run at the same time.
	semaphore := make(chan int, g.MaxThreads)

	for i := 0; i < length; i++ {
		semaphore <- 1
		wg.Add(1)

		go func(i int) {
			r := models.Repository{RepositoryName: repositories[i].Name, RepositoryUrl: repositories[i].Name, OpenIssueCount: "", ClosedIssueCount: "", OriginalCodebaseSize: "", LibraryCodebaseSize: "", RepositoryType: "", PrimaryLanguage: ""}
			g.DatabaseClient.Create(&r)

			defer func() { <-semaphore }()
		}(i)
		wg.Done()
	}
	wg.Wait()

	// When the Channel Length is not 0, there is still running goroutines.
	for !(len(semaphore) == 0) {
		continue
	}
}

// TODO: Refactor, to not to use HTTP requests (?)
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
