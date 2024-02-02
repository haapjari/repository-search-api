package main

import (
	"github.com/gin-gonic/gin"
	"github.com/haapjari/glass-api/api"
	"github.com/haapjari/glass-api/internal/pkg/server"
)

func main() {
	router := gin.Default()

	api.RegisterHandlers(router, server.NewServer())

	router.Run(":8080")
}
