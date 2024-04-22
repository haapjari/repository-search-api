package main

import (
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/cfg"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/handler"
	"log/slog"
	"net/http"
	"os"
)

const (
	host = "0.0.0.0"
)

func main() {
	conf, err := cfg.NewConfig()
	if err != nil {
		panic("unable to create config: " + err.Error())
		return
	}

	h := handler.NewHandler(conf)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/repos/search", h.RepositoryHandler)

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	slog.Info("REST API | " + host + ":" + conf.Port)

	err = http.ListenAndServe(host+":"+conf.Port, mux)
	if err != nil {
		panic("unable to start the server: " + err.Error())
	}
}
