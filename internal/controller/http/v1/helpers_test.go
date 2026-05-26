package v1_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	v1 "github.com/alkurbatov/golang-grpc-service-template/internal/controller/http/v1"
)

func newPublicTestRouter() *http.ServeMux {
	reg := prometheus.NewRegistry()

	return v1.RegisterPublicRoutes(reg)
}

func sendTestRequest(t *testing.T, router *http.ServeMux, method, path string) (int, []byte) {
	t.Helper()

	srv := httptest.NewServer(router)
	defer srv.Close()

	req, err := http.NewRequest(method, srv.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	return resp.StatusCode, extractResponse(t, resp)
}

func extractResponse(t *testing.T, resp *http.Response) []byte {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	resp.Body.Close()

	return body
}

func getFrom(t *testing.T, router *http.ServeMux, endpoint string) string {
	t.Helper()

	status, body := sendTestRequest(t, router, http.MethodGet, endpoint)
	require.Equal(t, http.StatusOK, status, endpoint)

	return string(body)
}
