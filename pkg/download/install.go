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
	if alreadyInstalled(libPath) {
		isLatest, version, err := alreadyLatestVersion(libPath)
		if err != nil {
			return fmt.Errorf("error checking version installed: %w", err)
		}

		if isLatest {
			return nil
		}

		return upgradeInstall(libPath, version)
	}

	return initialInstall(libPath)
}

func alreadyInstalled(libPath string) bool {
	os.MkdirAll(libPath, 0755)

	versionInfoPath := filepath.Join(libPath, versionFile)

	if _, err := os.Stat(versionInfoPath); os.IsNotExist(err) {
		return false
	}

	return true
}

func alreadyLatestVersion(libPath string) (bool, string, error) {
	versionInfoPath := filepath.Join(libPath, versionFile)

	d, err := os.ReadFile(versionInfoPath)
	if err != nil {
		return false, "", fmt.Errorf("error reading version info file: %w", err)
	}

	var tag tag
	if err := json.Unmarshal(d, &tag); err != nil {
		return false, "", fmt.Errorf("error unmarshalling version info: %w", err)
	}

	version, err := LlamaLatestVersion()
	if err != nil {
		return false, "", fmt.Errorf("error install: %w", err)
	}

	return version == tag.TagName, version, nil
}

func initialInstall(libPath string) error {
	version, err := downloadVersionFile(llamaCppVersionDocURL)
	if err != nil {
		return fmt.Errorf("error downloading llama.cpp version document: %w", err)
	}

	return upgradeInstall(libPath, version)
}

func downloadVersionFile(llamaCppVersionDocURL string) (string, error) {
	r, err := http.DefaultClient.Get(llamaCppVersionDocURL)
	if err != nil {
		return "", fmt.Errorf("error getting llama.cpp version document: %w", err)
	}
	defer r.Body.Close()

	var tag tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		return "", fmt.Errorf("error decoding llama.cpp version document: %w", err)
	}

	return tag.TagName, nil
}

func upgradeInstall(libPath string, version string) error {
	if err := installLlamaCpp(libPath, version); err != nil {
		return fmt.Errorf("error installing %q of llama.cpp: %w", version, err)
	}

	if err := createVersionFile(libPath, version); err != nil {
		return fmt.Errorf("error creating version file: %w", err)
	}

	return nil
}

func installLlamaCpp(libPath string, version string) error {
	if _, err := os.Stat(libPath); !os.IsNotExist(err) {
		os.RemoveAll(libPath)
	}

	if err := Get(runtime.GOOS, "cpu", version, libPath); err != nil {
		return fmt.Errorf("error downloading llama.cpp: %w", err)
	}

	return nil
}

func createVersionFile(libPath string, version string) error {
	versionInfoPath := filepath.Join(libPath, versionFile)

	f, err := os.Create(versionInfoPath)
	if err != nil {
		return fmt.Errorf("error creating version info file: %w", err)
	}
	defer f.Close()

	t := tag{
		TagName: version,
	}

	d, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("error marshalling version info: %w", err)
	}

	if _, err := f.Write(d); err != nil {
		return fmt.Errorf("error writing version info: %w", err)
	}

	return nil
}
