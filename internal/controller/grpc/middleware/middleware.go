// Package middleware provides various interceptors for gRPC API.
package middleware

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/alkurbatov/golang-grpc-service-template/internal/infra/logging"
)

// Default prefix used to generate request ID.
// Usually stands to abbreviation of the service name.
const servicePrefix = "GS"

// RequestIDKey is gRPC metadata key to get/set request ID.
const RequestIDKey = "request_id"

// shouldIgnore returns true, if the provided gRPC method shouldn't be handled
// by logging middleware, e.g. Reflection API requests, Health probes etc.
func shouldIgnore(method string) bool {
	return strings.HasSuffix(method, "Health/Check") ||
		strings.HasSuffix(method, "ServerReflection/ServerReflectionInfo")
}

// generateRequestID generates new request ID in form of service prefix + UUID.
func generateRequestID() string {
	return servicePrefix + "-" + uuid.New().String()
}

// injectRequestID injects request ID into gRPC metadata.
func injectRequestID(ctx context.Context, reqID string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		// NB (alkurbatov): Rather weird case with missing metadata.
		// Replace it with our data.
		md = metadata.MD{}
	}

	md.Set(RequestIDKey, reqID)

	return metadata.NewIncomingContext(ctx, md)
}

// prepareLoggingContext puts request tracing info into the provided context so that
// it was automagically mixed into log records during further calls to slog.
func prepareLoggingContext(ctx context.Context) context.Context {
	reqID := generateRequestID()
	ctx = injectRequestID(ctx, reqID)

	return logging.Context(ctx, slog.String(RequestIDKey, reqID))
}

// LoggingUnaryInterceptor logs incoming unary-unary requests and results.
func LoggingUnaryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		if shouldIgnore(info.FullMethod) {
			return handler(ctx, req)
		}

		start := time.Now()

		ctx = prepareLoggingContext(ctx)
		l := logger.With(logging.GRPCMethod(info.FullMethod))

		l.LogAttrs(ctx, slog.LevelInfo, "Incoming request")

		resp, err = handler(ctx, req)

		l.LogAttrs(ctx, slog.LevelInfo, "Processing ended",
			logging.GRPCErr(err),
			slog.Duration("duration", time.Since(start)),
		)

		return resp, err
	}
}

// serverStream wraps gRPC stream to inject custom context.
type serverStream struct {
	grpc.ServerStream

	ctx context.Context
}

// Context returns context overridden by serverStream wrapper.
func (s *serverStream) Context() context.Context {
	return s.ctx
}

// LoggingStreamInterceptor logs incoming stream-stream, stream-unary and unary-stream
// requests and responses.
func LoggingStreamInterceptor(logger *slog.Logger) grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if shouldIgnore(info.FullMethod) {
			return handler(srv, ss)
		}

		start := time.Now()

		ctx := prepareLoggingContext(ss.Context())
		l := logger.With(logging.GRPCMethod(info.FullMethod))

		l.LogAttrs(ctx, slog.LevelInfo, "Incoming request")

		err := handler(srv, &serverStream{ss, ctx})

		l.LogAttrs(ctx, slog.LevelInfo, "Processing ended",
			logging.GRPCErr(err),
			slog.Duration("duration", time.Since(start)),
		)

		return err
	}
}
