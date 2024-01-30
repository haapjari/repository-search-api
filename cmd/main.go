package main

import (
	"log/slog"

	"github.com/haapjari/glass/pkg/router"
)

func main() {
	if err := router.SetupRouter(); err != nil {
		slog.Error(err.Error())
	}
}
