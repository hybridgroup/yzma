//go:build js && wasm

package llamawasm

import (
	"encoding/binary"
	"fmt"
)

// The calls in this file arrived in ABI version 13. They read the parts of a
// chunk list, make audio bitmaps, and run the projector one chunk at a time.
//
// A chunk from [MtmdInputChunksGet] belongs to its list, and image tokens
// belong to their chunk. Freeing the owner makes those handles invalid.

type (
	// MtmdInputChunk is a handle to one chunk of a prompt.
	MtmdInputChunk int32

	// MtmdImageTokens is a handle to the image tokens of a chunk.
	MtmdImageTokens int32

	// MtmdBatch is a handle to a batch of media chunks for the projector.
	MtmdBatch int32
)

// MtmdInputChunkType is the type of a chunk.
type MtmdInputChunkType int32

const (
	MtmdInputChunkTypeText MtmdInputChunkType = iota
	MtmdInputChunkTypeImage
	MtmdInputChunkTypeAudio
)

// MtmdDecoderPos is the M-RoPE position of one image token.
type MtmdDecoderPos struct {
	T uint32
	X uint32
	Y uint32
	Z uint32
}

// MtmdBatchAddResult is the result of MtmdBatchAddChunk.
type MtmdBatchAddResult int32

const (
	MtmdBatchAddSuccess  MtmdBatchAddResult = 0
	MtmdBatchAddError    MtmdBatchAddResult = 1
	MtmdBatchAddTooLarge MtmdBatchAddResult = 2
)

// MtmdPartsSupported reports whether the module has the calls in this file.
func MtmdPartsSupported() bool {
	return has("_yzma_mtmd_input_chunks_get")
}

func mtmdPartsReady() error {
	if !Loaded() {
		return ErrNotLoaded
	}
	if !MtmdPartsSupported() {
		return ErrNoMtmdParts
	}
	return nil
}

// MtmdDefaultMarker returns the marker of mtmd for a projector that sets none.
func MtmdDefaultMarker() string {
	return callString("_yzma_mtmd_default_marker", 64)
}

// MtmdDecodeUseMRope reports whether the model uses M-RoPE positions.
func MtmdDecodeUseMRope(mctx MtmdContext) bool {
	return MtmdPartsSupported() && call("_yzma_mtmd_decode_use_mrope", int(mctx)) == 1
}

// MtmdDecodeUseNonCausal reports whether the chunk needs a non-causal mask. A
// chunk of 0 stands for an image chunk.
func MtmdDecodeUseNonCausal(mctx MtmdContext, chunk MtmdInputChunk) bool {
	return MtmdPartsSupported() && call("_yzma_mtmd_decode_use_non_causal", int(mctx), int(chunk)) == 1
}

// MtmdGetAudioSampleRate returns the sample rate of the audio encoder, or -1
// if the projector has no audio.
func MtmdGetAudioSampleRate(mctx MtmdContext) int {
	if !MtmdPartsSupported() {
		return -1
	}
	n := call("_yzma_mtmd_get_audio_sample_rate", int(mctx))
	if n < 0 {
		return -1
	}
	return int(n)
}

// MtmdBitmapInitFromAudio creates a bitmap from mono PCM F32 samples. They
// must be at the rate of [MtmdGetAudioSampleRate]. A page decodes and
// resamples them with an OfflineAudioContext.
func MtmdBitmapInitFromAudio(samples []float32) (MtmdBitmap, error) {
	if err := mtmdPartsReady(); err != nil {
		return 0, err
	}
	if len(samples) == 0 {
		return 0, fmt.Errorf("llamawasm: an audio bitmap needs samples")
	}

	ptr, done, err := allocFloats(samples)
	if err != nil {
		return 0, err
	}
	defer done()

	handle, err := callErr("_yzma_mtmd_bitmap_init_from_audio", len(samples), ptr)
	if err != nil {
		return 0, err
	}
	return MtmdBitmap(handle), nil
}

