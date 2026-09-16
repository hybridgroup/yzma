//go:build js && wasm

package llamawasm

import (
	"strings"
	"syscall/js"
	"testing"
)

// The fakes are in JavaScript, because a filesystem and a fetch of the browser
// are objects with many functions.
const fakeSource = `
globalThis.__yzmaFake = (function () {
	const files = {};

	function setFetch(chunks, headers) {
		globalThis.fetch = function () {
			let next = 0;
			return Promise.resolve({
				ok: true,
				status: 200,
				headers: { get: (name) => headers[name.toLowerCase()] ?? null },
				body: {
					getReader: () => ({
						read: () => Promise.resolve(
							next < chunks.length
								? { done: false, value: Uint8Array.from(chunks[next++]) }
								: { done: true }
						),
					}),
				},
			});
		};
	}

	return {
		setFetch,
		size: (name) => (files[name] ? files[name].length : -1),
		module: {
			FS: {
				mkdirTree() {},
				open(name) { files[name] = []; return { name }; },
				close() {},
				write(stream, buf, offset, length, position) {
					for (let i = 0; i < length; i++) {
						files[stream.name][position + i] = buf[offset + i];
					}
					return length;
				},
			},
		},
	};
})();
`

// fake puts the module and the fetch of the test in place and gives the
// helper back. It repairs the globals when the test ends.
func fake(t *testing.T, chunks [][]byte, headers map[string]any) js.Value {
	t.Helper()

	js.Global().Call("eval", fakeSource)
	helper := js.Global().Get("__yzmaFake")

	parts := make([]any, len(chunks))
	for i, chunk := range chunks {
		bytes := make([]any, len(chunk))
		for j, b := range chunk {
			bytes[j] = int(b)
		}
		parts[i] = bytes
	}
	helper.Call("setFetch", parts, headers)

	previous := mod
	mod = helper.Get("module")
	t.Cleanup(func() { mod = previous })

	return helper
}

func TestFetchModelFile(t *testing.T) {
	helper := fake(t, [][]byte{{1, 2, 3}, {4, 5}}, map[string]any{"content-length": "5"})

	if err := FetchModelFile("/models/a.gguf", "http://test/a.gguf", nil); err != nil {
		t.Fatalf("FetchModelFile gave %v, want no error", err)
	}
	if size := helper.Call("size", "/models/a.gguf").Int(); size != 5 {
		t.Errorf("the file has %d bytes, want 5", size)
	}
}

func TestFetchModelFileShortBody(t *testing.T) {
	fake(t, [][]byte{{1, 2, 3}}, map[string]any{"content-length": "5"})

	err := FetchModelFile("/models/b.gguf", "http://test/b.gguf", nil)
	if err == nil {
		t.Fatal("FetchModelFile gave no error, want one for a body that stops early")
	}
	if !strings.Contains(err.Error(), "3 bytes of 5") {
		t.Errorf("the error is %q, want the two counts in it", err)
	}
}

func TestFetchModelFileNoLength(t *testing.T) {
	// A server that gives no length leaves nothing to compare.
	if err := fakeFetch(t, map[string]any{}); err != nil {
		t.Fatalf("FetchModelFile gave %v, want no error", err)
	}
}

func TestFetchModelFileCompressed(t *testing.T) {
	// The length is of the compressed body, thus the counts do not agree and
	// the check has to stay away.
	headers := map[string]any{"content-length": "2", "content-encoding": "gzip"}
	if err := fakeFetch(t, headers); err != nil {
		t.Fatalf("FetchModelFile gave %v, want no error", err)
	}
}

// fakeFetch runs one download of three bytes with the headers.
func fakeFetch(t *testing.T, headers map[string]any) error {
	t.Helper()

	fake(t, [][]byte{{1, 2, 3}}, headers)
	return FetchModelFile("/models/c.gguf", "http://test/c.gguf", nil)
}
