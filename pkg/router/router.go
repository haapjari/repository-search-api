package router

import (
	"github.com/haapjari/glass/pkg/controllers/commit"
	"github.com/haapjari/glass/pkg/controllers/repository"
	"github.com/haapjari/glass/pkg/database"
	"github.com/haapjari/glass/pkg/metrics/prom"

	"github.com/gin-gonic/gin"
)

func SetupRouter() {
	r := gin.Default()
	prom := prom.NewProm()

	db := database.SetupDatabase()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.GET("/api/gls/v1/commit", commit.GetCommits)
	r.POST("/api/gls/v1/commit", commit.CreateCommit)
	r.GET("/api/gls/v1/commit/:id", commit.GetCommitById)
	r.DELETE("/api/gls/v1/commit/:id", commit.DeleteCommitById)
	r.PATCH("/api/gls/v1/commit/:id", commit.UpdateCommitById)

	r.GET("/api/gls/v1/repository", repository.GetRepositories)
	r.POST("/api/gls/v1/repository", repository.CreateRepository)
	r.GET("/api/gls/v1/repository/:id", repository.GetRepositoryById)
	r.DELETE("/api/gls/v1/repository/:id", repository.DeleteRepositoryById)
	r.PATCH("/api/gls/v1/repository/:id", repository.UpdateRepositoryById)
	r.GET("/api/gls/v1/repository/fetch/metadata", repository.GetRepositoryMetadata)

	r.GET("/api/gls/v1/repository/generate/csv", repository.GenerateCsv)

	r.GET("/api/gls/v1/metrics", prom.Handler)

	r.Run()
}
