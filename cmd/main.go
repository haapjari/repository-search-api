package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/haapjari/glass-api/api"
	"github.com/haapjari/glass-api/internal/pkg/server"
)

func main() {
	router := gin.Default()

	api.RegisterHandlers(router, server.NewServer())

	err := router.Run(":8080")
	if err != nil {
		slog.Error("unable to run the api: %v", err)
		return
	}
}
