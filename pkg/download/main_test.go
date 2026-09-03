//go:build !manifest

package download

import (
	"os"
	"testing"
)

// TestMain points the digest manifests at an address that answers nothing, so an
// ordinary test run does not reach the network. A test that needs a manifest sets
// digestsURL to its own server. The tests behind the "manifest" build tag talk to the
// real site on purpose, so this does not apply to them.
func TestMain(m *testing.M) {
	digestsURL = "http://127.0.0.1:0/digests/%s.json"
	os.Exit(m.Run())
}