// MtmdBitmapGetNx returns the width of an image, or the sample count of audio.
func MtmdBitmapGetNx(bitmap MtmdBitmap) uint32 {
	return uint32(nonNegative(partsCall("_yzma_mtmd_bitmap_get_nx", int(bitmap))))
}

// MtmdBitmapGetNy returns the height of an image, or 1 for audio.
func MtmdBitmapGetNy(bitmap MtmdBitmap) uint32 {
	return uint32(nonNegative(partsCall("_yzma_mtmd_bitmap_get_ny", int(bitmap))))
}

// MtmdBitmapGetNBytes returns the size of the bitmap data.
func MtmdBitmapGetNBytes(bitmap MtmdBitmap) uint64 {
	return uint64(nonNegative(partsCall("_yzma_mtmd_bitmap_get_n_bytes", int(bitmap))))
}

// MtmdBitmapGetData returns a copy of the bitmap data.
func MtmdBitmapGetData(bitmap MtmdBitmap) []byte {
	n := int(MtmdBitmapGetNBytes(bitmap))
	if n == 0 {
		return nil
	}

	ptr, err := malloc(n)
	if err != nil {
		return nil
	}
	defer free(ptr)

	got := call("_yzma_mtmd_bitmap_get_data", int(bitmap), ptr, n)
	if got <= 0 {
		return nil
	}
	return readBytes(ptr, int(got))
}

// MtmdBitmapIsAudio reports whether the bitmap holds audio.
func MtmdBitmapIsAudio(bitmap MtmdBitmap) bool {
	return partsCall("_yzma_mtmd_bitmap_is_audio", int(bitmap)) == 1
}

// MtmdBitmapGetId returns the id of a bitmap, which the KV cache uses.
func MtmdBitmapGetId(bitmap MtmdBitmap) string {
	return callString("_yzma_mtmd_bitmap_get_id", 64, int(bitmap))
}

// MtmdBitmapSetId sets the id of a bitmap.
func MtmdBitmapSetId(bitmap MtmdBitmap, id string) error {
	if err := mtmdPartsReady(); err != nil {
		return err
	}
	ptr, done, err := allocString(id)
	if err != nil {
		return err
	}
	defer done()

	_, err = callErr("_yzma_mtmd_bitmap_set_id", int(bitmap), ptr)
	return err
}

// MtmdInputChunksGet returns chunk idx of the list, or 0 if idx is out of
// range.
func MtmdInputChunksGet(chunks MtmdInputChunks, idx uint64) MtmdInputChunk {
	return MtmdInputChunk(nonNegative(partsCall("_yzma_mtmd_input_chunks_get", int(chunks), int(idx))))
}

// MtmdInputChunkGetType returns the type of a chunk.
func MtmdInputChunkGetType(chunk MtmdInputChunk) MtmdInputChunkType {
	return MtmdInputChunkType(nonNegative(partsCall("_yzma_mtmd_input_chunk_get_type", int(chunk))))
}

// MtmdInputChunkGetTokensText returns the tokens of a text chunk, or nil for
// a media chunk.
func MtmdInputChunkGetTokensText(chunk MtmdInputChunk) []Token {
	n := int(MtmdInputChunkGetNTokens(chunk))
	if n == 0 || MtmdInputChunkGetType(chunk) != MtmdInputChunkTypeText {
		return nil
	}

	ptr, err := tokenScratch.reserve(n * 4)
	if err != nil {
		return nil
	}
	got := call("_yzma_mtmd_input_chunk_get_tokens_text", int(chunk), ptr, n)
	if got <= 0 {
		return nil
	}
	return readTokens(ptr, int(got))
}

// MtmdInputChunkGetNTokens returns the token count of a chunk.
func MtmdInputChunkGetNTokens(chunk MtmdInputChunk) uint64 {
	return uint64(nonNegative(partsCall("_yzma_mtmd_input_chunk_get_n_tokens", int(chunk))))
}

// MtmdInputChunkGetId returns the id of a media chunk.
func MtmdInputChunkGetId(chunk MtmdInputChunk) string {
	return callString("_yzma_mtmd_input_chunk_get_id", 64, int(chunk))
}

