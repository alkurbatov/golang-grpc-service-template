package config_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"

	"github.com/alkurbatov/golang-grpc-service-template/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	sut := config.New()

	snaps.MatchSnapshot(t, sut.String())
}
