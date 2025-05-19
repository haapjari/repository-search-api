package main

import (
	"log/slog"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/haapjari/repository-search-api/internal/pkg/cfg"
	"github.com/haapjari/repository-search-api/internal/pkg/handler"
)

const (
	host = "0.0.0.0"
)

func main() {
	conf := cfg.NewConfig()

	h := handler.NewHandler(conf)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/repos/search", h.RepositoryHandler)

	if conf.EnablePprof {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		mux.HandleFunc("/debug/pprof/heap", pprof.Handler("heap").ServeHTTP)
		mux.HandleFunc("/debug/pprof/goroutine", pprof.Handler("goroutine").ServeHTTP)
		mux.HandleFunc("/debug/pprof/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
		mux.HandleFunc("/debug/pprof/block", pprof.Handler("block").ServeHTTP)
		mux.HandleFunc("/debug/pprof/mutex", pprof.Handler("mutex").ServeHTTP)
		mux.HandleFunc("/debug/pprof/allocs", pprof.Handler("allocs").ServeHTTP)
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))
	slog.Info("REST API | " + host + ":" + conf.Port)

	err := http.ListenAndServe(host+":"+conf.Port, mux)
	if err != nil {
		panic("unable to start the server: " + err.Error())
	}
}
