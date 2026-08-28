package download

// DefaultVersion is the llama.cpp release that this yzma release installs when no
// version is asked for. An empty value means the most recent nightly build, which
// is the state on the main branch between yzma releases.
//
// Set it to the llama.cpp tag for the release in the same commit that bumps
// version.go, then set it back to "" after the release is tagged.
var DefaultVersion = ""
