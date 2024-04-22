package handler

import (
	"encoding/json"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/models"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/svc"
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

func (h *Handler) ReposHandler(w http.ResponseWriter, r *http.Request) {
	q := &models.QueryParameters{
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

	service := svc.NewGitHubService(strings.SplitAfter(token, " ")[1], q)
	repos := service.Query()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(&models.RepositoryResponse{
		TotalCount: len(repos),
		Items:      repos,
	})
}
