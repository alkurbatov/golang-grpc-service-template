package v1_test

import (
	"context"
	"log/slog"
	"net"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	v1 "github.com/alkurbatov/golang-grpc-service-template/internal/controller/grpc/v1"
)

const bufSize = 1024 * 1024

type GRPCTestSuite struct {
	suite.Suite

	srv *grpc.Server

	// Connection to gRPC server, should be used to initialize new clients.
	Conn *grpc.ClientConn
}

func (s *GRPCTestSuite) SetupTest() {
	lis := bufconn.Listen(bufSize)
	s.srv = grpc.NewServer()

	v1.RegisterRoutes(slog.Default(), s.srv)

	go func() {
		if err := s.srv.Serve(lis); err != nil {
			panic(err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.NewClient(
		"passthrough:///",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	s.Conn = conn
}

func (s *GRPCTestSuite) TearDownTest() {
	if s.Conn != nil {
		s.Conn.Close()
	}

	if s.srv != nil {
		s.srv.Stop()
	}
}
