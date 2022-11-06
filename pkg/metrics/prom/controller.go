package prom

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Wraps the functionality of Prometheus.

type Prom struct {
	Handler gin.HandlerFunc
}

func NewProm() *Prom {
	p := new(Prom)

	p.Handler = PrometheusHandler()

	return p
}

func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
