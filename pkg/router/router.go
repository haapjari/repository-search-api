package router

import (
	"github.com/haapjari/glass/pkg/controllers/data"
	"github.com/haapjari/glass/pkg/controllers/library"
	"github.com/haapjari/glass/pkg/controllers/metadata"
	"github.com/haapjari/glass/pkg/database"

	"github.com/gin-gonic/gin"
)

func SetupRouter() {

	// Setting Up Router and Database
	r := gin.Default()
	db := database.SetupDatabase()

	// Provide Database to Controllers
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Endpoint: /api/glass/repositories
	r.GET("/api/glass/v1/repositories", metadata.GetRepositories)
	r.POST("/api/glass/v1/repositories", metadata.CreateRepository)
	r.GET("/api/glass/v1/repositories/:id", metadata.GetRepositoryById)
	r.DELETE("/api/glass/v1/repositories/:id", metadata.DeleteRepositoryById)
	r.PATCH("/api/glass/v1/repositories/:id", metadata.UpdateRepositoryById)

	// Endpoint: /api/glass/v1/libraries
	r.GET("/api/glass/v1/libraries", library.GetLibraries)
	r.POST("/api/glass/v1/libraries", library.CreateLibrary)
	r.GET("/api/glass/v1/libraries/:id", library.GetLibraryById)
	r.DELETE("/api/glass/v1/libraries/:id", library.DeleteLibraryById)
	r.PATCH("/api/glass/v1/libraries/:id", library.UpdateLibraryById)

	// Endpoint: /api/glass/v1/data
	r.GET("/api/glass/v1/data", data.GetDataEntities)
	r.POST("/api/glass/v1/data", data.CreateDataEntity)
	r.GET("/api/glass/v1/data/:id", data.GetDataEntityById)
	r.DELETE("/api/glass/v1/data/:id", data.DeleteDataEntityById)
	r.PATCH("/api/glass/v1/data/:id", data.UpdateDataEntityById)
	r.POST("/api/glass/v1/data/fetch", data.FetchRepositories)
	r.POST("/api/glass/v1/data/enrich", data.EnrichMetadata)
	r.POST("/api/glass/v1/data/csv", data.GenerateCsvFile)

	r.Run()
}
