package router

import (
	"github.com/haapjari/glass/pkg/controllers/data"
	"github.com/haapjari/glass/pkg/controllers/library"
	"github.com/haapjari/glass/pkg/controllers/metadata"
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

	r.GET("/api/glass/v1/metadata", metadata.GetRepositories)
	r.POST("/api/glass/v1/metadata", metadata.CreateRepository)
	r.GET("/api/glass/v1/metadata/:id", metadata.GetRepositoryById)
	r.DELETE("/api/glass/v1/metadata/:id", metadata.DeleteRepositoryById)
	r.PATCH("/api/glass/v1/metadata/:id", metadata.UpdateRepositoryById)

	r.GET("/api/glass/v1/libraries", library.GetLibraries)
	r.POST("/api/glass/v1/libraries", library.CreateLibrary)
	r.GET("/api/glass/v1/libraries/:id", library.GetLibraryById)
	r.DELETE("/api/glass/v1/libraries/:id", library.DeleteLibraryById)
	r.PATCH("/api/glass/v1/libraries/:id", library.UpdateLibraryById)

	r.GET("/api/glass/v1/entity", data.GetEntities)
	r.POST("/api/glass/v1/entity", data.CreateEntity)
	r.GET("/api/glass/v1/entity/:id", data.GetEntityById)
	r.DELETE("/api/glass/v1/entity/:id", data.DeleteEntityById)
	r.PATCH("/api/glass/v1/entity/:id", data.UpdateEntityById)
	r.POST("/api/glass/v1/entity/fetch", data.FetchRepositories)
	r.POST("/api/glass/v1/entity/enrich", data.EnrichMetadata)
	r.POST("/api/glass/v1/entity/csv", data.GenerateCsvFile)

	r.Run()
}
