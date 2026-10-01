//go:build js && wasm

package llamawasm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"syscall/js"
)

// The shim interface version, which is YZMA_ABI_VERSION in wasm/yzma_wasm.cpp
// of the llama-cpp-builder repo.
//
// This package supports every module version from abiVersionMin to
// abiVersion, so a new yzma works with modules from an older release. A call
// from a later version exists only in modules that have it, so this package
// checks for it before each use.
const (
	abiVersionMin = 1  // 1 has the calls for text generation and embeddings
	abiVersion    = 10 // 2 adds yzma_gpu_device, 3 the multimodal calls, 4 the
	//                    bounds of the tokens of an image, 5 the rest of the
	//                    vocabulary and of the samplers, 6 batches with
	//                    positions and the calls for the memory of a sequence,
	//                    7 the calls that read the logits and the
	//                    embeddings of a batch, 8 yzma_backend_check, 9 the
	//                    metadata, the state, and the performance counters,
	//                    10 contexts with a unified cache and a cap on outputs
)

// Error codes that the shim returns. These match the values in
// wasm/yzma_wasm.cpp.
const (
	errGeneric  = -1
	errHandle   = -2
	errAlloc    = -3
	errLoad     = -4
	errTooSmall = -5
)

var (
	// ErrNotLoaded means Load did not run or failed.
	ErrNotLoaded = errors.New("llamawasm: the llama.cpp module is not loaded, call Load first")

	// ErrNoModule means the JavaScript glue did not run before Load.
	ErrNoModule = errors.New("llamawasm: globalThis.yzmaReady is missing, load yzma-loader.js first")

	// ErrNoMultimodal means the module predates the multimodal calls, which
	// arrived in ABI version 3.
	ErrNoMultimodal = errors.New("llamawasm: this llama.cpp module has no multimodal calls, install a newer build")

	// ErrNoOutputs means the module predates the calls that read the logits
	// and embeddings of a batch, which arrived in ABI version 7.
	ErrNoOutputs = errors.New("llamawasm: this llama.cpp module has no calls for the logits and the embeddings of a batch, install a newer build")

	// ErrNoBackendSampling means the module predates the backend sampling
	// calls, which arrived in ABI version 7.
	ErrNoBackendSampling = errors.New("llamawasm: this llama.cpp module has no calls for the sampling of the backend, install a newer build")

	// ErrNoContextFlags means the module predates the calls that change a
	// context after it is created, which arrived in ABI version 7.
	ErrNoContextFlags = errors.New("llamawasm: this llama.cpp module cannot change a context after it is made, install a newer build")

	// ErrNoKVUnified means the module predates contexts with a unified cache,
	// which arrived in ABI version 10.
	ErrNoKVUnified = errors.New("llamawasm: this llama.cpp module cannot make a context with a unified cache, install a newer build")

	// ErrNoPerf means the module predates the performance counters, which
	// arrived in ABI version 9.
	ErrNoPerf = errors.New("llamawasm: this llama.cpp module has no performance counters, install a newer build")
)

// mod is the Emscripten module instance of llama.cpp.
var mod js.Value

// threaded reports whether the page module uses more than one thread.
var threaded bool

// moduleABI is the interface version of the loaded module.
var moduleABI int

// gpuDevice is the name of the non CPU llama.cpp device, or an empty string if
// there is none. Init sets it.
var gpuDevice string

// backendOK reports whether the llama.cpp device matches the CPU. Init sets it.
var backendOK = true

