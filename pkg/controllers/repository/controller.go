package repository

import (
	"fmt"
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
	e := h.HandleGetRepositories()

	c.JSON(http.StatusOK, gin.H{"data": e})
}

func GetRepositoryById(c *gin.Context) {
	h := NewHandler(c)
	e := h.HandleGetRepositories()

	c.JSON(http.StatusOK, gin.H{"data": e})
}

func CreateRepository(c *gin.Context) {
	h := NewHandler(c)
	b := h.HandleCreateRepository()

	c.JSON(http.StatusOK, gin.H{"data: ": b})
}

func DeleteRepositoryById(c *gin.Context) {
	h := NewHandler(c)
	b := h.HandleDeleteRepositoryById()

	c.JSON(http.StatusOK, gin.H{"data ": b})
}

func UpdateRepositoryById(c *gin.Context) {
	h := NewHandler(c)
	e := h.HandleUpdateRepositoryById()

	c.JSON(http.StatusOK, gin.H{"data": e})
}

// TODO
func FetchRepositories(c *gin.Context) {
	h := NewHandler(c)

	if c.Query("project_type") == "npm" {

		// Method fetches desired amount of repositories from GitHub, through GitHub API
		// and saves those repositories to "repositories" -table in the local PostgreSQL.
		b := h.HandleFetchRepositories()
		fmt.Println(b)

		// Creates new instance of NpmCounter, which is a plugin, that is capable of calculating
		// the amount of Code Lines in the npm -repository. All the data is saved in the
		// database.
		// c := npm.NewNpmCounter()
		// c.CalculateCodeLines()
	}

	if !(c.Query("project_type") == "npm") {
		c.JSON(http.StatusOK, gin.H{"data": "not supported"})
	}
}
