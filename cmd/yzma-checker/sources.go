package main

// Resolution of the two inputs the audit needs: the yzma tree to analyse and
// the llama.cpp headers to analyse it against. Everything here is designed so
// that a bare `go run .` inside a yzma checkout works with no arguments and no
// assumptions about the surrounding filesystem.

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// yzmaModulePath is the module the checker audits. A directory is accepted as
// a yzma checkout when its go.mod declares this module.
const yzmaModulePath = "github.com/hybridgroup/yzma"

// llamaRepo is the upstream llama.cpp repository headers are fetched from when
// no local checkout is available.
const llamaRepo = "ggml-org/llama.cpp"

// versionURLs are the endpoints yzma's own installer consults to decide which
// llama.cpp build to download, in the order it tries them. Auditing against
// the same ref keeps the checker aligned with the library yzma installs.
var versionURLs = []string{
	"https://hybridgroup.github.io/llama-cpp-builder/version.json",
	"https://hybridgroup.github.io/llama-cpp-builder/previous.json",
}

// headerFiles maps the local basename the parser expects to the path inside
// the llama.cpp tree. These are every header that declares a symbol yzma binds.
var headerFiles = []struct{ local, inTree string }{
	{"llama.h", "include/llama.h"},
	{"ggml.h", "ggml/include/ggml.h"},
	{"ggml-backend.h", "ggml/include/ggml-backend.h"},
	{"ggml-cpu.h", "ggml/include/ggml-cpu.h"},
	{"mtmd.h", "tools/mtmd/mtmd.h"},
	{"mtmd-helper.h", "tools/mtmd/mtmd-helper.h"},
}

// findYzmaRoot locates the yzma tree to audit.
//
// It first walks up from the working directory and from the program's own
// directory, which covers running inside a yzma checkout. Failing that it asks
// the go tool which yzma the surrounding module compiles against, so the
// checker is also usable from a consumer's repo without pointing -yzma at
// anything by hand.
func findYzmaRoot() (string, string, error) {
	starts := []string{mustAbs(".")}
	if exe, err := os.Executable(); err == nil {
		starts = append(starts, filepath.Dir(exe))
	}

	for _, start := range starts {
		for d := start; ; {
			if isYzmaModule(d) {
				return d, "checkout", nil
			}

			parent := filepath.Dir(d)
			if parent == d {
				break
			}

			d = parent
		}
	}

	if dir, err := goListYzma("."); err == nil {
		return dir, "module cache (via go list)", nil
	}

	return "", "", fmt.Errorf("cannot locate a %s checkout; pass -yzma <dir>", yzmaModulePath)
}

func isYzmaModule(dir string) bool {
	b, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return false
	}

	for line := range strings.SplitSeq(string(b), "\n") {
		if strings.TrimSpace(line) == "module "+yzmaModulePath {
			return true
		}
	}

	return false
}

// goListYzma reports the yzma module directory the module rooted at dir
// resolves, which is the tree that actually compiles for that consumer.
func goListYzma(dir string) (string, error) {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", yzmaModulePath)
	cmd.Dir = dir

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	got := strings.TrimSpace(string(out))
	if got == "" {
		return "", fmt.Errorf("go list -m %s returned no directory", yzmaModulePath)
	}

	return got, nil
}

// resolveRef reports the llama.cpp ref to audit against, asking the same
// endpoints yzma's installer uses so the audit tracks the build yzma ships.
func resolveRef() (string, error) {
	var errs []string

	for _, u := range versionURLs {
		tag, err := fetchTag(u)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", u, err))
			continue
		}

		return tag, nil
	}

	return "", fmt.Errorf("cannot determine the current llama.cpp release (%s); pass -ref explicitly",
		strings.Join(errs, "; "))
}

func fetchTag(url string) (string, error) {
	body, err := httpGet(url)
	if err != nil {
		return "", err
	}

	var v struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &v); err != nil {
		return "", err
	}

	if v.TagName == "" {
		return "", fmt.Errorf("no tag_name in response")
	}

	return v.TagName, nil
}

