//go:build js && wasm

package llamawasm

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"syscall/js"
)

// The llama.cpp module has its own in memory filesystem and opens the model as
// a normal file in it. So a program must put the file there before
// ModelLoadFromFile.

// WriteModelFile writes data to name in the llama.cpp module filesystem. It
// creates any missing directories in the path.
func WriteModelFile(name string, data []byte) error {
	fs, err := filesystem()
	if err != nil {
		return err
	}
	if err := makeDir(fs, name); err != nil {
		return err
	}

	buf := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(buf, data)

	_, err = fsCall(fs, "writeFile", name, buf)
	return err
}

// RemoveModelFile removes name from the llama.cpp module filesystem and frees
// the file memory.
func RemoveModelFile(name string) error {
	fs, err := filesystem()
	if err != nil {
		return err
	}
	_, err = fsCall(fs, "unlink", name)
	return err
}

// filesystem returns the FS object of the llama.cpp module.
func filesystem() (js.Value, error) {
	if !Loaded() {
		return js.Undefined(), ErrNotLoaded
	}

	fs := mod.Get("FS")
	if fs.IsUndefined() || fs.IsNull() {
		return js.Undefined(), errors.New("llamawasm: the llama.cpp module has no filesystem")
	}
	return fs, nil
}

// makeDir creates the directories in the path of name. The module filesystem
// starts almost empty, so a path such as /models/model.gguf needs its
// directory first.
func makeDir(fs js.Value, name string) error {
	dir := path.Dir(name)
	if dir == "." || dir == "/" || dir == "" {
		return nil
	}
	if _, err := fsCall(fs, "mkdirTree", dir); err != nil {
		// A directory that already exists is not a failure.
		if strings.Contains(err.Error(), "EEXIST") {
			return nil
		}
		return err
	}
	return nil
}

// fsCall calls a module filesystem function and turns a JavaScript exception
// into an error. A failed call throws, and a throw in a call from Go is a panic
// that stops the program.
func fsCall(fs js.Value, name string, args ...any) (result js.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("llamawasm: FS.%s failed: %v", name, r)
		}
	}()

	return fs.Call(name, args...), nil
}

// FetchModelFile downloads a model and writes it to name in the llama.cpp
// module filesystem. It creates any missing directories in the path.
//
// The response body goes into the file in small chunks. So memory use stays
// near the model size and not twice that.
//
// If progress is not nil, it receives the bytes read so far and the total
// bytes. The total is 0 if the server sends no length.
//
// A body that ends before the length the server sent is an error, because
// llama.cpp can load an incomplete model file.
//
// One JavaScript ArrayBuffer holds at most 2 GB, so a larger model must be
// split.
func FetchModelFile(name, url string, progress func(done, total int64)) error {
	fs, err := filesystem()
	if err != nil {
		return err
	}
	if err := makeDir(fs, name); err != nil {
		return err
	}

	// Firefox writes the response to its cache while the program reads the
	// stream. A model is larger than the maximum cache entry size, so the
	// browser aborts the stream with an error. no-store keeps the model out of
	// the cache.
	options := js.Global().Get("Object").New()
	options.Set("cache", "no-store")

	response, err := await(js.Global().Call("fetch", url, options))
	if err != nil {
		return fmt.Errorf("llamawasm: cannot fetch %s: %w", url, err)
	}
	if !response.Get("ok").Bool() {
		return fmt.Errorf("llamawasm: cannot fetch %s: status %d", url, response.Get("status").Int())
	}

	// A compressed body has a content-length of the compressed size, which is
	// not the number of bytes that arrive here.
	encoded := response.Get("headers").Call("get", "content-encoding").Truthy()

	var total int64
	if length := response.Get("headers").Call("get", "content-length"); length.Truthy() {
		total = int64(js.Global().Get("parseInt").Invoke(length, 10).Int())
	}

	body := response.Get("body")
	if body.IsUndefined() || body.IsNull() {
		return errors.New("llamawasm: the response has no body to read")
	}

	stream, err := fsCall(fs, "open", name, "w")
	if err != nil {
		return err
	}
	defer func() { _, _ = fsCall(fs, "close", stream) }()

	reader := body.Call("getReader")

	var done int64
	for {
		chunk, err := await(reader.Call("read"))
		if err != nil {
			return fmt.Errorf("llamawasm: cannot read %s: %w", url, err)
		}
		if chunk.Get("done").Bool() {
			break
		}

		value := chunk.Get("value")
		n := value.Get("length").Int()

		// FS.write takes the file position, so the full model is never in
		// memory at once.
		if _, err := fsCall(fs, "write", stream, value, 0, n, done); err != nil {
			return err
		}

		done += int64(n)
		if progress != nil {
			progress(done, total)
		}
	}

	// A body that ends early leaves an incomplete file. llama.cpp can load
	// such a file and then compute wrong values, so fail here.
	if total > 0 && !encoded && done != total {
		return fmt.Errorf("llamawasm: %s gave %d bytes of %d", url, done, total)
	}

	return nil
}

// ReleaseScratch returns this package's scratch memory to the llama.cpp
// module. The next call allocates it again, so use this only after the last
// inference.
func ReleaseScratch() {
	if !Loaded() {
		return
	}
	for _, s := range []*scratch{
		&tokenScratch, &textScratch, &pieceScratch, &embdScratch, &errScratch,
		&posScratch, &nSeqScratch, &seqScratch, &logitScratch, &outScratch,
	} {
		s.release()
	}
}

// await waits for a JavaScript promise and returns its value.
func await(promise js.Value) (js.Value, error) {
	type result struct {
		value js.Value
		err   error
	}
	done := make(chan result, 1)

	onOK := js.FuncOf(func(this js.Value, args []js.Value) any {
		var v js.Value
		if len(args) > 0 {
			v = args[0]
		}
		done <- result{value: v}
		return nil
	})
	defer onOK.Release()

	onErr := js.FuncOf(func(this js.Value, args []js.Value) any {
		msg := "unknown error"
		if len(args) > 0 {
			msg = args[0].Call("toString").String()
		}
		done <- result{err: errors.New(msg)}
		return nil
	})
	defer onErr.Release()

	promise.Call("then", onOK, onErr)

	r := <-done
	return r.value, r.err
}
