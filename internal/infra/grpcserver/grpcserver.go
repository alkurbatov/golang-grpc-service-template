// Package grpcserver implements handy wrap around gRPC server
// to group common settings and tasks inside single entity.
package grpcserver

import (
	"context"
	"math"
	"net"
	"slices"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const defaultMaxReceiveMessageSize = math.MaxInt32

// Server wraps gRPC server entity and provides handy means
// to simplify work with it.
type Server struct {
	Server *grpc.Server
}

// New provides new instance of gRPC server.
func New(withReflection bool, opts ...grpc.ServerOption) *Server {
	srvOpts := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(defaultMaxReceiveMessageSize),
	}
	srvOpts = slices.Concat(srvOpts, opts)

	server := grpc.NewServer(srvOpts...)

	// NB (alkurbatov): Enable reflection to use tools like grpcurl.
	// See: https://github.com/grpc/grpc/blob/master/doc/server-reflection.md
	//
	// Warning: don't enable in production as it cause security risk
	// (attacker can identify the binary API calls and use them for own needs).
	if withReflection {
		reflection.Register(server)
	}

	return &Server{Server: server}
}

// Start launches the gRPC server.
func (s *Server) Start(address string) (context.Context, error) {
	ctx, cancel := context.WithCancelCause(context.Background())

	lc := net.ListenConfig{}

	listen, err := lc.Listen(ctx, "tcp", address)
	if err != nil {
		cancel(err)
		return nil, err
	}

	go func() {
		if err = s.Server.Serve(listen); err != nil {
			cancel(err)
			return
		}

		cancel(grpc.ErrServerStopped)
	}()

	return ctx, nil
}

// Shutdown gracefully stops the server.
// Returns true, if service stopped gracefully.
// If timeout has expired, stop the server forcibly and return false.
func (s *Server) Shutdown(timeout time.Duration) bool {
	stopped := make(chan struct{})

	go func() {
		s.Server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-time.After(timeout):
		s.Server.Stop()
		return false

	case <-stopped:
		return true
	}
}
