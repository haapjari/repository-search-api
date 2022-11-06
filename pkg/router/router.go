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

	r.GET("/api/glass/v1/commit", commit.GetCommits)
	r.POST("/api/glass/v1/commit", commit.CreateCommit)
	r.GET("/api/glass/v1/commit/:id", commit.GetCommitById)
	r.DELETE("/api/glass/v1/commit/:id", commit.DeleteCommitById)
	r.PATCH("/api/glass/v1/commit/:id", commit.UpdateCommitById)

	r.GET("/api/glass/v1/repository", repository.GetRepositories)
	r.POST("/api/glass/v1/repository", repository.CreateRepository)
	r.GET("/api/glass/v1/repository/:id", repository.GetRepositoryById)
	r.DELETE("/api/glass/v1/repository/:id", repository.DeleteRepositoryById)
	r.PATCH("/api/glass/v1/repository/:id", repository.UpdateRepositoryById)

	r.GET("/api/glass/v1/repository/fetch", repository.FetchRepositories)

	r.GET("/api/glass/v1/metrics", prom.Handler)

	r.Run()
}