// obtainHeaders materialises the six headers for ref into a directory and
// reports it, along with a human-readable description of where they came from.
//
// Sources are tried in order of how much the caller has pinned down: an
// explicit directory, then a local git checkout, then upstream over HTTP.
func obtainHeaders(hdrDir, llamaDir, ref string) (dir, source string, cleanup func(), err error) {
	noop := func() {}

	if hdrDir != "" {
		if err := checkHeaderDir(hdrDir); err != nil {
			return "", "", noop, err
		}

		return hdrDir, "pre-extracted (-hdrs)", noop, nil
	}

	if llamaDir != "" {
		dir, err := extractHeaders(llamaDir, ref)
		if err != nil {
			return "", "", noop, err
		}

		return dir, fmt.Sprintf("git show %s: in %s", ref, llamaDir), func() { os.RemoveAll(dir) }, nil
	}

	dir, cached, err := downloadHeaders(ref)
	if err != nil {
		return "", "", noop, fmt.Errorf("%w\n"+
			"Pass -llama with a llama.cpp checkout or -hdrs with pre-extracted headers to work offline.", err)
	}

	src := "downloaded from " + llamaRepo + "@" + ref
	if cached {
		src = "cache of " + llamaRepo + "@" + ref
	}

	return dir, src, noop, nil
}

func checkHeaderDir(dir string) error {
	var missing []string

	for _, h := range headerFiles {
		if _, err := os.Stat(filepath.Join(dir, h.local)); err != nil {
			missing = append(missing, h.local)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("%s is missing required header(s): %s", dir, strings.Join(missing, ", "))
	}

	return nil
}

// extractHeaders materialises the headers for ref out of a llama.cpp checkout
// into a scratch dir, using `git show` so the checkout's own worktree state
// (which may be at a different tag) is irrelevant.
func extractHeaders(llamaDir, ref string) (string, error) {
	if _, err := os.Stat(filepath.Join(llamaDir, ".git")); err != nil {
		return "", fmt.Errorf("%s is not a git checkout: %w", llamaDir, err)
	}

	dir, err := os.MkdirTemp("", "yzma-checker-hdrs-")
	if err != nil {
		return "", err
	}

	for _, h := range headerFiles {
		out, err := exec.Command("git", "-C", llamaDir, "show", ref+":"+h.inTree).Output()
		if err != nil {
			os.RemoveAll(dir)
			return "", fmt.Errorf("git show %s:%s in %s: %w", ref, h.inTree, llamaDir, err)
		}

		if err := os.WriteFile(filepath.Join(dir, h.local), out, 0o600); err != nil {
			os.RemoveAll(dir)
			return "", err
		}
	}

	return dir, nil
}

// downloadHeaders fetches the headers for ref from upstream into a persistent
// cache, so only the first run for a given ref needs the network. It reports
// whether the cache already held a complete set.
func downloadHeaders(ref string) (string, bool, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		base = os.TempDir()
	}

	dir := filepath.Join(base, "yzma-checker", "headers", ref)
	if checkHeaderDir(dir) == nil {
		return dir, true, nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", false, err
	}

	for _, h := range headerFiles {
		path := filepath.Join(dir, h.local)
		if _, err := os.Stat(path); err == nil {
			continue
		}

		url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", llamaRepo, ref, h.inTree)

		body, err := httpGet(url)
		if err != nil {
			return "", false, fmt.Errorf("fetch %s: %w", url, err)
		}

		if err := os.WriteFile(path, body, 0o600); err != nil {
			return "", false, err
		}
	}

	return dir, false, nil
}

func httpGet(url string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %s", resp.Status)
	}

	return io.ReadAll(io.LimitReader(resp.Body, 8<<20))
}

func mustAbs(p string) string {
	a, _ := filepath.Abs(p)
	return a
}
