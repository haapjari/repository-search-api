package handler

import (
	"github.com/haapjari/repository-search-api/internal/pkg/cfg"
)

type Handler struct {
	Config *cfg.Config
}

func NewHandler(config *cfg.Config) *Handler {
	return &Handler{
		Config: config,
	}
}
