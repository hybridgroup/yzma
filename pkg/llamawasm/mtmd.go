//go:build js && wasm

package llamawasm

import "fmt"

// The calls in this file follow the mtmd library of llama.cpp, which is what
// gives a model its eyes. The names are the ones of the mtmd package with an
// Mtmd prefix, because this package holds the llama calls as well:
//
//	mtmd.InitFromFile     -> llamawasm.MtmdInitFromFile
//	mtmd.InputChunksInit  -> llamawasm.MtmdInputChunksInit
//	mtmd.Tokenize         -> llamawasm.MtmdTokenize
//	mtmd.HelperEvalChunks -> llamawasm.MtmdHelperEvalChunks
//
// The order of the calls is the same as in the examples/vlm program:
//
//	mctx, err := llamawasm.MtmdInitFromFile("mmproj.gguf", model, 0, true)
//	bitmap, err := llamawasm.MtmdBitmapInit(width, height, rgb)
//	chunks, err := llamawasm.MtmdInputChunksInit()
//	llamawasm.MtmdTokenize(mctx, chunks, prompt, true, true, []MtmdBitmap{bitmap})
//	nPast, err := llamawasm.MtmdHelperEvalChunks(mctx, ctx, chunks, 0, 0, nBatch, true)
//	// then the same loop of SamplerSample and Decode as for text alone
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

// MtmdSupported tells if the module has the multimodal calls. A module from a
// release before they came has none.
func MtmdSupported() bool {
	return has("_yzma_mtmd_init_from_file")
}

// MtmdInitFromFile loads the projector of a multimodal model. The file must be
// in the filesystem of the module, the same as the model, and the model must be
// loaded first.
//
// useGPU puts the projector on the GPU where there is one. nThreads of 0 leaves
// the number of threads to the module.
func MtmdInitFromFile(mmprojPath string, model Model, nThreads int32, useGPU bool) (MtmdContext, error) {
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

	handle, err := callErr("_yzma_mtmd_init_from_file", ptr, int(model), int(nThreads), boolToInt(useGPU))
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

// MtmdSupportVision tells if the projector takes images.
func MtmdSupportVision(mctx MtmdContext) bool {
	if !Loaded() || !MtmdSupported() {
		return false
	}
	return call("_yzma_mtmd_support_vision", int(mctx)) == 1
}

// MtmdSupportAudio tells if the projector takes audio. This package has no way
// to give it any, so this is here to report what the model is.
func MtmdSupportAudio(mctx MtmdContext) bool {
	if !Loaded() || !MtmdSupported() {
		return false
	}
	return call("_yzma_mtmd_support_audio", int(mctx)) == 1
}

// MtmdMarker gives the marker that stands for a piece of media in the text of a
// prompt, such as "<__media__>".
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
// The pixels must be RGB, three bytes for each one, with no padding between the
// rows, so the length of rgb must be width*height*3. A page gets them from a
// canvas: read the RGBA of the image and drop every fourth byte.
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

	// The image does not go in the scratch memory: it is large, and it stays
	// only until the module copies it.
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

// errBitmapSize says what a bitmap of this size needs.
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

// MtmdInputChunksInit makes an empty list of chunks for MtmdTokenize to fill.
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
// must hold one marker for each bitmap, and [MtmdMarker] gives the marker.
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

	// The handles of the bitmaps go in their own memory, because the text is in
	// the scratch already.
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

// MtmdHelperEvalChunks runs every chunk through the model: the text as tokens,
// and each image through the projector first. It gives the position after the
// last chunk, which the generation loop carries on from.
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
