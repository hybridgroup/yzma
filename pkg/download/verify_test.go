package download

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	getter "github.com/hashicorp/go-getter"
)

// installedBundle is a small stand-in for an extracted llama.cpp bundle. It has a
// file, a file in a subdirectory, and the symbolic link chain the real bundles have.
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

// bundleManifest gives the manifest that describes a bundle, and the assets it covers.
// The assets are the ones the resolver names, so a check that resolves them again finds
// the same ones. upstream is the nightly build that holds the binaries of a tagged
// release, and is empty for a nightly tag.
func bundleManifest(t *testing.T, bundle installedBundle, tag, upstream string) ([]byte, []Asset) {
	t.Helper()

	urls, err := defaultResolve(Target{
		Arch: ARM64, OS: Linux, Processor: CPU, Version: tag, UpstreamVersion: upstream,
	})
	if err != nil {
		t.Fatal(err)
	}

	sources := map[string]manifestSource{}
	assets := make([]Asset, len(urls))
	for i, url := range urls {
		parts := assetURLPattern.FindStringSubmatch(url)
		if parts == nil {
			t.Fatalf("cannot read the asset URL %s", url)
		}

		source := sources[parts[1]]
		source.Tag = parts[2]
		if source.Assets == nil {
			source.Assets = map[string]manifestAsset{}
		}

		// Only the first asset holds the bundle.
		entry := manifestAsset{SHA256: "aaaa"}
		if i == 0 {
			entry.Files, entry.Links = bundle.files, bundle.links
		}
		source.Assets[parts[3]] = entry
		sources[parts[1]] = source

		assets[i] = Asset{URL: url}
	}

	upstreamTag := upstream
	if upstreamTag == "" {
		upstreamTag = tag
	}

	body, err := json.Marshal(manifest{
		Version: 1, Tag: tag, UpstreamTag: upstreamTag, Sources: sources,
	})
	if err != nil {
		t.Fatal(err)
	}

	return body, assets
}

// serveManifestBody publishes manifest bytes for the length of the test.
func serveManifestBody(t *testing.T, body []byte) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	}))
	t.Cleanup(server.Close)

	original := digestsURL
	digestsURL = server.URL + "/digests/%s.json"
	t.Cleanup(func() { digestsURL = original })

	return server
}

// writeBundleRecord writes the install record of a bundle.
func writeBundleRecord(t *testing.T, bundle installedBundle, tag, upstream string, assets []Asset, manifestDigest string) {
	t.Helper()

	record := InstallRecord{
		Tag: tag, UpstreamTag: upstream, Arch: "arm64", OS: "linux", Processor: "cpu",
		ManifestSHA256: manifestDigest, Assets: assets,
	}
	if err := WriteInstallRecord(bundle.libPath, record); err != nil {
		t.Fatal(err)
	}
}

// serveBundleManifest publishes a manifest that describes a bundle, and writes the
// install record that points at it. The record has no manifest beside it, as the records
// that earlier releases wrote do.
func serveBundleManifest(t *testing.T, bundle installedBundle, tag string) (*httptest.Server, []byte) {
	t.Helper()

	body, assets := bundleManifest(t, bundle, tag, "")
	server := serveManifestBody(t, body)
	writeBundleRecord(t, bundle, tag, "", assets, "")

	return server, body
}

// installBundle writes the record and the manifest that an install leaves beside a
// bundle, and gives the digest of the manifest to pin with. Nothing is published, so a
// check that reaches the network fails.
func installBundle(t *testing.T, bundle installedBundle, tag, upstream string) string {
	t.Helper()

	body, assets := bundleManifest(t, bundle, tag, upstream)
	sum := sha256.Sum256(body)
	digest := hex.EncodeToString(sum[:])

	writeBundleRecord(t, bundle, tag, upstream, assets, digest)
	if err := WriteInstallManifest(bundle.libPath, body); err != nil {
		t.Fatal(err)
	}

	return digest
}

