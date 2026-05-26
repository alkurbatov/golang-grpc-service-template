package v1_test

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestMetrics(t *testing.T) {
	router := newPublicTestRouter()

	result := getFrom(t, router, "/metrics")

	snaps.MatchSnapshot(t, result)
}
