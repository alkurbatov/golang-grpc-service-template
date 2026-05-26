package logging_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/alkurbatov/golang-grpc-service-template/internal/infra/logging"
)

func TestErr(t *testing.T) {
	result := logging.Err(io.EOF)

	require.Equal(t, "error=EOF", result.String())
}

func TestGRPCErr(t *testing.T) {
	tt := []struct {
		name     string
		code     codes.Code
		expected string
	}{
		{
			name:     "No error",
			code:     codes.OK,
			expected: "status=OK",
		},
		{
			name:     "Some error",
			code:     codes.InvalidArgument,
			expected: "status=InvalidArgument",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := status.Error(tc.code, "error reason")
			result := logging.GRPCErr(err)

			require.Equal(t, tc.expected, result.String())
		})
	}
}

func TestGRPCErrOnOtherErrors(t *testing.T) {
	tt := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "No error",
			expected: "status=OK",
		},
		{
			name:     "Error from other package",
			err:      io.EOF,
			expected: "error=EOF",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := logging.GRPCErr(tc.err)

			require.Equal(t, tc.expected, result.String())
		})
	}
}

func TestGRPCMethod(t *testing.T) {
	tt := []struct {
		name       string
		fullMethod string
		expected   string
	}{
		{
			name:     "Empty input",
			expected: "method=Unknown/Unknown",
		},
		{
			name:       "Malformed input",
			fullMethod: "something-unexpected",
			expected:   "method=Unknown/Unknown",
		},
		{
			name:       "Method without namespace and version",
			fullMethod: "/AudioRecorder/AddRecognizeRequestData",
			expected:   "method=AudioRecorder/AddRecognizeRequestData",
		},
		{
			name:       "Method with namespace and version",
			fullMethod: "/examples.v1.Example/StreamStream",
			expected:   "method=Example/StreamStream",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := logging.GRPCMethod(tc.fullMethod)

			require.Equal(t, tc.expected, result.String())
		})
	}
}