// MtmdInputChunkGetNPos returns the number of positions that the chunk takes.
func MtmdInputChunkGetNPos(chunk MtmdInputChunk) Pos {
	return Pos(nonNegative(partsCall("_yzma_mtmd_input_chunk_get_n_pos", int(chunk))))
}

// MtmdInputChunkCopy returns a copy of a chunk that the caller owns. Free it
// with [MtmdInputChunkFree].
func MtmdInputChunkCopy(chunk MtmdInputChunk) MtmdInputChunk {
	return MtmdInputChunk(nonNegative(partsCall("_yzma_mtmd_input_chunk_copy", int(chunk))))
}

// MtmdInputChunkFree frees a copy. For a chunk of a list it only releases the
// handle.
func MtmdInputChunkFree(chunk MtmdInputChunk) {
	if MtmdPartsSupported() {
		callVoid("_yzma_mtmd_input_chunk_free", int(chunk))
	}
}

// MtmdInputChunkGetTokensImage returns the image tokens of a media chunk, or 0
// for a text chunk.
func MtmdInputChunkGetTokensImage(chunk MtmdInputChunk) MtmdImageTokens {
	return MtmdImageTokens(nonNegative(partsCall("_yzma_mtmd_input_chunk_get_tokens_image", int(chunk))))
}

// MtmdImageTokensGetNTokens returns the token count of an image.
func MtmdImageTokensGetNTokens(it MtmdImageTokens) uint64 {
	return uint64(nonNegative(partsCall("_yzma_mtmd_image_tokens_get_n_tokens", int(it))))
}

// MtmdImageTokensGetNX returns the width of an image in tokens.
//
// Deprecated: use [MtmdImageTokensGetDecoderPos].
func MtmdImageTokensGetNX(it MtmdImageTokens) uint64 {
	return uint64(nonNegative(partsCall("_yzma_mtmd_image_tokens_get_nx", int(it))))
}

// MtmdImageTokensGetNY returns the height of an image in tokens.
//
// Deprecated: use [MtmdImageTokensGetDecoderPos].
func MtmdImageTokensGetNY(it MtmdImageTokens) uint64 {
	return uint64(nonNegative(partsCall("_yzma_mtmd_image_tokens_get_ny", int(it))))
}

// MtmdImageTokensGetId returns the id of an image.
func MtmdImageTokensGetId(it MtmdImageTokens) string {
	return callString("_yzma_mtmd_image_tokens_get_id", 64, int(it))
}

// MtmdImageTokensGetNPos returns the number of positions that an image takes.
func MtmdImageTokensGetNPos(it MtmdImageTokens) Pos {
	return Pos(nonNegative(partsCall("_yzma_mtmd_image_tokens_get_n_pos", int(it))))
}

// MtmdImageTokensGetDecoderPos returns the M-RoPE position of token i of an
// image that starts at pos0.
func MtmdImageTokensGetDecoderPos(it MtmdImageTokens, pos0 Pos, i uint64) MtmdDecoderPos {
	if !MtmdPartsSupported() {
		return MtmdDecoderPos{}
	}
	ptr, err := outScratch.reserve(16)
	if err != nil {
		return MtmdDecoderPos{}
	}
	if call("_yzma_mtmd_image_tokens_get_decoder_pos", int(it), int(pos0), int(i), ptr) != 4 {
		return MtmdDecoderPos{}
	}
	b := readBytes(ptr, 16)
	return MtmdDecoderPos{
		T: binary.LittleEndian.Uint32(b[0:]),
		X: binary.LittleEndian.Uint32(b[4:]),
		Y: binary.LittleEndian.Uint32(b[8:]),
		Z: binary.LittleEndian.Uint32(b[12:]),
	}
}

// MtmdEncodeChunk runs one media chunk through the projector. Read the result
// with [MtmdGetOutputEmbd].
func MtmdEncodeChunk(mctx MtmdContext, chunk MtmdInputChunk) error {
	if err := mtmdPartsReady(); err != nil {
		return err
	}
	_, err := callErr("_yzma_mtmd_encode_chunk", int(mctx), int(chunk))
	return err
}