// blockNightlyTag points the nightly tag asset at a port that answers nothing, so a
// lookup that must not happen fails the test rather than reaching GitHub.
func blockNightlyTag(t *testing.T) {
	t.Helper()

	original := nightlyTagURL
	nightlyTagURL = "http://127.0.0.1:0/%s/nightly-tag.txt"
	t.Cleanup(func() { nightlyTagURL = original })
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

func TestVerifyInstallReadsTheManifestBesideTheRecord(t *testing.T) {
	bundle := writeBundle(t)
	installBundle(t, bundle, "b10783", "")

	// Nothing is published, so a fetch of the manifest fails the test.
	report, err := VerifyInstall(context.Background(), bundle.libPath, "")
	if err != nil {
		t.Fatalf("VerifyInstall() failed: %v", err)
	}
	if !report.OK() {
		t.Errorf("report is not OK: %+v", report)
	}
	// The manifest belongs to the install, the same as the record.
	if report.Unexpected != 0 {
		t.Errorf("Unexpected = %d, want 0", report.Unexpected)
	}
}

func TestVerifyInstallOfAPinnedTaggedReleaseNeedsNoNetwork(t *testing.T) {
	blockNightlyTag(t)

	bundle := writeBundle(t)
	digest := installBundle(t, bundle, "v0.4.0", "b10783")

	// The manifest names the nightly build, so neither it nor the release page is
	// fetched.
	report, err := VerifyInstall(context.Background(), bundle.libPath, "v0.4.0@sha256:"+digest)
	if err != nil {
		t.Fatalf("VerifyInstall() failed: %v", err)
	}
	if !report.OK() {
		t.Errorf("report is not OK: %+v", report)
	}
	if report.Tag != "v0.4.0" {
		t.Errorf("Tag = %q, want v0.4.0", report.Tag)
	}
}

func TestVerifyInstallFetchesWhenTheManifestOnDiskIsBad(t *testing.T) {
	bundle := writeBundle(t)
	installBundle(t, bundle, "b10783", "")
	body, _ := bundleManifest(t, bundle, "b10783", "")
	serveManifestBody(t, body)

	if err := WriteInstallManifest(bundle.libPath, []byte("not a manifest")); err != nil {
		t.Fatal(err)
	}

	report, err := VerifyInstall(context.Background(), bundle.libPath, "")
	if err != nil {
		t.Fatalf("VerifyInstall() failed: %v", err)
	}
	if !report.OK() {
		t.Errorf("report is not OK: %+v", report)
	}
}

func TestVerifyInstallIgnoresAManifestForAnotherRelease(t *testing.T) {
	bundle := writeBundle(t)
	installBundle(t, bundle, "b10783", "")
	body, _ := bundleManifest(t, bundle, "b10783", "")
	serveManifestBody(t, body)

	// A manifest of another release does not describe this install, even when the
	// bytes are whole.
	other, _ := bundleManifest(t, bundle, "b10780", "")
	if err := WriteInstallManifest(bundle.libPath, other); err != nil {
		t.Fatal(err)
	}

	report, err := VerifyInstall(context.Background(), bundle.libPath, "")
	if err != nil {
		t.Fatalf("VerifyInstall() failed: %v", err)
	}
	if !report.OK() {
		t.Errorf("report is not OK: %+v", report)
	}
}

func TestVerifyInstallKeepsTheManifestItFetches(t *testing.T) {
	bundle := writeBundle(t)
	server, body := serveBundleManifest(t, bundle, "b10783")

	before, err := ReadInstallRecord(bundle.libPath)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := VerifyInstall(context.Background(), bundle.libPath, ""); err != nil {
		t.Fatalf("VerifyInstall() failed: %v", err)
	}

	kept, err := ReadInstallManifest(bundle.libPath)
	if err != nil {
		t.Fatalf("the manifest was not kept: %v", err)
	}
	if !bytes.Equal(kept, body) {
		t.Error("the manifest that was kept is not the one that was fetched")
	}

	sum := sha256.Sum256(body)
	record, err := ReadInstallRecord(bundle.libPath)
	if err != nil {
		t.Fatal(err)
	}
	if record.ManifestSHA256 != hex.EncodeToString(sum[:]) {
		t.Errorf("record ManifestSHA256 = %q, want the digest of the manifest", record.ManifestSHA256)
	}
	if !record.Installed.Equal(before.Installed) || len(record.Assets) != len(before.Assets) {
		t.Errorf("record = %+v, want the rest of it unchanged", record)
	}

	// The next check reads what was kept.
	server.Close()
	report, err := VerifyInstall(context.Background(), bundle.libPath, "")
	if err != nil {
		t.Fatalf("VerifyInstall() after the manifest was kept failed: %v", err)
	}
	if !report.OK() {
		t.Errorf("report is not OK: %+v", report)
	}
}

func TestInstallKeepsTheManifest(t *testing.T) {
	manifestDigest := serveManifest(t, "b10783", map[string]string{
		"llama-b10783-bin-ubuntu-cpu-arm64.tar.gz": "abcd",
	})
	recordAssets(t)

	dest := t.TempDir()
	target := Target{Arch: ARM64, OS: Linux, Processor: CPU, Version: "b10783@sha256:" + manifestDigest}
	if err := Install(context.Background(), target, dest, nil, nil); err != nil {
		t.Fatalf("Install() failed: %v", err)
	}

	body, err := ReadInstallManifest(dest)
	if err != nil {
		t.Fatalf("the manifest was not kept: %v", err)
	}
	sum := sha256.Sum256(body)
	if got := hex.EncodeToString(sum[:]); got != manifestDigest {
		t.Errorf("the manifest that was kept has the digest %s, want %s", got, manifestDigest)
	}

	record, err := ReadInstallRecord(dest)
	if err != nil {
		t.Fatal(err)
	}
	if record.ManifestSHA256 != manifestDigest {
		t.Errorf("record ManifestSHA256 = %q, want %s", record.ManifestSHA256, manifestDigest)
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
