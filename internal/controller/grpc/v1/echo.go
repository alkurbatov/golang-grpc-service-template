package v1

import (
	"context"
	"io"
	"log/slog"
	"strings"

	echopb "github.com/alkurbatov/golang-grpc-service-template/pkg/echopb/v1"
)

// EchoServer provides implementation of the example Echo API.
type EchoServer struct {
	echopb.UnimplementedEchoServer

	log *slog.Logger
}

// NewEchoServer creates new EchoServer.
func NewEchoServer(logger *slog.Logger) *EchoServer {
	return &EchoServer{log: logger}
}

// StreamStream provides example of stream-stream gRPC request handling.
func (s *EchoServer) StreamStream(stream echopb.Echo_StreamStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		s.log.InfoContext(stream.Context(), "Received chunk", slog.Any("req", req))

		resp := &echopb.Response{Text: req.GetText()}
		s.log.InfoContext(stream.Context(), "Sending chunk", slog.Any("resp", resp))

		if err = stream.Send(resp); err != nil {
			return err
		}
	}
}

// StreamUnary provides example of stream-unary gRPC request handling.
func (s *EchoServer) StreamUnary(stream echopb.Echo_StreamUnaryServer) error {
	var result []string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		s.log.InfoContext(stream.Context(), "Received chunk", slog.Any("req", req))
		result = append(result, req.GetText())
	}

	resp := &echopb.Response{Text: strings.Join(result, " ")}
	s.log.InfoContext(stream.Context(), "Sending response", slog.Any("resp", resp))

	return stream.SendAndClose(resp)
}

// UnaryStream provides example of unary-stream gRPC request handling.
func (s *EchoServer) UnaryStream(req *echopb.Request, stream echopb.Echo_UnaryStreamServer) error {
	s.log.InfoContext(stream.Context(), "Received request", slog.Any("req", req))

	for word := range strings.FieldsSeq(req.GetText()) {
		resp := &echopb.Response{Text: word}
		s.log.InfoContext(stream.Context(), "Sending chunk", slog.Any("resp", resp))

		if err := stream.Send(resp); err != nil {
			return err
		}
	}

	return nil
}

// UnaryUnary provides example of unary-unary gRPC request handling.
func (s *EchoServer) UnaryUnary(
	ctx context.Context,
	req *echopb.Request,
) (*echopb.Response, error) {
	s.log.InfoContext(ctx, "Received request", slog.Any("req", req))

	resp := &echopb.Response{Text: req.GetText()}
	s.log.InfoContext(ctx, "Sending response", slog.Any("resp", resp))

	return resp, nil
}
