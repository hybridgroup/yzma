package download

// DefaultVersion is the llama.cpp release that this yzma release installs when no
// version is given. An empty value means the most recent nightly build, which is
// how the main branch stays between yzma releases.
//
// It holds the full pin, "<tag>@sha256:<manifest digest>", so the default install
// authenticates its digest manifest. The llama-cpp-builder release notes for the
// tag print that value, and so does
// https://hybridgroup.github.io/llama-cpp-builder/version.json for the newest build.
//
// Set it in the same commit that bumps version.go, then set it back to "" after the
// release is tagged.
var DefaultVersion = "v0.6.0@sha256:eae6acae7a2044bb27981d2bd72481780b5675d92b6adaefe5035364160d8153"

// DefaultTag returns [DefaultVersion] without its digest, for display to users. It
// returns "" when there is no default version.
func DefaultTag() string {
	tag, _, err := ParsePinnedVersion(DefaultVersion)
	if err != nil {
		return ""
	}

	return tag
}
