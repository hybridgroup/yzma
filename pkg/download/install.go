package download

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func InstallLibraries(libPath string) error {
	var tag tag

	switch doesVersionFileExist(libPath) {
	case true:
		isLatest, version, err := isLatestVersionInstalled(libPath)
		if err != nil {
			return fmt.Errorf("error checking version installed: %w", err)
		}

		if isLatest {
			return nil
		}

		tag.TagName = version

	case false:
		var err error
		tag, err = downloadVersionFile(llamaCppVersionDocURL)
		if err != nil {
			return fmt.Errorf("error downloading llama.cpp version document: %w", err)
		}
	}

	if err := installLlamaCpp(libPath, tag.TagName); err != nil {
		return fmt.Errorf("error installing %q of llama.cpp: %w", tag.TagName, err)
	}

	// -------------------------------------------------------------------------

	if err := createVersionFile(libPath, tag); err != nil {
		return fmt.Errorf("error creating version file: %w", err)
	}

	return nil
}

func doesVersionFileExist(libPath string) bool {
	os.MkdirAll(libPath, 0755)

	versionInfoPath := filepath.Join(libPath, versionFile)

	if _, err := os.Stat(versionInfoPath); os.IsNotExist(err) {
		return false
	}

	return true
}

func isLatestVersionInstalled(libPath string) (bool, string, error) {
	var tag struct {
		TagName string `json:"tag_name"`
	}

	versionInfoPath := filepath.Join(libPath, versionFile)

	d, err := os.ReadFile(versionInfoPath)
	if err != nil {
		return false, "", fmt.Errorf("error reading version info file: %w", err)
	}

	if err := json.Unmarshal(d, &tag); err != nil {
		return false, "", fmt.Errorf("error unmarshalling version info: %w", err)
	}

	version, err := LlamaLatestVersion()
	if err != nil {
		return false, "", fmt.Errorf("error install: %w", err)
	}

	return version == tag.TagName, version, nil
}

func downloadVersionFile(llamaCppVersionDocURL string) (tag, error) {
	r, err := http.DefaultClient.Get(llamaCppVersionDocURL)
	if err != nil {
		return tag{}, fmt.Errorf("error getting llama.cpp version document: %w", err)
	}
	defer r.Body.Close()

	var t tag
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return tag{}, fmt.Errorf("error decoding llama.cpp version document: %w", err)
	}

	return t, nil
}

func installLlamaCpp(libPath string, tagName string) error {
	if _, err := os.Stat(libPath); !os.IsNotExist(err) {
		os.RemoveAll(libPath)
	}

	if err := Get(runtime.GOOS, "cpu", tagName, libPath); err != nil {
		return fmt.Errorf("error downloading llama.cpp: %w", err)
	}

	return nil
}

func createVersionFile(libPath string, tag tag) error {
	versionInfoPath := filepath.Join(libPath, versionFile)

	f, err := os.Create(versionInfoPath)
	if err != nil {
		return fmt.Errorf("error creating version info file: %w", err)
	}
	defer f.Close()

	d, err := json.Marshal(tag)
	if err != nil {
		return fmt.Errorf("error marshalling version info: %w", err)
	}

	if _, err := f.Write(d); err != nil {
		return fmt.Errorf("error writing version info: %w", err)
	}

	return nil
}
