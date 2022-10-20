package repository

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RepositoryController struct {
	Handler        *Handler
	GitHubApiToken string
	GitHubName     string
}

func NewRepositoryController() *RepositoryController {
	return new(RepositoryController)
}

func GetRepositories(c *gin.Context) {
	var (
		h = NewHandler(c)
		e = h.HandleGetRepositories()
	)

	c.JSON(http.StatusOK, gin.H{"data": e})
}

func GetRepositoryById(c *gin.Context) {
	var (
		h = NewHandler(c)
		e = h.HandleGetRepositoryById()
	)

	c.JSON(http.StatusOK, gin.H{"data": e})
}

func CreateRepository(c *gin.Context) {
	var (
		h = NewHandler(c)
		b = h.HandleCreateRepository()
	)

	c.JSON(http.StatusOK, gin.H{"data: ": b})
}

func DeleteRepositoryById(c *gin.Context) {
	var (
		h = NewHandler(c)
		b = h.HandleDeleteRepositoryById()
	)

	c.JSON(http.StatusOK, gin.H{"data ": b})
}

func UpdateRepositoryById(c *gin.Context) {
	var (
		h = NewHandler(c)
		e = h.HandleUpdateRepositoryById()
	)

	c.JSON(http.StatusOK, gin.H{"data": e})
}
