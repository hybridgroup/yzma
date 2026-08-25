package download

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	getter "github.com/hashicorp/go-getter"
)

var (
	ErrUnknownArch      = errors.New("unknown architecture")
	ErrUnknownOS        = errors.New("unknown OS")
	ErrUnknownProcessor = errors.New("unknown processor")
	ErrInvalidVersion   = errors.New("invalid version")
	ErrFileNotFound     = errors.New("could not download file: the requested llama.cpp version may still be building for your platform.")
)

var (
	// RetryCount is how many times the package will retry to obtain the latest llama.cpp version.
	RetryCount = 3
	// RetryDelay is the delay between retries when obtaining the latest llama.cpp version.
	RetryDelay = 3 * time.Second
	// versionURL is the URL for fetching the latest llama.cpp version.
	// We use the llama-cpp-builder repo instead of the original llama.cpp repo because
	// we need the precompiled binaries for certain platforms, and the build server might be
	// up to 1 hour out of sync with the latest commits to the original llama.cpp repo.
	//
	// Actual downloads will be from the llama.cpp repo for any builds that are available there,
	// and from the llama-cpp-builder repo for builds that are not available in the original repo
	// (e.g. ARM64 CUDA builds). This is handled in the getDownloadLocationAndFilename function.
	currentVersionURL = "https://hybridgroup.github.io/llama-cpp-builder/version.json"

	// previousVersionURL is the URL for fetching the previous llama.cpp version. This is used as a fallback
	// if the current version URL is not available or does not contain a valid version.
	// This is necessary because the build server for the current version might be building the latest version
	// and might not have it available yet, while the previous version is likely to be available and can be used as a fallback.
	previousVersionURL = "https://hybridgroup.github.io/llama-cpp-builder/previous.json"

	// nightlyPattern is the format of a llama.cpp nightly build tag, for example "b10620".
	nightlyPattern = regexp.MustCompile(`^b[0-9]+$`)

	// releasePattern is the format of a llama.cpp tagged release, for example "v0.3.0".
	releasePattern = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+`)

	// nightlyTagURL is the URL of the asset in a tagged llama.cpp release that gives
	// the nightly build tag. A tagged release has no binaries of its own.
	// https://github.com/ggml-org/llama.cpp/releases
	nightlyTagURL = "https://github.com/ggml-org/llama.cpp/releases/download/%s/nightly-tag.txt"
)

// LlamaLatestVersion fetches the latest release tag of llama.cpp from the version URL.
func LlamaLatestVersion() (string, error) {
	var version string
	var err error
	for range RetryCount {
		version, err = getLatestVersion()
		if err == nil {
			return version, nil
		}
		time.Sleep(RetryDelay)
	}

	return "", errors.New("unable to fetch latest version")
}

func getLatestVersion() (string, error) {
	req, err := http.NewRequest("GET", currentVersionURL, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("received status code %d from version URL: %s", resp.StatusCode, string(body))
	}

	var result struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if err := VersionIsValid(result.TagName); err != nil {
		return "", fmt.Errorf("%w: %s", err, result.TagName)
	}

	return result.TagName, nil
}

// LlamaPreviousVersion fetches the previous release tag of llama.cpp from the version URL.
func LlamaPreviousVersion() (string, error) {
	var version string
	var err error
	for range RetryCount {
		version, err = getPreviousVersion()
		if err == nil {
			return version, nil
		}
		time.Sleep(RetryDelay)
	}

	return "", errors.New("unable to fetch previous version")
}

func getPreviousVersion() (string, error) {
	req, err := http.NewRequest("GET", previousVersionURL, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("received status code %d from version URL: %s", resp.StatusCode, string(body))
	}

	var result struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if err := VersionIsValid(result.TagName); err != nil {
		return "", fmt.Errorf("%w: %s", err, result.TagName)
	}

	return result.TagName, nil
}

// getDownloadLocationAndFilename returns the download location and filename for the
// given parameters.
//
// Deprecated: use [DefaultResolver] or [Install]. Unlike the resolver, it downloads
// auxiliary assets (the Windows CUDA runtime) itself.
func getDownloadLocationAndFilename(arch Arch, os OS, prcssr Processor, version string, dest string) (location, filename string, err error) {
	urls, err := defaultResolve(Target{Arch: arch, OS: os, Processor: prcssr, Version: version})
	if err != nil {
		return "", "", err
	}
	for _, url := range urls[:len(urls)-1] {
		if err := get(context.Background(), url, dest, ProgressTracker); err != nil {
			return "", "", err
		}
	}
	last := urls[len(urls)-1]
	i := strings.LastIndex(last, "/")
	return last[:i], last[i+1:], nil
}

// getFunc is the function used to download files. It can be overridden for testing.
var getFunc = get

// Get downloads the llama.cpp precompiled binaries for the desired arch/OS/processor.
// arch can be one of the following values: "amd64", "arm64".
// os can be one of the following values: "linux", "darwin", "windows", "bookworm", "trixie".
// processor can be one of the following values: "cpu", "cuda", "metal", "rocm", "vulkan".
// version should be the desired llama.cpp version, either a `b1234` nightly build
// or a `v1.2.3` tagged release. If an empty
// string ("") or "latest" is provided, the latest release will be downloaded,
// with an automatic fallback to the previous version if the latest is still building.
// dest in the destination directory for the downloaded binaries.
func Get(architecture string, operatingSystem string, processor string, version string, dest string) error {
	return GetWithProgress(architecture, operatingSystem, processor, version, dest, ProgressTracker)
}

// GetWithProgress downloads the llama.cpp precompiled binaries for the desired arch/OS/processor
// using the provided progress tracker.
// arch can be one of the following values: "amd64", "arm64".
// os can be one of the following values: "linux", "darwin", "windows", "bookworm", "trixie".
// processor can be one of the following values: "cpu", "cuda", "metal", "rocm", "vulkan".
// version should be the desired llama.cpp version, either a `b1234` nightly build
// or a `v1.2.3` tagged release. If an empty
// string ("") or "latest" is provided, the latest release will be downloaded,
// with an automatic fallback to the previous version if the latest is still building.
// dest in the destination directory for the downloaded binaries.
func GetWithProgress(architecture string, operatingSystem string, processor string, version string, dest string, progress getter.ProgressTracker) error {
	return GetWithContext(context.Background(), architecture, operatingSystem, processor, version, dest, progress)
}

// GetWithContext downloads the llama.cpp precompiled binaries for the desired arch/OS/processor
// using the provided context and progress tracker.
// arch can be one of the following values: "amd64", "arm64".
// os can be one of the following values: "linux", "darwin", "windows", "bookworm", "trixie".
// processor can be one of the following values: "cpu", "cuda", "metal", "rocm", "vulkan".
// version should be the desired llama.cpp version, either a `b1234` nightly build
// or a `v1.2.3` tagged release. If an empty
// string ("") or "latest" is provided, the latest release will be downloaded,
// with an automatic fallback to the previous version if the latest is still building.
// dest in the destination directory for the downloaded binaries.
func GetWithContext(ctx context.Context, architecture string, operatingSystem string, processor string, version string, dest string, progress getter.ProgressTracker) error {
	arch, err := ParseArch(architecture)
	if err != nil {
		return ErrUnknownArch
	}

	os, err := ParseOS(operatingSystem)
	if err != nil {
		return ErrUnknownOS
	}

	prcssr, err := ParseProcessor(processor)
	if err != nil {
		return ErrUnknownProcessor
	}

	return Install(ctx, Target{Arch: arch, OS: os, Processor: prcssr, Version: version}, dest, progress, nil)
}

func get(ctx context.Context, url, dest string, progress getter.ProgressTracker) error {
	// Check if it's a .tar.gz file
	if strings.HasSuffix(url, ".tar.gz") {
		err := downloadAndExtractTarGz(url, dest, progress)
		if err != nil && strings.Contains(err.Error(), "404") {
			return fmt.Errorf("%w: %s", ErrFileNotFound, url)
		}
		return err
	}

	// Use go-getter for other file types (e.g., .zip)
	client := &getter.Client{
		Ctx:  ctx,
		Src:  url,
		Dst:  dest,
		Mode: getter.ClientModeAny,
	}

	if progress != nil {
		client.ProgressListener = progress
	}

	if err := client.Get(); err != nil {
		if strings.Contains(err.Error(), "404") {
			return fmt.Errorf("%w: %s", ErrFileNotFound, url)
		}
		return err
	}

	return nil
}

// downloadAndExtractTarGz downloads a .tar.gz file and extracts it to the destination directory.
func downloadAndExtractTarGz(url, dest string, progress getter.ProgressTracker) error {
	downloadFile := filepath.Join(dest, filepath.Base(url))

	client := &getter.Client{
		Ctx:  context.Background(),
		Src:  url + "?archive=false",
		Dst:  dest,
		Mode: getter.ClientModeAny,
	}

	if progress != nil {
		client.ProgressListener = progress
	}

	if err := client.Get(); err != nil {
		// Check for 404 errors specifically
		if strings.Contains(err.Error(), "404") {
			return fmt.Errorf("404 not found: %s", url)
		}
		return err
	}
	defer os.Remove(downloadFile)

	resp, err := os.Open(downloadFile)
	if err != nil {
		return fmt.Errorf("failed to open downloaded file: %w", err)
	}
	defer resp.Close()

	// Create gzip reader
	gzr, err := gzip.NewReader(resp)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	// Create tar reader
	tr := tar.NewReader(gzr)

	// Extract files
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// Strip the top-level directory (e.g., "llama-b1234/")
		name := header.Name
		if idx := strings.Index(name, "/"); idx != -1 {
			name = name[idx+1:]
		}

		// Skip empty names (the top-level directory itself)
		if name == "" {
			continue
		}

		target := filepath.Join(dest, filepath.Clean(name))

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			// Remove any existing entry first, so an upgrade replaces it instead of
			// writing through a stale symlink left by a previous install.
			if err := removeExisting(target); err != nil {
				return err
			}

			// Create the file
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}

			// Copy contents
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return fmt.Errorf("failed to write file: %w", err)
			}
			f.Close()
		case tar.TypeSymlink:
			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			// Remove any existing entry first. Keeping it would leave the version
			// symlinks (e.g. libllama.dylib) pointing at the previously installed
			// build, so an upgrade would have no effect at load time.
			if err := removeExisting(target); err != nil {
				return err
			}

			if err := os.Symlink(header.Linkname, target); err != nil {
				return fmt.Errorf("failed to create symlink: %w", err)
			}
		}
	}

	return nil
}

// removeExisting removes target unless it is already absent or a directory,
// which is left alone so extraction can populate it.
func removeExisting(target string) error {
	fi, err := os.Lstat(target)
	switch {
	case os.IsNotExist(err):
		return nil
	case err != nil:
		return fmt.Errorf("failed to stat %s: %w", target, err)
	case fi.IsDir():
		return nil
	}

	if err := os.Remove(target); err != nil {
		return fmt.Errorf("failed to remove existing %s: %w", target, err)
	}

	return nil
}

// VersionIsValid checks if the provided version string is valid.
func VersionIsValid(version string) error {
	if !nightlyPattern.MatchString(version) && !releasePattern.MatchString(version) {
		return ErrInvalidVersion
	}

	return nil
}

// IsTaggedRelease tells if version is a tagged llama.cpp release such as "v0.3.0",
// which needs [LlamaNightlyTag] to find the build with the binaries.
func IsTaggedRelease(version string) bool {
	return releasePattern.MatchString(version)
}

// LlamaNightlyTag returns the nightly build tag that has the binaries for a llama.cpp
// version. A nightly tag such as "b10620" gives itself. A tagged release such as
// "v0.3.0" has no binaries of its own, so the nightly build tag comes from the
// nightly-tag.txt asset of that release.
func LlamaNightlyTag(version string) (string, error) {
	if nightlyPattern.MatchString(version) {
		return version, nil
	}

	if err := VersionIsValid(version); err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(fmt.Sprintf(nightlyTagURL, version))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received status code %d for nightly tag of %s", resp.StatusCode, version)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 64))
	if err != nil {
		return "", err
	}

	tag := strings.TrimSpace(string(body))
	if !nightlyPattern.MatchString(tag) {
		return "", fmt.Errorf("%w: %s", ErrInvalidVersion, tag)
	}

	return tag, nil
}

// LibraryName returns the name for the llama.cpp library for any given OS.
func LibraryName(operatingSystem string) string {
	os, err := ParseOS(operatingSystem)
	if err != nil {
		return "unknown"
	}

	switch os {
	case Linux:
		return "libllama.so"
	case Windows:
		return "llama.dll"
	case Darwin:
		return "libllama.dylib"
	default:
		return "unknown"
	}
}
