//go:build js && wasm

package llamawasm

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"syscall/js"
)

// The llama.cpp module has its own filesystem in memory, and llama.cpp opens
// the model as an ordinary file in it. A program must therefore put the file
// there before ModelLoadFromFile.

// WriteModelFile writes data to name in the filesystem of the llama.cpp
// module. It makes the directories of the path that are not there.
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

// RemoveModelFile removes name from the filesystem of the llama.cpp module,
// which gives the memory of the file back.
func RemoveModelFile(name string) error {
	fs, err := filesystem()
	if err != nil {
		return err
	}
	_, err = fsCall(fs, "unlink", name)
	return err
}

// filesystem gives the FS object of the llama.cpp module.
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

// makeDir makes the directories of the path of name. The filesystem of the
// module starts with almost nothing in it, so a path such as
// /models/model.gguf needs its directory first.
func makeDir(fs js.Value, name string) error {
	dir := path.Dir(name)
	if dir == "." || dir == "/" || dir == "" {
		return nil
	}
	if _, err := fsCall(fs, "mkdirTree", dir); err != nil {
		// A directory that is already there is not a failure.
		if strings.Contains(err.Error(), "EEXIST") {
			return nil
		}
		return err
	}
	return nil
}

// fsCall calls a function of the filesystem of the module and turns a
// JavaScript exception into an error. A call that fails throws, and a throw in
// a call from Go is a panic, which would stop the program.
func fsCall(fs js.Value, name string, args ...any) (result js.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("llamawasm: FS.%s failed: %v", name, r)
		}
	}()

	return fs.Call(name, args...), nil
}

// FetchModelFile gets a model over the network and writes it to name in the
// filesystem of the llama.cpp module. It makes the directories of the path that
// are not there.
//
// The body of the response goes into the file piece by piece, so the memory
// that this needs stays near the size of the model and not twice that size.
// The progress function, if it is not nil, is called with the number of bytes
// so far and the total number of bytes. The total is 0 if the server does not
// give a length.
//
// One JavaScript ArrayBuffer holds at most 2 GB, so a larger model must be in
// splits.
func FetchModelFile(name, url string, progress func(done, total int64)) error {
	fs, err := filesystem()
	if err != nil {
		return err
	}
	if err := makeDir(fs, name); err != nil {
		return err
	}

	response, err := await(js.Global().Call("fetch", url))
	if err != nil {
		return fmt.Errorf("llamawasm: cannot fetch %s: %w", url, err)
	}
	if !response.Get("ok").Bool() {
		return fmt.Errorf("llamawasm: cannot fetch %s: status %d", url, response.Get("status").Int())
	}

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

		// FS.write takes the position in the file, so the whole model never
		// needs to be in memory at one time.
		if _, err := fsCall(fs, "write", stream, value, 0, n, done); err != nil {
			return err
		}

		done += int64(n)
		if progress != nil {
			progress(done, total)
		}
	}

	return nil
}

// ReleaseScratch gives the scratch memory of this package back to the
// llama.cpp module. The next call takes it again, so this is only useful after
// the last inference.
func ReleaseScratch() {
	if !Loaded() {
		return
	}
	for _, s := range []*scratch{&tokenScratch, &textScratch, &pieceScratch, &embdScratch, &errScratch} {
		s.release()
	}
}

// await waits for a JavaScript promise and gives its value.
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
