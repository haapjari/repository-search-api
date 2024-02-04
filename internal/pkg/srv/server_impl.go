package srv

import "github.com/haapjari/glass-api/internal/pkg/logger"

// Server implements
type Server struct {
	log logger.Logger
}

func NewServer(log logger.Logger) *Server {
	return &Server{
		log: log,
	}
}
