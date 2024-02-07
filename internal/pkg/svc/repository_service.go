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
	// What is the backend we need to call here?

	// GitHub Search API restricts a single query to return only up-to 1000 results. We need to go over 15 times over the limit.
	// To work around this restriction on the API, we do date range query splitting.

	// First we determine, when was the first and last repository created, this gives us a date range.
	// Then we split this date range to weeks.

	// First, find out the creation dates of the oldest and newest repositories that match your search criteria.
	// This can give you the overall date range to segment your queries.

	// We need to redesign the APIs to first return the date-range, then an API for the date range.

	// 2. Segment Your Search by Week

	// Based on the date range, you can create a series of queries that each cover a week.
	// Adjust your search query to include a date range using the created qualifier with the format YYYY-MM-DD..YYYY-MM-DD.

	// Connect to GitHub
	// Relay the Query
	// Parse the Query
	// Return the Query
}
