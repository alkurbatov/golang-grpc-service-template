package httpserver_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/alkurbatov/golang-grpc-service-template/internal/infra/httpserver"
)

func TestNewWithBadAddress(t *testing.T) {
	tt := []struct {
		name string
		sut  *httpserver.Server
	}{
		{
			name: "Production server",
			sut:  httpserver.NewProdServer(nil, "100"),
		},
		{
			name: "Default server",
			sut:  httpserver.New(nil, "100"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.sut.Start()

			<-ctx.Done()

			require.Error(t, context.Cause(ctx))
		})
	}
}

func TestStop(t *testing.T) {
	sut := httpserver.NewProdServer(nil, "")

	// NB (alkurbatov): Do second start to test that it returns the same context.
	sut.Start()
	ctx := sut.Start()

	// NB (alkurbatov): Give the server some time to start.
	time.Sleep(10 * time.Millisecond)

	err := sut.Stop()
	require.NoError(t, err)

	<-ctx.Done()

	require.ErrorIs(t, context.Cause(ctx), http.ErrServerClosed)
}

func TestStop_ServerNotRunning(t *testing.T) {
	sut := httpserver.NewProdServer(nil, "")

	err := sut.Stop()

	require.NoError(t, err)
}
