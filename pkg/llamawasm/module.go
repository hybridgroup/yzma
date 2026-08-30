//go:build js && wasm

package llamawasm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"syscall/js"
)

// The version of the interface of the shim, which is YZMA_ABI_VERSION in
// wasm/yzma_wasm.cpp of the llama-cpp-builder repo.
//
// This package drives a module of any version from abiVersionMin to
// abiVersion, so a new yzma still works with the modules of an older release.
// A call that came with a later version is there only if the module has it, so
// each one is tested for before it is used.
const (
	abiVersionMin = 1 // 1 has the calls for text generation and embeddings
	abiVersion    = 3 // 2 adds yzma_gpu_device, 3 adds the multimodal calls
)

// Error codes that the shim returns. These are the same as the values in
// wasm/yzma_wasm.cpp.
const (
	errGeneric  = -1
	errHandle   = -2
	errAlloc    = -3
	errLoad     = -4
	errTooSmall = -5
)

var (
	// ErrNotLoaded says that Load did not run, or that it failed.
	ErrNotLoaded = errors.New("llamawasm: the llama.cpp module is not loaded, call Load first")

	// ErrNoModule says that the JavaScript glue did not run before Load.
	ErrNoModule = errors.New("llamawasm: globalThis.yzmaReady is missing, load yzma-loader.js first")

	// ErrNoMultimodal says that the module came from a release before the
	// multimodal calls, which are in ABI version 3 and later.
	ErrNoMultimodal = errors.New("llamawasm: this llama.cpp module has no multimodal calls, install a newer build")
)

// mod is the Emscripten module instance of llama.cpp.
var mod js.Value

// threaded tells if the module that the page took uses more than one thread.
var threaded bool

// moduleABI is the version of the interface of the module that is loaded.
var moduleABI int

// gpuDevice is the name of the device of llama.cpp that is not the CPU, or an
// empty string if there is none. Init fills it in.
var gpuDevice string

// Load waits for the llama.cpp WebAssembly module and attaches to it.
//
// The path argument is not used. It is here to keep the same shape as
// llama.Load, so that the same code can build for a native platform and for a
// browser.
//
// The JavaScript glue must run first. It puts a promise in
// globalThis.yzmaReady, and the value of that promise is the module.
func Load(path string) error {
	global := js.Global()

	ready := global.Get("yzmaReady")
	if ready.IsUndefined() || ready.IsNull() {
		// The glue may have put the instance in place without a promise.
		if m := global.Get("yzmaModule"); !m.IsUndefined() && !m.IsNull() {
			return attach(m)
		}
		return ErrNoModule
	}

	if ready.Type() != js.TypeObject || ready.Get("then").IsUndefined() {
		return attach(ready)
	}

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
		done <- result{err: fmt.Errorf("llamawasm: the llama.cpp module failed to load: %s", msg)}
		return nil
	})
	defer onErr.Release()

	ready.Call("then", onOK, onErr)

	r := <-done
	if r.err != nil {
		return r.err
	}
	return attach(r.value)
}

// attach keeps the module and makes sure that it has the interface that this
// package needs.
func attach(m js.Value) error {
	if m.IsUndefined() || m.IsNull() {
		return ErrNoModule
	}
	if m.Get("_yzma_abi_version").IsUndefined() {
		return errors.New("llamawasm: the module is not a yzma build of llama.cpp")
	}

	mod = m

	got := int(settle(m.Call("_yzma_abi_version")).Int())
	if got < abiVersionMin || got > abiVersion {
		mod = js.Undefined()
		return fmt.Errorf("llamawasm: the llama.cpp module has ABI version %d, this build of yzma drives %d to %d",
			got, abiVersionMin, abiVersion)
	}
	moduleABI = got

	threaded = js.Global().Get("yzmaThreaded").Truthy()
	gpuDevice = ""

	return nil
}

