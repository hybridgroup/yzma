package download

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

var llamaCppVersionDocURL = "https://api.github.com/repos/ggml-org/llama.cpp/releases/latest"

func InstallLibraries(libPath string) error {
	var tag struct {
		TagName string `json:"tag_name"`
	}

	os.MkdirAll(libPath, 0755)
	versionInfoPath := filepath.Join(libPath, "version.json")

	// -------------------------------------------------------------------------

	switch _, err := os.Stat(versionInfoPath); os.IsNotExist(err) {
	case false:
		d, err := os.ReadFile(versionInfoPath)
		if err != nil {
			return fmt.Errorf("error reading version info file: %w", err)
		}

		if err := json.Unmarshal(d, &tag); err != nil {
			return fmt.Errorf("error unmarshalling version info: %w", err)
		}

		version, err := LlamaLatestVersion()
		if err != nil {
			return fmt.Errorf("error install: %w", err)
		}

		if version == tag.TagName {
			return nil
		}

		tag.TagName = version

	case true:
		r, err := http.DefaultClient.Get(llamaCppVersionDocURL)
		if err != nil {
			return fmt.Errorf("error getting llama.cpp version document: %w", err)
		}
		defer r.Body.Close()

		if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
			return fmt.Errorf("error decoding llama.cpp version document: %w", err)
		}
	}

	// -------------------------------------------------------------------------

	if _, err := os.Stat(libPath); !os.IsNotExist(err) {
		os.RemoveAll(libPath)
	}

	if err := Get(runtime.GOOS, "cpu", tag.TagName, libPath); err != nil {
		return fmt.Errorf("error downloading llama.cpp: %w", err)
	}

	// -------------------------------------------------------------------------

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
