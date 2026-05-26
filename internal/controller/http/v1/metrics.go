package v1

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func newMetricsResource(router *http.ServeMux, reg *prometheus.Registry) {
	h := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})
	router.Handle("/", h)
	router.Handle("/metrics", h)
}
