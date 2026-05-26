package smoke_test

import (
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SmokeTestSuite struct {
	suite.Suite

	// Address of GRPC API of the service under tests.
	GRPCAddress string

	// Connection to gRPC server.
	Conn *grpc.ClientConn
}

func NewSmokeTestSuite(GRPCAddress string) *SmokeTestSuite {
	return &SmokeTestSuite{GRPCAddress: GRPCAddress}
}

func (s *SmokeTestSuite) SetupTest() {
	conn, err := grpc.NewClient(s.GRPCAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}

	s.Conn = conn
}

func (s *SmokeTestSuite) TearDownTest() {
	if s.Conn != nil {
		s.Conn.Close()
	}
}
