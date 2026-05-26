package logging

import (
	"log/slog"
	"regexp"

	"google.golang.org/grpc/status"
)

var methodInfo = regexp.MustCompile(`[./](\w+/\w+)$`)

// Err converts error to slog.Attr.
func Err(err error) slog.Attr {
	return slog.Any("error", err)
}

// GRPCErr converts gRPC error to slog.Attr.
// Returns "OK" if error is nil.
func GRPCErr(err error) slog.Attr {
	if status, ok := status.FromError(err); ok {
		return slog.String("status", status.Code().String())
	}

	return Err(err)
}

// GRPCMethod extracts service and call names from full description of gRPC method
// and converts them to slog.Attr.
// E.g.
// /examples.v1.Example/StreamStream -> Example/StreamStream
// /AudioRecorder/AddRecognizeRequestData -> AudioRecorder/AddRecognizeRequestData.
func GRPCMethod(fullMethod string) slog.Attr {
	result := methodInfo.FindStringSubmatch(fullMethod)
	if len(result) < 2 {
		return slog.String("method", "Unknown/Unknown")
	}

	return slog.String("method", result[1])
}
