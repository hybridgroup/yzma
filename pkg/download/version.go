package download

// DefaultVersion is the llama.cpp release that this yzma release installs when no
// version is asked for. An empty value means the most recent nightly build, which
// is the state on the main branch between yzma releases.
//
// It holds the complete pin, "<tag>@sha256:<manifest digest>", so that the default
// install authenticates its digest manifest. The llama-cpp-builder release notes for
// the tag print that value, and so does
// https://hybridgroup.github.io/llama-cpp-builder/version.json for the newest build.
//
// Set it in the same commit that bumps version.go, then set it back to "" after the
// release is tagged.
var DefaultVersion = ""

// DefaultTag gives [DefaultVersion] without its digest, which is the value to show a
// person. It gives "" when there is no default version.
func DefaultTag() string {
	tag, _, err := ParsePinnedVersion(DefaultVersion)
	if err != nil {
		return ""
	}

	return tag
}
