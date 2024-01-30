package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/haapjari/glass/pkg/handlers"
)

type RepositoryController struct {
	Handler        *handlers.Handler
	Context        *gin.Context
	GitHubApiToken string
	GitHubName     string
}

func FetchRepositoryData(c *gin.Context) {
	h := handlers.NewHandler(c)

	if c.Query("type") == "go" {
		h.HandleGetRepositoryMetadata()
	}
}
