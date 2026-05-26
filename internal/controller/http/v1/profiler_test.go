package v1_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	v1 "github.com/alkurbatov/golang-grpc-service-template/internal/controller/http/v1"
)

func TestProfiler(t *testing.T) {
	tt := []struct {
		name     string
		endpoint string
	}{
		{
			name:     "Direct access",
			endpoint: "/debug/pprof",
		},
		{
			name:     "Redirect from root endpoint",
			endpoint: "/",
		},
		{
			name:     "Redirect from /pprof",
			endpoint: "/pprof",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			router := v1.RegisterPrivateRoutes()

			result := getFrom(t, router, tc.endpoint)

			require.Contains(t, result, "full goroutine stack dump")
		})
	}
}
