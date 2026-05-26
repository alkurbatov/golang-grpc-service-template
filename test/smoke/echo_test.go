package smoke_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	echopb "github.com/alkurbatov/golang-grpc-service-template/pkg/echopb/v1"
)

type streamSender interface {
	Send(req *echopb.Request) error
}

func sendData(t *testing.T, stream streamSender, data []string) {
	t.Helper()

	for _, word := range data {
		err := stream.Send(&echopb.Request{Text: word})
		require.NoError(t, err)
	}
}

type streamReader interface {
	Recv() (*echopb.Response, error)
}

func drainStream(t *testing.T, stream streamReader) []string {
	t.Helper()

	var result []string

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return result
		}

		require.NoError(t, err)

		result = append(result, resp.GetText())
	}
}

type EchoSmokeTestSuite struct {
	*SmokeTestSuite

	// Client to gRPC echo API.
	client echopb.EchoClient
}

func NewEchoSmokeTestSuite(GRPCAddress string) *EchoSmokeTestSuite {
	return &EchoSmokeTestSuite{SmokeTestSuite: NewSmokeTestSuite(GRPCAddress)}
}

func (s *EchoSmokeTestSuite) SetupTest() {
	s.SmokeTestSuite.SetupTest()

	s.client = echopb.NewEchoClient(s.Conn)
}

func (s *EchoSmokeTestSuite) TestStreamStream() {
	stream, err := s.client.StreamStream(s.T().Context())
	s.Require().NoError(err)

	sendData(s.T(), stream, []string{"hello", "hellooo"})
	err = stream.CloseSend()
	s.Require().NoError(err)

	result := drainStream(s.T(), stream)

	s.Require().Equal([]string{"hello", "hellooo"}, result)
}

func (s *EchoSmokeTestSuite) TestStreamUnary() {
	stream, err := s.client.StreamUnary(s.T().Context())
	s.Require().NoError(err)

	sendData(s.T(), stream, []string{"hello", "hellooo"})
	result, err := stream.CloseAndRecv()

	s.Require().NoError(err)
	s.Require().Equal("hello hellooo", result.GetText())
}

func (s *EchoSmokeTestSuite) TestUnaryStream() {
	req := &echopb.Request{Text: "hello hellooo"}

	stream, err := s.client.UnaryStream(s.T().Context(), req)
	s.Require().NoError(err)

	result := drainStream(s.T(), stream)

	s.Require().Equal([]string{"hello", "hellooo"}, result)
}

func (s *EchoSmokeTestSuite) TestUnaryUnary() {
	req := &echopb.Request{Text: "hello"}

	result, err := s.client.UnaryUnary(s.T().Context(), req)

	s.Require().NoError(err)
	s.Require().Equal("hello", result.GetText())
}
