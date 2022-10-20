package router

import (
	"github.com/haapjari/glass/pkg/controllers/repository"
	"github.com/haapjari/glass/pkg/database"

	"github.com/gin-gonic/gin"
)

func SetupRouter() {
	r := gin.Default()
	db := database.SetupDatabase()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.GET("/api/glass/v1/npm/analysis/repository", repository.GetRepositories)
	r.POST("/api/glass/v1/npm/analysis/repository", repository.CreateRepository)
	r.GET("/api/glass/v1/npm/analysis/repository/:id", repository.GetRepositoryById)
	r.DELETE("/api/glass/v1/npm/analysis/repository:id", repository.DeleteRepositoryById)
	r.PATCH("/api/glass/v1/npm/analysis/repository:id", repository.UpdateRepositoryById)

	r.Run()
}
