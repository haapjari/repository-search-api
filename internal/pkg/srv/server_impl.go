package srv

import (
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/cfg"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/logger"
)

type Server struct {
	log    logger.Logger
	config *cfg.Config
}

func NewServer(log logger.Logger, conf *cfg.Config) *Server {
	return &Server{
		log:    log,
		config: conf,
	}
}
