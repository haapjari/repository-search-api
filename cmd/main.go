package main

import (
	"github.com/gin-gonic/gin"
	"github.com/haapjari/glass-api/api"
	"github.com/haapjari/glass-api/internal/pkg/cfg"
	"github.com/haapjari/glass-api/internal/pkg/logger"
	"github.com/haapjari/glass-api/internal/pkg/srv"
)

func main() {
	r := gin.Default()
	c := cfg.NewConfig()
	logger := logger.NewLogger(logger.DebugLevel)

	api.RegisterHandlers(r, srv.NewServer(logger))

	logger.Debugf("GIN_MODE=%v", c.GinMode)

	// TODO: GIN MODE

	err := r.Run("127.0.0.1" + ":" + c.Port)
	if err != nil {
		logger.Errorf("unable to run the server on port: %v, error: %v", c.Port, err.Error())
		return
	}
}
