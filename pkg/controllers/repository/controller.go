package repository

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RepositoryController struct {
	Handler        *Handler
	Context        *gin.Context
	GitHubApiToken string
	GitHubName     string
}

func GetRepositories(c *gin.Context) {
	h := NewHandler(c)
	h.HandleGetRepositories()
}

func GetRepositoryById(c *gin.Context) {
	h := NewHandler(c)
	h.HandleGetRepositoryById()
}

func CreateRepository(c *gin.Context) {
	h := NewHandler(c)
	h.HandleCreateRepository()
}

func DeleteRepositoryById(c *gin.Context) {
	h := NewHandler(c)
	h.HandleDeleteRepositoryById()
}

func UpdateRepositoryById(c *gin.Context) {
	h := NewHandler(c)
	h.HandleUpdateRepositoryById()
}

// TODO
func FetchRepositories(c *gin.Context) {
	h := NewHandler(c)

	if c.Query("project_type") == "go" {

		// Method fetches desired amount of repositories from GitHub, through GitHub API
		// and saves those repositories to "repositories" -table in the local PostgreSQL.
		h.HandleFetchRepositories()

		// Creates new instance of NpmCounter, which is a plugin, that is capable of calculating
		// the amount of Code Lines in the npm -repository. All the data is saved in the
		// database.
		// c := npm.NewNpmCounter()
		// c.CalculateCodeLines()
	}

	if !(c.Query("project_type") == "go") {
		c.JSON(http.StatusOK, gin.H{"data": "not supported"})
	}
}
