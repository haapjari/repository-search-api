package handler

import (
	"encoding/json"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/model"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/service"
	"log/slog"
	"net/http"
	"strings"
)

const (
	FirstCreationDate string = "firstCreationDate"
	LastCreationDate  string = "lastCreationDate"
	Language          string = "language"
	MinStars          string = "minStars"
	MaxStars          string = "maxStars"
	Order             string = "order"
)

func (h *Handler) RepositoryHandler(w http.ResponseWriter, r *http.Request) {
	q := &model.QueryParameters{
		FirstCreationDate: r.URL.Query().Get(FirstCreationDate),
		LastCreationDate:  r.URL.Query().Get(LastCreationDate),
		Language:          r.URL.Query().Get(Language),
		MinStars:          r.URL.Query().Get(MinStars),
		MaxStars:          r.URL.Query().Get(MaxStars),
		Order:             r.URL.Query().Get(Order),
	}

	if !q.Validate() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token := r.Header.Get("Authorization")

	if token == "" || len(strings.TrimSpace(token)) < 2 {
		slog.Warn("empty or malformed authorization header")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodGet {
		slog.Warn("invalid request method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	slog.Debug(r.Method + " " + r.RequestURI)

	svc := service.NewGitHubService(token, q)

	repos, err := svc.Query()
	if err != nil {
		slog.Error("unable to query the repositories: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&model.RepositoryResponse{
		TotalCount: len(repos),
		Items:      repos,
	})
}