// Load waits for the llama.cpp WebAssembly module and attaches to it.
//
// The path argument is unused. It keeps the same signature as llama.Load, so
// the same code builds for a native platform and for a browser.
//
// The JavaScript glue must run first. It puts a promise in
// globalThis.yzmaReady, and that promise resolves to the module.
func Load(path string) error {
	global := js.Global()

	ready := global.Get("yzmaReady")
	if ready.IsUndefined() || ready.IsNull() {
		// The glue can set the instance directly without a promise.
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

// attach stores the module and checks that it has the required interface.
func attach(m js.Value) error {
	if m.IsUndefined() || m.IsNull() {
		return ErrNoModule
	}
	if m.Get("_yzma_abi_version").IsUndefined() {
		return errors.New("llamawasm: the module is not a yzma build of llama.cpp")
	}

	mod = m

	v, err := settle(m.Call("_yzma_abi_version"))
	if err != nil || v.Type() != js.TypeNumber {
		mod = js.Undefined()
		return fmt.Errorf("llamawasm: cannot read the ABI version of the llama.cpp module: %v", err)
	}
	got := v.Int()
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

// has reports whether the module has a call. An earlier module does not have
// the calls of a later version.
func has(name string) bool {
	return Loaded() && mod.Get(name).Type() == js.TypeFunction
}

// Loaded reports whether the llama.cpp module is ready to use.
func Loaded() bool {
	return !mod.IsUndefined() && !mod.IsNull()
}

// Threaded reports whether the page module uses more than one thread. A browser
// only allows threads on a page with the Cross-Origin-Opener-Policy and
// Cross-Origin-Embedder-Policy headers.
func Threaded() bool {
	return threaded
}

// Threads returns the number of threads the module can use, which the
// JavaScript glue reads from the machine. The value is 1 for a single thread
// build and for the WebGPU build, where the GPU does the work.
//
// llama.cpp asks for four threads unless a caller changes it. So
// [ContextDefaultParams] and [MtmdContextParamsDefault] send this value.
func Threads() int32 {
	if !Loaded() {
		return 0
	}

	n := js.Global().Get("yzmaThreads")
	if n.Type() != js.TypeNumber || n.Int() < 1 {
		return 1
	}
	return int32(n.Int())
}

// Init starts the llama.cpp backend. Call it after Load.
func Init() {
	if !Loaded() {
		return
	}
	callVoid("_yzma_backend_init")

	// llama.cpp devices only exist after the backend starts. This request
	// makes the WebGPU backend look for an adapter.
	gpuDevice = readGPUDevice()

	backendOK = readBackendCheck()
}

// readBackendCheck asks the shim to compare the non CPU device against the
// CPU. A module before ABI version 8 has no such call, so the answer is that
// nothing is known to be wrong.
func readBackendCheck() bool {
	if !has("_yzma_backend_check") {
		return true
	}
	return call("_yzma_backend_check") == 0
}

// BackendOK reports whether the llama.cpp device matches the CPU. Call it after
// Init.
//
// Some drivers expose a WebGPU adapter that llama.cpp accepts but that
// computes wrong values. A model on such a device answers with random
// vocabulary tokens. Init runs one small matrix multiply on the device and on
// the CPU and compares the results, so a page can warn or fall back to the CPU
// instead of showing nonsense.
//
// It returns true for a CPU only build, and for a module before ABI version 8,
// which has no such test.
func BackendOK() bool {
	return backendOK
}

// readGPUDevice asks the shim for the name of the non CPU device.
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

// GPUDevice returns the name of the non CPU llama.cpp device, or an empty
// string if the computation runs on the CPU. Call it after Init.
//
// A page can ask for WebGPU and still get the CPU, because the backend needs an
// adapter with f16 shaders. This reports what llama.cpp actually uses.
func GPUDevice() string {
	return gpuDevice
}

// Backend returns the name of the compute backend. The values are "webgpu",
// "cpu-threads" for the multithreaded CPU build, and "cpu". Call it
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

// call runs a shim function and returns the result as an int32. It returns
// errBadHandle when the result is not a number, such as after a rejected promise.
func call(name string, args ...any) int32 {
	n, ok := callNumber(name, args...)
	if !ok {
		return errBadHandle
	}
	return int32(n)
}

// callVoid runs a shim function that has no result.
func callVoid(name string, args ...any) {
	settle(mod.Call(name, args...))
}

// callNumber runs a shim function and reports false when the result is not a number.
//
// A module with the WebGPU backend needs the browser GPU, and a GPU request is
// asynchronous. Emscripten builds that module with JSPI, so such a call returns
// a promise and not a number. This function waits for the promise. A CPU build
// returns the number directly, which is the fast path.
func callNumber(name string, args ...any) (float64, bool) {
	v, err := settle(mod.Call(name, args...))
	if err != nil || v.Type() != js.TypeNumber {
		return 0, false
	}
	return v.Float(), true
}

// rejected is the reason of the last rejected promise. shimError reports it once.
var rejected error

// settle waits for a promise. It returns a value that is not a promise, or the
// reason the promise was rejected.
func settle(v js.Value) (js.Value, error) {
	if v.Type() != js.TypeObject || v.Get("then").Type() != js.TypeFunction {
		return v, nil
	}

	resolved, err := await(v)
	if err != nil {
		rejected = err
		return js.Undefined(), err
	}
	return resolved, nil
}

// callErr runs a shim function and turns a negative result into an error with
// the shim error text.
func callErr(name string, args ...any) (int32, error) {
	rc := call(name, args...)
	if rc < 0 {
		return rc, shimError(name, rc)
	}
	return rc, nil
}

// callString runs a call that copies a string into a buffer and reads the
// string. The args come before the pointer and size that the call takes
// last. It uses a larger buffer if the first buffer is too small.
func callString(name string, size int, args ...any) string {
	if !has(name) {
		return ""
	}

	for {
		ptr, err := pieceScratch.reserve(size)
		if err != nil {
			return ""
		}

		n := call(name, append(args, ptr, size)...)
		switch {
		case n == errTooSmall && size < 1<<20:
			size *= 4
			continue
		case n <= 0:
			return ""
		default:
			return string(readBytes(ptr, int(n)))
		}
	}
}

// shimError creates an error from a return code and the last shim error.
func shimError(name string, rc int32) error {
	if err := rejected; err != nil {
		rejected = nil
		return fmt.Errorf("llamawasm: %s: %w", name, err)
	}
	if text := lastError(); text != "" {
		return fmt.Errorf("llamawasm: %s: %s (%d)", name, text, rc)
	}
	return fmt.Errorf("llamawasm: %s failed with code %d", name, rc)
}

// lastError reads the text of the last shim error.
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
// The two WebAssembly modules do not share memory, so each value that goes
// into llama.cpp is a copy. ALLOW_MEMORY_GROWTH creates a new buffer each time
// the memory grows, which detaches the old views. So each read and write gets
// the view again.
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

// writeString writes s and a zero byte at ptr. The space at ptr must be at
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

// writeInt32s writes a slice of any 32-bit type that a batch carries.
func writeInt32s[T ~int32](ptr int, values []T) {
	b := make([]byte, len(values)*4)
	for i, v := range values {
		binary.LittleEndian.PutUint32(b[i*4:], uint32(v))
	}
	writeBytes(ptr, b)
}

func writeInt8s(ptr int, values []int8) {
	b := make([]byte, len(values))
	for i, v := range values {
		b[i] = byte(v)
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

func writeFloats(ptr int, values []float32) {
	b := make([]byte, len(values)*4)
	for i, v := range values {
		binary.LittleEndian.PutUint32(b[i*4:], math.Float32bits(v))
	}
	writeBytes(ptr, b)
}

// writeStrings writes each string followed by a zero byte and returns the
// number of bytes written. The shim reads string lists in this layout, because
// an array of pointers cannot cross the boundary.
func writeStrings(ptr int, values []string) int {
	n := 0
	for _, s := range values {
		n += len(s) + 1
	}

	b := make([]byte, 0, n)
	for _, s := range values {
		b = append(b, s...)
		b = append(b, 0)
	}
	writeBytes(ptr, b)

	return n
}

// stringsSize returns the number of bytes that writeStrings needs.
func stringsSize(values []string) int {
	n := 0
	for _, s := range values {
		n += len(s) + 1
	}
	return n
}

func readFloat64s(ptr, n int) []float64 {
	b := readBytes(ptr, n*8)
	values := make([]float64, n)
	for i := range values {
		values[i] = math.Float64frombits(binary.LittleEndian.Uint64(b[i*8:]))
	}
	return values
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
// memory for one call
//
// A sampler that takes a grammar or a list of strings is created once, outside
// the token generation loop. So each of these allocates and frees rather than
// hold a permanent scratch area.
//

// allocString copies s and a trailing zero byte into the module. The second
// result frees the memory.
func allocString(s string) (int, func(), error) {
	ptr, err := malloc(len(s) + 1)
	if err != nil {
		return 0, func() {}, err
	}
	writeString(ptr, s)
	return ptr, func() { free(ptr) }, nil
}

// allocStrings copies each string and a trailing zero byte into the module,
// which is the layout the shim reads. An empty list returns a pointer of 0,
// which the shim sees as a null pointer.
func allocStrings(values []string) (int, func(), error) {
	if len(values) == 0 {
		return 0, func() {}, nil
	}

	ptr, err := malloc(stringsSize(values))
	if err != nil {
		return 0, func() {}, err
	}
	writeStrings(ptr, values)
	return ptr, func() { free(ptr) }, nil
}

// allocTokens copies the tokens into the module. An empty list returns 0.
func allocTokens(tokens []Token) (int, func(), error) {
	if len(tokens) == 0 {
		return 0, func() {}, nil
	}

	ptr, err := malloc(len(tokens) * 4)
	if err != nil {
		return 0, func() {}, err
	}
	writeTokens(ptr, tokens)
	return ptr, func() { free(ptr) }, nil
}

// allocFloats copies the values into the module. An empty list returns 0.
func allocFloats(values []float32) (int, func(), error) {
	if len(values) == 0 {
		return 0, func() {}, nil
	}

	ptr, err := malloc(len(values) * 4)
	if err != nil {
		return 0, func() {}, err
	}
	writeFloats(ptr, values)
	return ptr, func() { free(ptr) }, nil
}

//
// scratch memory
//
// The generation loop calls into the module many times for each token. A
// permanent scratch area avoids allocations in the module, which also keeps the
// module memory from growing during generation.
//

type scratch struct {
	ptr  int
	size int
}

// reserve returns a pointer to at least n bytes. The contents of an earlier
// reserve on the same scratch are lost.
func (s *scratch) reserve(n int) (int, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}
	if s.size >= n {
		return s.ptr, nil
	}

	// Reserve more space than requested, so a loop that grows by small
	// amounts does not allocate each time.
	size := n * 2
	ptr, err := malloc(size)
	if err != nil {
		return 0, err
	}

	free(s.ptr)
	s.ptr, s.size = ptr, size
	return ptr, nil
}

// release returns the scratch memory to the module.
func (s *scratch) release() {
	free(s.ptr)
	s.ptr, s.size = 0, 0
}

// Each purpose has its own scratch, because several are in use at once.
// tokenScratch holds input tokens, textScratch holds an input string, and
// pieceScratch holds output bytes.
// posScratch, seqScratch, and logitScratch hold the other arrays of a batch,
// which go into the module beside the tokens. outScratch is different and holds
// the logits and embeddings that come out of the module.
var (
	tokenScratch scratch
	textScratch  scratch
	pieceScratch scratch
	embdScratch  scratch
	errScratch   scratch
	posScratch   scratch
	nSeqScratch  scratch
	seqScratch   scratch
	logitScratch scratch
	outScratch   scratch
)
