//go:build js && wasm

package llamawasm

import "fmt"

// The calls in this file follow the mtmd library of llama.cpp, which gives a
// model the ability to read images. The names are the names of the mtmd package
// with an Mtmd prefix, because this package also holds the llama calls.
//
//	mtmd.InitFromFile     -> llamawasm.MtmdInitFromFile
//	mtmd.InputChunksInit  -> llamawasm.MtmdInputChunksInit
//	mtmd.Tokenize         -> llamawasm.MtmdTokenize
//	mtmd.HelperEvalChunks -> llamawasm.MtmdHelperEvalChunks
//
// The order of the calls agrees with the examples/vlm program.
//
//	mctx, err := llamawasm.MtmdInitFromFile("mmproj.gguf", model, llamawasm.MtmdContextParamsDefault())
//	bitmap, err := llamawasm.MtmdBitmapInit(width, height, rgb)
//	chunks, err := llamawasm.MtmdInputChunksInit()
//	llamawasm.MtmdTokenize(mctx, chunks, prompt, true, true, []MtmdBitmap{bitmap})
//	nPast, err := llamawasm.MtmdHelperEvalChunks(mctx, ctx, chunks, 0, 0, nBatch, true)
//	// then the usual loop of SamplerSample and Decode, as for text
//
// The prompt must hold one marker for each bitmap. [MtmdMarker] gives the
// marker of the model.

// The handles of the multimodal objects.
type (
	// MtmdContext is a handle to a loaded projector.
	MtmdContext int32

	// MtmdBitmap is a handle to an image in the module.
	MtmdBitmap int32

	// MtmdInputChunks is a handle to a list of chunks of a prompt.
	MtmdInputChunks int32
)

// MtmdContextParams holds what the shim can set while it loads a projector.
//
// It is the small part of mtmd.ContextParamsType that a browser can use.
type MtmdContextParams struct {
	// NThreads is the number of threads for the projector. The projector is the
	// slowest part of an answer, thus this value is more important than the
	// threads of the context.
	NThreads int32

	// UseGPU puts the projector on the GPU if there is one.
	UseGPU bool

	// ImageMinTokens and ImageMaxTokens limit the number of tokens of one
	// image, for a model with a variable resolution. A value of 0 uses the
	// limits of the model. Fewer tokens give less work and less detail.
	//
	// A module of ABI version 3 ignores these values.
	ImageMinTokens int32
	ImageMaxTokens int32
}

// MtmdContextParamsDefault gives the parameters that a projector uses if the
// program changes nothing. These are each thread of the module and the GPU, if
// llama.cpp found one.
func MtmdContextParamsDefault() MtmdContextParams {
	return MtmdContextParams{
		NThreads: Threads(),
		UseGPU:   GPUDevice() != "",
	}
}

// MtmdSupported tells if the module has the multimodal calls. A module from an
// earlier release has none.
func MtmdSupported() bool {
	return has("_yzma_mtmd_init_from_file")
}

// MtmdInitFromFile loads the projector of a multimodal model. Put the file in
// the filesystem of the module, as for the model, and load the model first.
//
// Use [MtmdContextParamsDefault] unless you have a reason for other values.
func MtmdInitFromFile(mmprojPath string, model Model, params MtmdContextParams) (MtmdContext, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}
	if !MtmdSupported() {
		return 0, ErrNoMultimodal
	}

	ptr, err := textScratch.reserve(len(mmprojPath) + 1)
	if err != nil {
		return 0, err
	}
	writeString(ptr, mmprojPath)

	handle, err := callErr("_yzma_mtmd_init_from_file", ptr, int(model),
		int(params.NThreads), boolToInt(params.UseGPU),
		int(params.ImageMinTokens), int(params.ImageMaxTokens))
	if err != nil {
		return 0, err
	}
	return MtmdContext(handle), nil
}

// MtmdFree frees a projector.
func MtmdFree(mctx MtmdContext) error {
	if !Loaded() || !MtmdSupported() {
		return ErrNotLoaded
	}
	callVoid("_yzma_mtmd_free", int(mctx))
	return nil
}

// MtmdSupportVision tells if the projector accepts images.
func MtmdSupportVision(mctx MtmdContext) bool {
	if !Loaded() || !MtmdSupported() {
		return false
	}
	return call("_yzma_mtmd_support_vision", int(mctx)) == 1
}

// MtmdSupportAudio tells if the projector accepts audio. This package cannot
// send audio, thus this call only gives a property of the model.
func MtmdSupportAudio(mctx MtmdContext) bool {
	if !Loaded() || !MtmdSupported() {
		return false
	}
	return call("_yzma_mtmd_support_audio", int(mctx)) == 1
}

