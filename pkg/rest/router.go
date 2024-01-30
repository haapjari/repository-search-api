package router

import (
	"github.com/haapjari/glass/pkg/controllers"
	// "github.com/haapjari/glass/pkg/database"

	"github.com/gin-gonic/gin"
)

func SetupRouter() error {
	r := gin.Default()

	// TODO
	// db := database.SetupDatabase()

	// r.Use(func(c *gin.Context) {
	// 	c.Set("db", db)
	// 	c.Next()
	// })

	r.GET("/api/v1/repository/fetch", controllers.FetchRepositoryData)

	err := r.Run()
	if err != nil {
		return err
	}

	return nil
}
