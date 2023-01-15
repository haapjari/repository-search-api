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

func GetRepositoryMetadata(c *gin.Context) {
	h := NewHandler(c)

	if c.Query("type") == "go" {
		h.HandleGetRepositoryMetadata()
	}

	if !(c.Query("type") == "go") {
		c.JSON(http.StatusOK, gin.H{"data": "not supported"})
	}
}

func GenerateCsv(c *gin.Context) {
	h := NewHandler(c)
	h.HandleGenerateCsv()
}
