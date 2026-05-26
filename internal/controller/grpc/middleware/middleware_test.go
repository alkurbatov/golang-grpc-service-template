package middleware_test

import (
	"context"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/alkurbatov/golang-grpc-service-template/internal/controller/grpc/middleware"
	"github.com/alkurbatov/golang-grpc-service-template/internal/infra/logging"
)

func runLoggingUnaryInterceptor(
	ctx context.Context,
	t *testing.T,
	m *logging.HandlerMock,
	method string,
) {
	t.Helper()

	l := logging.NewLoggerMock(m)
	info := &grpc.UnaryServerInfo{FullMethod: method}

	_, err := middleware.LoggingUnaryInterceptor(l)(
		ctx,
		struct{}{},
		info,
		func(context.Context, any) (any, error) {
			return struct{}{}, nil
		},
	)

	require.NoError(t, err)
}

func runLoggingStreamInterceptor(
	ctx context.Context,
	t *testing.T,
	m *logging.HandlerMock,
	method string,
) {
	t.Helper()

	l := logging.NewLoggerMock(m)
	info := &grpc.StreamServerInfo{FullMethod: method}

	err := middleware.LoggingStreamInterceptor(l)(
		struct{}{},
		middleware.NewServerStreamMock(ctx),
		info,
		func(any, grpc.ServerStream) error {
			return nil
		},
	)

	require.NoError(t, err)
}

func assertRequestLogged(t *testing.T, m *logging.HandlerMock) {
	t.Helper()

	require.Len(t, m.Records, 2)
	snaps.MatchSnapshot(t, m.DumpRecordedMsg(0))
	snaps.MatchSnapshot(t, m.DumpRecordedMsg(1))
}

func TestLoggingInterceptorsIgnoreUnwantedMethods(t *testing.T) {
	tt := []struct {
		name   string
		method string
	}{
		{
			name:   "Reflection API",
			method: "/grpc.reflection.v1.ServerReflection/ServerReflectionInfo",
		},
		{
			name:   "HealthProbe API",
			method: "/grpc.health.v1.Health/Check",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := &logging.HandlerMock{}
			ctx := t.Context()

			runLoggingUnaryInterceptor(ctx, t, m, tc.method)
			runLoggingStreamInterceptor(ctx, t, m, tc.method)

			require.Empty(t, m.Records)
		})
	}
}

func TestLoggingInterceptors(t *testing.T) {
	ctx := t.Context()

	m := &logging.HandlerMock{}
	runLoggingUnaryInterceptor(ctx, t, m, "/examples.v1.Echo/UnaryUnary")
	assertRequestLogged(t, m)

	m = &logging.HandlerMock{}
	runLoggingStreamInterceptor(ctx, t, m, "/examples.v1.Echo/UnaryStream")
	assertRequestLogged(t, m)
}
