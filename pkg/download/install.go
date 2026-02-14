package download

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var (
	llamaCppVersionDocURL = "https://api.github.com/repos/ggml-org/llama.cpp/releases/latest"
	versionFile           = "version.json"
)

type tag struct {
	TagName string `json:"tag_name"`
}

// AlreadyInstalled checks if llama.cpp is already installed at the given libPath. It does this
// by checking for the presence of the library file corresponding to the current OS. If the
// library file exists, it returns true, indicating that llama.cpp is already installed. If the
// library file does not exist, it returns false, indicating that llama.cpp is not installed.
func AlreadyInstalled(libPath string) bool {
	if _, err := os.Stat(filepath.Join(libPath, LibraryName(runtime.GOOS))); !os.IsNotExist(err) {
		return true
	}
	return false
}

// InstallLibraries has been deprecated. Use the `GetXXX` functions directly.
func InstallLibraries(libPath string, processor Processor, allowUpgrade bool) error {
	return fmt.Errorf("InstallLibraries is deprecated. Use the GetXXX functions directly")
}