// MtmdEncode runs image tokens through the projector.
//
// Deprecated: use [MtmdEncodeChunk].
func MtmdEncode(mctx MtmdContext, it MtmdImageTokens) error {
	if err := mtmdPartsReady(); err != nil {
		return err
	}
	_, err := callErr("_yzma_mtmd_encode", int(mctx), int(it))
	return err
}

// MtmdGetOutputEmbd returns a copy of the output of the last encode. The
// embedSize is ModelNEmbdInp(model) * MtmdInputChunkGetNTokens(chunk).
func MtmdGetOutputEmbd(mctx MtmdContext, embedSize int32) ([]float32, error) {
	if err := mtmdPartsReady(); err != nil {
		return nil, err
	}
	return readEmbd(embedSize, func(ptr int) (int32, error) {
		return callErr("_yzma_mtmd_get_output_embd", int(mctx), ptr, int(embedSize))
	})
}

// MtmdBatchInit creates a batch that encodes several media chunks at once.
func MtmdBatchInit(mctx MtmdContext) (MtmdBatch, error) {
	if err := mtmdPartsReady(); err != nil {
		return 0, err
	}
	handle, err := callErr("_yzma_mtmd_batch_init", int(mctx))
	if err != nil {
		return 0, err
	}
	return MtmdBatch(handle), nil
}

// MtmdBatchFree frees a batch.
func MtmdBatchFree(batch MtmdBatch) {
	if MtmdPartsSupported() {
		callVoid("_yzma_mtmd_batch_free", int(batch))
	}
}

// MtmdBatchAddChunk adds a media chunk to a batch. The chunk must stay valid
// until the batch is encoded.
func MtmdBatchAddChunk(batch MtmdBatch, chunk MtmdInputChunk) MtmdBatchAddResult {
	rc := partsCall("_yzma_mtmd_batch_add_chunk", int(batch), int(chunk))
	if rc < 0 {
		return MtmdBatchAddError
	}
	return MtmdBatchAddResult(rc)
}

// MtmdBatchEncode runs every chunk of a batch through the projector.
func MtmdBatchEncode(batch MtmdBatch) error {
	if err := mtmdPartsReady(); err != nil {
		return err
	}
	_, err := callErr("_yzma_mtmd_batch_encode", int(batch))
	return err
}

// MtmdBatchGetOutputEmbd returns a copy of the output of one chunk after
// [MtmdBatchEncode]. The embedSize is the same as for [MtmdGetOutputEmbd].
func MtmdBatchGetOutputEmbd(batch MtmdBatch, chunk MtmdInputChunk, embedSize int32) ([]float32, error) {
	if err := mtmdPartsReady(); err != nil {
		return nil, err
	}
	return readEmbd(embedSize, func(ptr int) (int32, error) {
		return callErr("_yzma_mtmd_batch_get_output_embd", int(batch), int(chunk), ptr, int(embedSize))
	})
}

// readEmbd allocates embedSize floats, runs fill, and copies them out. The
// output of an image is large, so it does not go in the scratch memory.
func readEmbd(embedSize int32, fill func(ptr int) (int32, error)) ([]float32, error) {
	if embedSize <= 0 {
		return nil, fmt.Errorf("llamawasm: embedSize must be more than 0, got %d", embedSize)
	}
	ptr, err := malloc(int(embedSize) * 4)
	if err != nil {
		return nil, err
	}
	defer free(ptr)

	if _, err := fill(ptr); err != nil {
		return nil, err
	}
	return readFloats(ptr, int(embedSize)), nil
}

// partsCall runs a call of this file, or returns errHandle if the module does
// not have it.
func partsCall(name string, args ...any) int32 {
	if !MtmdPartsSupported() {
		return errHandle
	}
	return call(name, args...)
}

func nonNegative(n int32) int32 {
	if n < 0 {
		return 0
	}
	return n
}
