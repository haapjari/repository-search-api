package main

import (
	"github.com/gin-gonic/gin"
	"github.com/haapjari/repository-metadata-aggregator/api"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/cfg"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/logger"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/srv"
)

const (
	host = "0.0.0.0"
)

func main() {
	conf, err := cfg.NewConfig()
	if err != nil {
		return
	}

	router := gin.Default()

	log := logger.NewLogger(logger.DebugLevel)

	api.RegisterHandlers(router, srv.NewServer(log, conf))

	if err = router.Run(host + ":" + conf.Port); err != nil {
		log.Errorf("unable to run the server on port: %v, error: %v", conf.Port,
			err.Error())
		return
	}
}
