//go:build !windows

package loader

import "github.com/jupiterrider/ffi"

// Open loads the shared library at filename.
func Open(filename string) (ffi.Lib, error) {
	return ffi.Load(filename)
}

// PreloadBackends does nothing outside of Windows.
func PreloadBackends(path string) func() {
	return func() {}
}