// MtmdMarker gives the marker that replaces a media item in the text of a
// prompt, for example "<__media__>".
func MtmdMarker(mctx MtmdContext) string {
	if !Loaded() || !MtmdSupported() {
		return ""
	}

	const size = 64
	ptr, err := pieceScratch.reserve(size)
	if err != nil {
		return ""
	}

	n := call("_yzma_mtmd_get_marker", int(mctx), ptr, size)
	if n <= 0 {
		return ""
	}
	return string(readBytes(ptr, int(n)))
}

// MtmdBitmapInit makes a bitmap from the pixels of an image.
//
// The pixels must be RGB, three bytes for each pixel, with no padding between
// the rows. Thus the length of rgb must be width*height*3. A page reads the RGBA
// of the image from a canvas and removes each fourth byte.
func MtmdBitmapInit(width, height int32, rgb []byte) (MtmdBitmap, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}
	if !MtmdSupported() {
		return 0, ErrNoMultimodal
	}

	want := int(width) * int(height) * 3
	if width <= 0 || height <= 0 || len(rgb) != want {
		return 0, errBitmapSize(width, height, len(rgb), want)
	}

	// The image is large and stays only until the module copies it, thus it
	// does not go in the scratch memory.
	ptr, err := malloc(len(rgb))
	if err != nil {
		return 0, err
	}
	defer free(ptr)

	writeBytes(ptr, rgb)

	handle, err := callErr("_yzma_mtmd_bitmap_init", int(width), int(height), ptr)
	if err != nil {
		return 0, err
	}
	return MtmdBitmap(handle), nil
}

// errBitmapSize gives the necessary length for a bitmap of this size.
func errBitmapSize(width, height int32, got, want int) error {
	return fmt.Errorf("llamawasm: an image of %d by %d needs %d bytes of RGB, got %d",
		width, height, want, got)
}

// MtmdBitmapFree frees a bitmap.
func MtmdBitmapFree(bitmap MtmdBitmap) {
	if !Loaded() || !MtmdSupported() {
		return
	}
	callVoid("_yzma_mtmd_bitmap_free", int(bitmap))
}

// MtmdInputChunksInit makes an empty list of chunks for MtmdTokenize.
func MtmdInputChunksInit() (MtmdInputChunks, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}
	if !MtmdSupported() {
		return 0, ErrNoMultimodal
	}

	handle, err := callErr("_yzma_mtmd_input_chunks_init")
	if err != nil {
		return 0, err
	}
	return MtmdInputChunks(handle), nil
}

// MtmdInputChunksFree frees a list of chunks.
func MtmdInputChunksFree(chunks MtmdInputChunks) {
	if !Loaded() || !MtmdSupported() {
		return
	}
	callVoid("_yzma_mtmd_input_chunks_free", int(chunks))
}

// MtmdInputChunksSize gives the number of chunks in the list.
func MtmdInputChunksSize(chunks MtmdInputChunks) int32 {
	if !Loaded() || !MtmdSupported() {
		return 0
	}
	n := call("_yzma_mtmd_input_chunks_size", int(chunks))
	if n < 0 {
		return 0
	}
	return n
}

// MtmdTokenize puts the text and the bitmaps into the list of chunks. The text
// must hold one marker for each bitmap. [MtmdMarker] gives the marker.
func MtmdTokenize(mctx MtmdContext, chunks MtmdInputChunks, text string,
	addSpecial bool, parseSpecial bool, bitmaps []MtmdBitmap) error {
	if !Loaded() {
		return ErrNotLoaded
	}
	if !MtmdSupported() {
		return ErrNoMultimodal
	}

	textPtr, err := textScratch.reserve(len(text) + 1)
	if err != nil {
		return err
	}
	writeString(textPtr, text)

	// The handles of the bitmaps go in their own memory, because the text is
	// already in the scratch.
	handlesPtr, err := tokenScratch.reserve(max(len(bitmaps), 1) * 4)
	if err != nil {
		return err
	}
	handles := make([]Token, len(bitmaps))
	for i, bitmap := range bitmaps {
		handles[i] = Token(bitmap)
	}
	writeTokens(handlesPtr, handles)

	_, err = callErr("_yzma_mtmd_tokenize", int(mctx), int(chunks), textPtr, len(text),
		boolToInt(addSpecial), boolToInt(parseSpecial), handlesPtr, len(bitmaps))
	return err
}

// MtmdHelperEvalChunks runs each chunk through the model. The text goes as
// tokens and each image goes through the projector first. It gives the position
// after the last chunk, where the generation loop starts.
func MtmdHelperEvalChunks(mctx MtmdContext, ctx Context, chunks MtmdInputChunks,
	nPast Pos, seqID SeqId, nBatch int32, logitsLast bool) (Pos, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}
	if !MtmdSupported() {
		return 0, ErrNoMultimodal
	}

	newNPast, err := callErr("_yzma_mtmd_helper_eval_chunks", int(mctx), int(ctx), int(chunks),
		int(nPast), int(seqID), int(nBatch), boolToInt(logitsLast))
	if err != nil {
		return 0, err
	}
	return Pos(newNPast), nil
}
