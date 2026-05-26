// Package v1 implements version 1 of the gRPC API.
package v1

import (
	"log/slog"

	"google.golang.org/grpc"

	echopb "github.com/alkurbatov/golang-grpc-service-template/pkg/echopb/v1"
)

// RegisterRoutes injects new routes into the provided gRPC server.
func RegisterRoutes(logger *slog.Logger, server *grpc.Server) {
	echo := NewEchoServer(logger)
	echopb.RegisterEchoServer(server, echo)
}