// has tells if the module has a call. A module of an earlier version does not
// have the calls that came later.
func has(name string) bool {
	return Loaded() && mod.Get(name).Type() == js.TypeFunction
}

// Loaded tells if the llama.cpp module is ready to use.
func Loaded() bool {
	return !mod.IsUndefined() && !mod.IsNull()
}

// Threaded tells if the module that the page took uses more than one thread. A
// browser gives more than one thread only to a page that sets the
// Cross-Origin-Opener-Policy and Cross-Origin-Embedder-Policy headers.
func Threaded() bool {
	return threaded
}

// Init starts the llama.cpp backend. Call it after Load.
func Init() {
	if !Loaded() {
		return
	}
	callVoid("_yzma_backend_init")

	// The devices of llama.cpp exist only after the backend starts, and asking
	// for them is what makes the WebGPU backend look for an adapter.
	gpuDevice = readGPUDevice()
}

// readGPUDevice asks the shim for the name of the device that is not the CPU.
func readGPUDevice() string {
	if !has("_yzma_gpu_device") {
		return ""
	}

	const size = 256
	ptr, err := pieceScratch.reserve(size)
	if err != nil {
		return ""
	}

	n := call("_yzma_gpu_device", ptr, size)
	if n <= 0 {
		return ""
	}
	return string(readBytes(ptr, int(n)))
}

// GPUDevice gives the name of the device of llama.cpp that is not the CPU, or
// an empty string if the computation is on the CPU. Call it after Init.
//
// A page can ask for WebGPU and still land on the CPU, because the backend
// needs an adapter that supports f16 shaders. This says what llama.cpp really
// has.
func GPUDevice() string {
	return gpuDevice
}

// Backend gives the name of what does the computation: "webgpu", "cpu-threads"
// for the build on the CPU that uses more than one thread, or "cpu". Call it
// after Init.
func Backend() string {
	switch {
	case gpuDevice != "":
		return "webgpu"
	case threaded:
		return "cpu-threads"
	default:
		return "cpu"
	}
}

// Close stops the llama.cpp backend.
func Close() {
	BackendFree()
}

// BackendInit starts the llama.cpp backend.
func BackendInit() {
	Init()
}

// BackendFree stops the llama.cpp backend.
func BackendFree() {
	if !Loaded() {
		return
	}
	mod.Call("_yzma_backend_free")
}

//
// calls into the module
//

// call runs a function of the shim and returns the result as an int32.
func call(name string, args ...any) int32 {
	return int32(callValue(name, args...).Int())
}

// callVoid runs a function of the shim that has no result.
func callVoid(name string, args ...any) {
	callValue(name, args...)
}

// callValue runs a function of the shim and gives the result.
//
// A module with the WebGPU backend needs the GPU of the browser, and the calls
// that ask for a GPU are asynchronous. Emscripten builds that module with JSPI,
// so such a call gives a promise instead of a number. This waits for the
// promise. A build for the CPU gives the number itself, which stays the fast
// path.
func callValue(name string, args ...any) js.Value {
	return settle(mod.Call(name, args...))
}

// settle waits for a promise, or gives back a value that is not one.
func settle(v js.Value) js.Value {
	if v.Type() != js.TypeObject || v.Get("then").Type() != js.TypeFunction {
		return v
	}

	resolved, err := await(v)
	if err != nil {
		// The shim has no way to report this, so leave the value undefined and
		// let the caller see a result that is not a number.
		return js.Undefined()
	}
	return resolved
}

// callErr runs a function of the shim and turns a negative result into an
// error that holds the text from the shim.
func callErr(name string, args ...any) (int32, error) {
	rc := call(name, args...)
	if rc < 0 {
		return rc, shimError(name, rc)
	}
	return rc, nil
}

// shimError makes an error from a return code and the last error of the shim.
func shimError(name string, rc int32) error {
	if text := lastError(); text != "" {
		return fmt.Errorf("llamawasm: %s: %s (%d)", name, text, rc)
	}
	return fmt.Errorf("llamawasm: %s failed with code %d", name, rc)
}

