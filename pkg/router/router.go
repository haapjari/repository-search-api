package router

import (
	"github.com/haapjari/glass/pkg/controllers"
	// "github.com/haapjari/glass/pkg/database"

	"github.com/gin-gonic/gin"
)

func SetupRouter() {
	r := gin.Default()

    // TODO
	// db := database.SetupDatabase()

	// r.Use(func(c *gin.Context) {
	// 	c.Set("db", db)
	// 	c.Next()
	// })


	r.GET("/api/v1/repository", controllers.GetRepositories)
	r.POST("/api/v1/repository", controllers.CreateRepository)
	r.GET("/api/v1/repository/:id", controllers.GetRepositoryById)
	r.DELETE("/api/v1/repository/:id", controllers.DeleteRepositoryById)
	r.PATCH("/api/v1/repository/:id", controllers.UpdateRepositoryById)
	
    // TODO 
    r.GET("/api/v1/repository/fetch/metadata", repository.GetRepositoryMetadata)

	r.Run()
}
