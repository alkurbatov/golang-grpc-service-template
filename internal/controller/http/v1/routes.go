package v1

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

// RegisterPublicRoutes initializes new router aka http.ServerMux with
// supported public HTTP endpoints.
func RegisterPublicRoutes(reg *prometheus.Registry) *http.ServeMux {
	router := http.NewServeMux()

	newMetricsResource(router, reg)

	return router
}

// RegisterPrivateRoutes initializes new router aka http.ServerMux with
// supported private HTTP endpoints.
// The private endpoints are usually not exposed to external users due to
// security concerns.
func RegisterPrivateRoutes() *http.ServeMux {
	router := http.NewServeMux()

	newProfilerResource(router)

	return router
}
