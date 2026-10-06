//go:build windows

package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jupiterrider/ffi"
	"golang.org/x/sys/windows"
)

// Open loads the shared library at filename.
// Windows searches the folder of the library for its dependencies first.
func Open(filename string) (ffi.Lib, error) {
	abs, err := filepath.Abs(filename)
	if err != nil {
		return ffi.Lib{}, fmt.Errorf("%s: error loading library: %w", filename, err)
	}

	h, err := windows.LoadLibraryEx(abs, 0, windows.LOAD_WITH_ALTERED_SEARCH_PATH)
	if err != nil {
		return ffi.Lib{}, fmt.Errorf("%s: error loading library: %w", filename, err)
	}

	return ffi.Lib{Addr: uintptr(h)}, nil
}

// PreloadBackends loads the ggml backend libraries in path so that their dependencies
// come from path. Call the returned func after ggml has loaded the backends.
func PreloadBackends(path string) func() {
	entries, err := os.ReadDir(path)
	if err != nil {
		return func() {}
	}

	var libs []ffi.Lib
	for _, e := range entries {
		if match, _ := filepath.Match("ggml-*.dll", e.Name()); !match {
			continue
		}

		if l, err := Open(filepath.Join(path, e.Name())); err == nil {
			libs = append(libs, l)
		}
	}

	return func() {
		for _, l := range libs {
			l.Close()
		}
	}
}
