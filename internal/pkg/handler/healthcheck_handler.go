package handler

import (
	"log/slog"
	"net/http"
)

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug(r.Method + " " + r.RequestURI)

	if r.Method != http.MethodGet {
		slog.Warn("invalid request method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
