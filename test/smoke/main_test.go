package smoke_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSmoke(t *testing.T) {
	suite.Run(t, NewEchoSmokeTestSuite("templatesrv:50051"))
}
