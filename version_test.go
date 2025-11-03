package yzma

import (
	"testing"
)

func TestVersion(t *testing.T) {
	version := Version()
	if version == "" {
		t.Fatal("version returned an empty string, which is invalid")
	}
	t.Logf("version returned: %s", version)
}
