package download

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	getter "github.com/hashicorp/go-getter"
)

// installedBundle is a small stand-in for an extracted llama.cpp bundle: a file, a
// file in a subdirectory, and the symbolic link chain that the real bundles have.
type installedBundle struct {
	libPath string
	files   map[string]string
	links   map[string]string
}

// writeBundle puts a bundle on disk and gives the digests that describe it.
func writeBundle(t *testing.T) installedBundle {
	t.Helper()

	libPath := t.TempDir()
	bundle := installedBundle{
		libPath: libPath,
		files:   map[string]string{},
		links:   map[string]string{"libllama.so": "libllama.so.0", "libllama.so.0": "libllama.so.0.3.0"},
	}

	contents := map[string]string{
		"libllama.so.0.3.0": "the library",
		"lib/libggml.so":    "a backend",
	}
	for name, body := range contents {
		path := filepath.Join(libPath, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256([]byte(body))
		bundle.files[name] = hex.EncodeToString(sum[:])
	}

	for name, destination := range bundle.links {
		if err := os.Symlink(destination, filepath.Join(libPath, name)); err != nil {
			t.Fatal(err)
		}
	}

	return bundle
}

// serveBundleManifest publishes a manifest that describes a bundle, and writes the
// install record that points at it.
func serveBundleManifest(t *testing.T, bundle installedBundle, tag string) {
	t.Helper()

	const assetName = "llama-%s-bin-ubuntu-cpu-arm64.tar.gz"
	name := fmt.Sprintf(assetName, tag)

	body, err := json.Marshal(manifest{
		Version: 1,
		Tag:     tag,
		Sources: map[string]manifestSource{
			"hybridgroup/llama-cpp-builder": {
				Tag: tag,
				Assets: map[string]manifestAsset{
					name: {SHA256: "aaaa", Files: bundle.files, Links: bundle.links},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	}))
	t.Cleanup(server.Close)

	original := digestsURL
	digestsURL = server.URL + "/digests/%s.json"
	t.Cleanup(func() { digestsURL = original })

	record := InstallRecord{
		Tag: tag, Arch: "arm64", OS: "linux", Processor: "cpu",
		Assets: []Asset{{
			URL: fmt.Sprintf("https://github.com/hybridgroup/llama-cpp-builder/releases/download/%s/%s", tag, name),
		}},
	}
	if err := WriteInstallRecord(bundle.libPath, record); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyInstallOnAnUntouchedInstall(t *testing.T) {
	bundle := writeBundle(t)
	serveBundleManifest(t, bundle, "b10783")

	report, err := VerifyInstall(context.Background(), bundle.libPath, "")
	if err != nil {
		t.Fatalf("VerifyInstall() failed: %v", err)
	}

	if !report.OK() {
		t.Errorf("report is not OK: %+v", report)
	}
	// Two files and two links, and the record itself is not counted.
	if report.Verified != 4 || report.Changed != 0 || report.Missing != 0 || report.Unexpected != 0 {
		t.Errorf("got %d verified, %d changed, %d missing, %d unexpected; want 4/0/0/0",
			report.Verified, report.Changed, report.Missing, report.Unexpected)
	}
	if report.Tag != "b10783" {
		t.Errorf("Tag = %q, want b10783", report.Tag)
	}
}

func TestVerifyInstallFindsAChangedFile(t *testing.T) {
	bundle := writeBundle(t)
	serveBundleManifest(t, bundle, "b10783")

	path := filepath.Join(bundle.libPath, "libllama.so.0.3.0")
	if err := os.WriteFile(path, []byte("something else"), 0644); err != nil {
		t.Fatal(err)
	}

	report, err := VerifyInstall(context.Background(), bundle.libPath, "")
	if err != nil {
		t.Fatalf("VerifyInstall() failed: %v", err)
	}
	if report.OK() {
		t.Error("a changed library must not be OK")
	}
	if report.Changed != 1 {
		t.Errorf("Changed = %d, want 1", report.Changed)
	}
	if state := stateOf(report, "libllama.so.0.3.0"); state != FileChanged {
		t.Errorf("state = %v, want changed", state)
	}
}

func TestVerifyInstallFindsAMissingFile(t *testing.T) {
	bundle := writeBundle(t)
	serveBundleManifest(t, bundle, "b10783")

	if err := os.Remove(filepath.Join(bundle.libPath, "lib", "libggml.so")); err != nil {
		t.Fatal(err)
	}

	report, err := VerifyInstall(context.Background(), bundle.libPath, "")
	if err != nil {
		t.Fatalf("VerifyInstall() failed: %v", err)
	}
	if report.OK() {
		t.Error("a missing library must not be OK")
	}
	if state := stateOf(report, "lib/libggml.so"); state != FileMissing {
		t.Errorf("state = %v, want missing", state)
	}
}

func TestVerifyInstallFindsARepointedLink(t *testing.T) {
	bundle := writeBundle(t)
	serveBundleManifest(t, bundle, "b10783")

	path := filepath.Join(bundle.libPath, "libllama.so")
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("libevil.so", path); err != nil {
		t.Fatal(err)
	}

	report, err := VerifyInstall(context.Background(), bundle.libPath, "")
	if err != nil {
		t.Fatalf("VerifyInstall() failed: %v", err)
	}
	if report.OK() {
		t.Error("a link that points somewhere else must not be OK")
	}
	if state := stateOf(report, "libllama.so"); state != FileChanged {
		t.Errorf("state = %v, want changed", state)
	}
}

func TestVerifyInstallReportsAFileThatIsNotPartOfTheInstall(t *testing.T) {
	bundle := writeBundle(t)
	serveBundleManifest(t, bundle, "b10783")

	if err := os.WriteFile(filepath.Join(bundle.libPath, "libextra.so"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}

	report, err := VerifyInstall(context.Background(), bundle.libPath, "")
	if err != nil {
		t.Fatalf("VerifyInstall() failed: %v", err)
	}
	// A directory can hold more than one install, so this alone is not a failure.
	if !report.OK() {
		t.Error("a file that belongs to something else must not fail the check")
	}
	if report.Unexpected != 1 {
		t.Errorf("Unexpected = %d, want 1", report.Unexpected)
	}
	if state := stateOf(report, "libextra.so"); state != FileUnexpected {
		t.Errorf("state = %v, want unexpected", state)
	}
}

func TestVerifyInstallDoesNotTrustTheRecordTag(t *testing.T) {
	bundle := writeBundle(t)
	serveBundleManifest(t, bundle, "b10783")

	// Rewrite the record the way something that changed the libraries would.
	record, err := ReadInstallRecord(bundle.libPath)
	if err != nil {
		t.Fatal(err)
	}
	record.Tag = "b10780"
	if err := WriteInstallRecord(bundle.libPath, *record); err != nil {
		t.Fatal(err)
	}

	// The operator says which release must be there, so the assets are resolved
	// again for that release instead of taken from the record.
	report, err := VerifyInstall(context.Background(), bundle.libPath, "b10783")
	if err != nil {
		t.Fatalf("VerifyInstall() failed: %v", err)
	}
	if !report.OK() {
		t.Errorf("report is not OK: %+v", report)
	}
	if report.Tag != "b10783" {
		t.Errorf("Tag = %q, want the tag the caller asked for", report.Tag)
	}
}

func TestVerifyInstallWithNoRecord(t *testing.T) {
	_, err := VerifyInstall(context.Background(), t.TempDir(), "")
	if !errors.Is(err, ErrNoInstallRecord) {
		t.Errorf("VerifyInstall() = %v, want ErrNoInstallRecord", err)
	}
}

func TestVerifyInstallWithNoFileDigests(t *testing.T) {
	bundle := writeBundle(t)
	bundle.files = nil
	bundle.links = nil
	serveBundleManifest(t, bundle, "b10783")

	_, err := VerifyInstall(context.Background(), bundle.libPath, "")
	if !errors.Is(err, ErrNoFileDigests) {
		t.Errorf("VerifyInstall() = %v, want ErrNoFileDigests", err)
	}
}

func TestVerifyInstallWithARecordForAnotherRelease(t *testing.T) {
	bundle := writeBundle(t)
	serveBundleManifest(t, bundle, "b10783")

	record, err := ReadInstallRecord(bundle.libPath)
	if err != nil {
		t.Fatal(err)
	}
	record.Assets = []Asset{{
		URL: "https://github.com/hybridgroup/llama-cpp-builder/releases/download/b10000/llama-b10000-bin-ubuntu-cpu-arm64.tar.gz",
	}}
	if err := WriteInstallRecord(bundle.libPath, *record); err != nil {
		t.Fatal(err)
	}

	_, err = VerifyInstall(context.Background(), bundle.libPath, "")
	if !errors.Is(err, ErrRecordMismatch) {
		t.Errorf("VerifyInstall() = %v, want ErrRecordMismatch", err)
	}
}

func TestInstallWritesARecord(t *testing.T) {
	originalGet := getFunc
	getFunc = func(context.Context, Asset, string, getter.ProgressTracker) error { return nil }
	defer func() { getFunc = originalGet }()

	dest := t.TempDir()
	target := Target{Arch: ARM64, OS: Linux, Processor: CPU, Version: "b10783"}
	if err := Install(context.Background(), target, dest, nil, nil, WithVerify(VerifyOff)); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	record, err := ReadInstallRecord(dest)
	if err != nil {
		t.Fatalf("reading the record failed: %v", err)
	}
	if record.Tag != "b10783" || record.Arch != "arm64" || record.OS != "linux" || record.Processor != "cpu" {
		t.Errorf("record = %+v, want the target that was installed", record)
	}
	if len(record.Assets) != 1 {
		t.Errorf("record names %d assets, want 1", len(record.Assets))
	}
}

// stateOf gives the state a report holds for a name.
func stateOf(report *VerifyReport, name string) FileState {
	for _, file := range report.Files {
		if file.Name == name {
			return file.State
		}
	}
	return -1
}
