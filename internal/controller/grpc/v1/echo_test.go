package v1_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

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

type EchoTestSuite struct {
	GRPCTestSuite

	sut echopb.EchoClient
}

func TestEcho(t *testing.T) {
	s := EchoTestSuite{}
	suite.Run(t, &s)
}

func (s *EchoTestSuite) SetupTest() {
	s.GRPCTestSuite.SetupTest()
	s.sut = echopb.NewEchoClient(s.Conn)
}

func (s *EchoTestSuite) TestStreamStream() {
	stream, err := s.sut.StreamStream(s.T().Context())
	s.Require().NoError(err)

	sendData(s.T(), stream, []string{"hello", "hellooo"})
	err = stream.CloseSend()
	s.Require().NoError(err)

	result := drainStream(s.T(), stream)

	s.Require().Equal([]string{"hello", "hellooo"}, result)
}

func (s *EchoTestSuite) TestStreamUnary() {
	stream, err := s.sut.StreamUnary(s.T().Context())
	s.Require().NoError(err)

	sendData(s.T(), stream, []string{"hello", "hellooo"})
	result, err := stream.CloseAndRecv()

	s.Require().NoError(err)
	s.Require().Equal("hello hellooo", result.GetText())
}

func (s *EchoTestSuite) TestUnaryStream() {
	req := &echopb.Request{Text: "hello hellooo"}

	stream, err := s.sut.UnaryStream(s.T().Context(), req)
	s.Require().NoError(err)

	result := drainStream(s.T(), stream)

	s.Require().Equal([]string{"hello", "hellooo"}, result)
}

func (s *EchoTestSuite) TestUnaryUnary() {
	req := &echopb.Request{Text: "hello"}

	result, err := s.sut.UnaryUnary(s.T().Context(), req)

	s.Require().NoError(err)
	s.Require().Equal("hello", result.GetText())
}