// lastError reads the text of the last error of the shim.
func lastError() string {
	const size = 512
	ptr, err := errScratch.reserve(size)
	if err != nil {
		return ""
	}
	n := call("_yzma_last_error", ptr, size)
	if n <= 0 {
		return ""
	}
	return string(readBytes(ptr, int(n)))
}

//
// memory of the module
//
// The two WebAssembly modules do not share memory, so every value that goes
// into llama.cpp is a copy. ALLOW_MEMORY_GROWTH makes a new buffer each time
// the memory of the module grows, and that leaves the old views detached, so
// each read and write takes the view again.
//

func heapU8() js.Value {
	return mod.Get("HEAPU8")
}

func view(ptr, n int) js.Value {
	return heapU8().Call("subarray", ptr, ptr+n)
}

func malloc(n int) (int, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}
	ptr := mod.Call("_malloc", n).Int()
	if ptr == 0 {
		return 0, fmt.Errorf("llamawasm: cannot allocate %d bytes in the llama.cpp module", n)
	}
	return ptr, nil
}

func free(ptr int) {
	if ptr != 0 && Loaded() {
		mod.Call("_free", ptr)
	}
}

func writeBytes(ptr int, b []byte) {
	if len(b) == 0 {
		return
	}
	js.CopyBytesToJS(view(ptr, len(b)), b)
}

func readBytes(ptr, n int) []byte {
	b := make([]byte, n)
	if n > 0 {
		js.CopyBytesToGo(b, view(ptr, n))
	}
	return b
}

// writeString writes s and a zero byte at ptr. The room at ptr must be at
// least len(s)+1 bytes.
func writeString(ptr int, s string) {
	b := make([]byte, len(s)+1)
	copy(b, s)
	writeBytes(ptr, b)
}

func writeTokens(ptr int, tokens []Token) {
	b := make([]byte, len(tokens)*4)
	for i, t := range tokens {
		binary.LittleEndian.PutUint32(b[i*4:], uint32(t))
	}
	writeBytes(ptr, b)
}

func readTokens(ptr, n int) []Token {
	b := readBytes(ptr, n*4)
	tokens := make([]Token, n)
	for i := range tokens {
		tokens[i] = Token(int32(binary.LittleEndian.Uint32(b[i*4:])))
	}
	return tokens
}

func readFloats(ptr, n int) []float32 {
	b := readBytes(ptr, n*4)
	values := make([]float32, n)
	for i := range values {
		values[i] = math.Float32frombits(binary.LittleEndian.Uint32(b[i*4:]))
	}
	return values
}

//
// scratch memory
//
// The loop that makes tokens calls into the module many times for each token.
// A scratch area that stays keeps the loop from allocating in the module, which
// also keeps the memory of the module from growing during generation.
//

type scratch struct {
	ptr  int
	size int
}

// reserve gives a pointer to at least n bytes. The contents of an earlier
// reserve on the same scratch are lost.
func (s *scratch) reserve(n int) (int, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}
	if s.size >= n {
		return s.ptr, nil
	}

	// Take room for more than the request, so that a loop that grows a little
	// each time does not allocate each time.
	size := n * 2
	ptr, err := malloc(size)
	if err != nil {
		return 0, err
	}

	free(s.ptr)
	s.ptr, s.size = ptr, size
	return ptr, nil
}

// release gives the memory of the scratch back to the module.
func (s *scratch) release() {
	free(s.ptr)
	s.ptr, s.size = 0, 0
}

// Each purpose has its own scratch, because more than one is in use at the
// same time. tokenScratch holds tokens on the way in, textScratch holds a
// string on the way in, and pieceScratch holds bytes on the way out.
var (
	tokenScratch scratch
	textScratch  scratch
	pieceScratch scratch
	embdScratch  scratch
	errScratch   scratch
)
