//go:build !manifest

package download

import (
	"os"
	"testing"
)

// TestMain points every digest source at an address that does not answer, so a normal
// test run stays off the network. A test that needs one points it at its own server.
// The tests behind the "manifest" build tag use the real site on purpose, so this does
// not apply to them.
func TestMain(m *testing.M) {
	manifestAssetURL = "http://127.0.0.1:0/releases/download/%[1]s/%[1]s.json"
	digestsURL = "http://127.0.0.1:0/digests/%s.json"
	releaseAPIURL = "http://127.0.0.1:0/releases/tags/%s"
	os.Exit(m.Run())
}
