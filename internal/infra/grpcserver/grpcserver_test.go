package grpcserver_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/alkurbatov/golang-grpc-service-template/internal/infra/grpcserver"
)

func startServer(t *testing.T) (*grpcserver.Server, context.Context) {
	t.Helper()

	srv := grpcserver.New(true)

	ctx, err := srv.Start("")
	require.NoError(t, err)

	return srv, ctx
}

func TestStartWithBadAddress(t *testing.T) {
	sut := grpcserver.New(true)

	_, err := sut.Start("100")

	require.Error(t, err)
}

func TestShutdown(t *testing.T) {
	tt := []struct {
		name     string
		timeout  time.Duration
		expected bool
	}{
		{
			name:     "Server shutdowns before timeout expiration",
			timeout:  2 * time.Second,
			expected: true,
		},
		{
			name:     "Timeout expires before service shutdown",
			timeout:  time.Nanosecond,
			expected: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			sut, ctx := startServer(t)

			result := sut.Shutdown(tc.timeout)

			<-ctx.Done()

			require.Equal(t, tc.expected, result)
			require.ErrorIs(t, context.Cause(ctx), grpc.ErrServerStopped)
		})
	}
}
