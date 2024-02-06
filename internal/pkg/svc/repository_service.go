package svc

import (
	"github.com/haapjari/glass-api/internal/pkg/logger"
)

// github service knows how to authenticate and communicate with github

type RepositoryService struct {
	log logger.Logger
}

func NewRepositoryService(logger logger.Logger) *RepositoryService {
	return &RepositoryService{
		log: logger,
	}
}

// SearchPages is a service method for a Repository Handler.
func (s *RepositoryService) SearchPages(language string, stars string) {
	// TODO
	s.log.Debugf("Querying GitHub for %v languaged repositories, with %v definition for stars", language, stars)

	// gitHubService := github.NewGitHubService()

	// Connect to GitHub
	// Relay the Query
	// Parse the Query
	// Return the Query
}
