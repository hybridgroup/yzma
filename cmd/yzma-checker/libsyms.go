package main

// Symbol presence in the installed shared library, which is the one input the
// header-only pipeline structurally cannot see.
//
// Everything else here is verified against llama.cpp's headers, and a header is
// a claim about what the library exports rather than the export itself. What
// purego actually does at load time is dlsym, so a declaration the installed
// build does not export - because the library predates it, was built with the
// feature that defines it disabled, or never had it exported at all - is a
// failure at binding time no header comparison can predict. Nor can the ref this
// audit resolves: it names the build yzma's installer *would* fetch, not the one
// sitting on this machine.
//
// It is opt-in (-lib) on purpose. The installed library is a property of the
// machine rather than of the repository, so a check that ran by default would
// make the report's contents depend on whether somebody had run yzma before, and
// a CI box with no library would either fail or print a section explaining that
// it could not. Nothing here ever produces a violation or a skip either: a
// missing symbol is a fact about a local installation, not a defect in the
// bindings.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

// libShortNames are the libraries yzma dlopens, in pkg/llama/loader.go and
// pkg/mtmd/loader.go. The union of their exports is what a binding's dlsym can
// possibly resolve against.
var libShortNames = []string{"ggml", "ggml-base", "ggml-cpu", "llama", "mtmd"}

// libFilename mirrors pkg/loader.GetLibraryFilename, which is what decides the
// file dlopen is handed.
func libFilename(dir, lib string) string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(dir, lib+".dll")
	case "darwin":
		return filepath.Join(dir, "lib"+lib+".dylib")
	default:
		return filepath.Join(dir, "lib"+lib+".so")
	}
}

// libSymbols reads the exported, defined symbols of every library in dir that
// yzma loads, and reports which files it read them from.
//
// An error means the comparison did not happen - no symbol tooling, or no
// library there - which is a reason to say nothing rather than to fail: see the
// header comment.
func libSymbols(dir string) (map[string]bool, []string, error) {
	syms := map[string]bool{}

	var read []string
	for _, short := range libShortNames {
		path := libFilename(dir, short)
		if _, err := os.Stat(path); err != nil {
			continue
		}

		out, err := nmOutput(path)
		if err != nil {
			return nil, nil, err
		}

		n := 0
		for s := range nmSymbols(out) {
			syms[s] = true
			n++
		}

		read = append(read, fmt.Sprintf("%s (%d symbols)", filepath.Base(path), n))
	}

	if len(read) == 0 {
		return nil, nil, fmt.Errorf("no library yzma loads found in %s", dir)
	}

	return syms, read, nil
}

// nmOutput runs the platform's symbol dumper. `nm -gU` means "external, defined
// only" to both llvm's nm and binutils'; the bare -g run is the fallback for an
// nm too old for -U, whose output nmSymbols filters itself anyway.
func nmOutput(path string) (string, error) {
	out, err := exec.Command("nm", "-gU", path).Output()
	if err == nil {
		return string(out), nil
	}

	out, err2 := exec.Command("nm", "-g", path).Output()
	if err2 != nil {
		return "", fmt.Errorf("nm %s: %w", filepath.Base(path), err)
	}

	return string(out), nil
}

// nmSymbols yields the defined symbol names in nm output.
//
// A line is "<addr> <type> <name>", or "<type> <name>" for an undefined one, and
// the type letter is what separates a symbol the library defines from one it
// merely imports - a U line names a symbol dlsym would not find here either, and
// neither would it find a lowercase (local) one. The leading underscore Mach-O
// adds is stripped so both platforms yield the C name.
func nmSymbols(out string) func(func(string) bool) {
	return func(yield func(string) bool) {
		for line := range strings.SplitSeq(out, "\n") {
			f := strings.Fields(line)
			if len(f) < 2 {
				continue
			}

			kind, name := f[len(f)-2], f[len(f)-1]
			if len(kind) != 1 || kind == "U" || kind != strings.ToUpper(kind) {
				continue
			}

			if runtime.GOOS == "darwin" {
				name = strings.TrimPrefix(name, "_")
			}

			if !yield(name) {
				return
			}
		}
	}
}

// checkLibSymbols reports which bound symbols the installed library does not
// export. The returned note describes what was compared, or why nothing was.
func checkLibSymbols(dir string, bindings []*Binding) (missing []string, note string) {
	syms, read, err := libSymbols(dir)
	if err != nil {
		return nil, fmt.Sprintf("not compared: %v", err)
	}

	for _, b := range bindings {
		if !syms[b.CName] {
			missing = append(missing, fmt.Sprintf("%s (%s)", b.CName, shortPos(b.PrepPos)))
		}
	}

	slices.Sort(missing)

	return slices.Compact(missing), "read " + strings.Join(read, ", ")
}
