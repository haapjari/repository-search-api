package model

import (
	"encoding/json"
	"github.com/haapjari/repository-search-api/internal/pkg/util"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type RepositoryResponse struct {
	TotalCount int           `json:"total_count"`
	Items      []*Repository `json:"items"`
}

type Repository struct {
	Name                   string `json:"name"`
	FullName               string `json:"full_name"`
	CreatedAt              string `json:"created_at"`
	StargazerCount         int    `json:"stargazer_count"`
	Language               string `json:"language"`
	OpenIssues             int    `json:"open_issues"`
	ClosedIssues           int    `json:"closed_issues"`
	OpenPullRequestCount   int    `json:"open_pull_request_count"`
	ClosedPullRequestCount int    `json:"closed_pull_request_count"`
	Forks                  int    `json:"forks"`
	WatcherCount           int    `json:"watcher_count"`
	CommitCount            int    `json:"commit_count"`
	NetworkCount           int    `json:"network_count"`
	LatestRelease          string `json:"latest_release"`
	TotalReleasesCount     int    `json:"total_releases_count"`
	ContributorCount       int    `json:"contributor_count"`
	ThirdPartyLOC          int    `json:"third_party_loc"`
	SelfWrittenLOC         int    `json:"self_written_loc"`
}

type QueryParameters struct {
	FirstCreationDate string
	LastCreationDate  string
	Language          string
	MinStars          string
	MaxStars          string
	Order             string
}

func (q *QueryParameters) ToString() string {
	b, _ := json.Marshal(q)
	return string(b)
}

func (q *QueryParameters) Validate() bool {
	if q.FirstCreationDate == "" {
		slog.Warn("empty first creation date parameter")
		return false
	}

	if _, err := time.Parse("2006-01-02", q.FirstCreationDate); err != nil {
		slog.Warn("invalid or malformed first creation date")
		return false
	}

	if q.LastCreationDate == "" {
		slog.Warn("empty last creation date parameter")
		return false
	}

	if _, err := time.Parse("2006-01-02", q.LastCreationDate); err != nil {
		slog.Warn("invalid last creation date")
		return false
	}

	if q.Language == "" || !util.OnlyLetters(q.Language) {
		slog.Warn("empty or malformed language parameter")
		return false
	}

	if q.MinStars == "" {
		slog.Warn("empty min stars parameter")
		return false
	}

	if _, err := strconv.Atoi(q.MinStars); err != nil {
		slog.Warn("invalid min stars parameter")
		return false
	}

	if q.MaxStars == "" {
		slog.Warn("empty max stars parameter")
		return false
	}

	if _, err := strconv.Atoi(q.MaxStars); err != nil {
		slog.Warn("invalid max stars parameter")
		return false
	}

	if q.Order == "" || !(strings.EqualFold(q.Order, "asc") ||
		strings.EqualFold(q.Order, "desc")) {
		slog.Warn("invalid or missing order parameter")
		return false
	}

	return true
}
