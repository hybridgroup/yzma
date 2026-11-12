//go:build js && wasm

package llamawasm

import (
	"errors"
	"syscall/js"
	"testing"
)

// The fake ABI version 13 module holds one audio bitmap, one text chunk, and
// one image chunk, and records the arguments it gets.
const fakeABI13Source = `
globalThis.__yzmaABI13 = (function () {
	const heap = new Uint8Array(1 << 16);
	const flagPtr = 16;
	let next = 64;
	let last = {};

	function i32(p) { return new Int32Array(heap.buffer, p, 1)[0]; }
	function f32s(p, n) { return Array.from(new Float32Array(heap.buffer, p, n)); }
	function str(p) { let s = ""; while (heap[p]) { s += String.fromCharCode(heap[p++]); } return s; }
	function put(p, s) { for (let i = 0; i < s.length; i++) heap[p + i] = s.charCodeAt(i); heap[p + s.length] = 0; return s.length; }

	return {
		module: {
			HEAPU8: heap,
			_malloc: (n) => { const p = next; next += n + (8 - (n % 8)); return p; },
			_free: () => {},
			_yzma_abi_version: () => 13,
			_yzma_last_error: () => 0,
			_yzma_set_abort_callback: (ctx, on) => { last.abort = [ctx, on]; return 0; },
			_yzma_abort_flag: () => flagPtr,
			_yzma_abort_set: (v) => { new Int32Array(heap.buffer, flagPtr, 1)[0] = v ? 1 : 0; },
			_yzma_abort_get: () => i32(flagPtr),
			_yzma_mtmd_get_audio_sample_rate: () => 16000,
			_yzma_mtmd_bitmap_init_from_audio: (n, p) => { last.audio = f32s(p, n); return 5; },
			_yzma_mtmd_bitmap_get_nx: () => 3,
			_yzma_mtmd_bitmap_get_ny: () => 1,
			_yzma_mtmd_bitmap_get_n_bytes: () => 12,
			_yzma_mtmd_bitmap_get_data: (b, p, cap) => { for (let i = 0; i < 12; i++) heap[p + i] = i; return 12; },
			_yzma_mtmd_bitmap_is_audio: () => 1,
			_yzma_mtmd_bitmap_set_id: (b, p) => { last.id = str(p); return 0; },
			_yzma_mtmd_bitmap_get_id: (b, p, cap) => put(p, "audio-1"),
			_yzma_mtmd_input_chunks_get: (list, idx) => idx < 2 ? 10 + idx : -1,
			_yzma_mtmd_input_chunk_get_type: (c) => c === 10 ? 0 : 1,
			_yzma_mtmd_input_chunk_get_n_tokens: (c) => c === 10 ? 3 : 64,
			_yzma_mtmd_input_chunk_get_tokens_text: (c, p, cap) => {
				if (c !== 10) return 0;
				new Int32Array(heap.buffer, p, 3).set([7, 8, 9]);
				return 3;
			},
			_yzma_mtmd_input_chunk_get_tokens_image: (c) => c === 11 ? 20 : 0,
			_yzma_mtmd_image_tokens_get_decoder_pos: (it, pos0, i, p) => {
				new Int32Array(heap.buffer, p, 4).set([pos0, i % 8, Math.floor(i / 8), 0]);
				return 4;
			},
			_yzma_mtmd_get_output_embd: (m, p, n) => {
				const out = new Float32Array(heap.buffer, p, n);
				for (let i = 0; i < n; i++) out[i] = i / 2;
				return n;
			},
			_yzma_mtmd_batch_add_chunk: () => 2,
		},
		last: () => JSON.stringify(last),
	};
})();
`

func fakeABI13(t *testing.T) js.Value {
	t.Helper()

	js.Global().Call("eval", fakeABI13Source)
	helper := js.Global().Get("__yzmaABI13")

	previous := mod
	mod = helper.Get("module")
	t.Cleanup(func() {
		mod = previous
		for _, s := range []*scratch{&tokenScratch, &outScratch, &pieceScratch} {
			s.ptr, s.size = 0, 0
		}
	})

	return helper
}

func lastOf(helper js.Value, key string) js.Value {
	return js.Global().Get("JSON").Call("parse", helper.Call("last")).Get(key)
}

func TestAbortFlag(t *testing.T) {
	helper := fakeABI13(t)

	if err := SetAbortCallback(Context(3), true); err != nil {
		t.Fatalf("SetAbortCallback gave %v", err)
	}
	got := lastOf(helper, "abort")
	if got.Index(0).Int() != 3 || got.Index(1).Int() != 1 {
		t.Fatalf("shim got %v", got)
	}

	if Aborted() {
		t.Fatal("flag is set before Abort")
	}
	Abort()
	if !Aborted() {
		t.Fatal("flag is not set after Abort")
	}
	AbortReset()
	if Aborted() {
		t.Fatal("flag is set after AbortReset")
	}

	if v := AbortFlagView(); !v.IsUndefined() {
		t.Fatal("AbortFlagView must be undefined without threads")
	}
}

func TestMtmdAudioBitmap(t *testing.T) {
	helper := fakeABI13(t)

	if got := MtmdGetAudioSampleRate(1); got != 16000 {
		t.Fatalf("sample rate %d", got)
	}

	b, err := MtmdBitmapInitFromAudio([]float32{0.5, -0.25, 1})
	if err != nil || b != 5 {
		t.Fatalf("MtmdBitmapInitFromAudio gave %d, %v", b, err)
	}
	samples := lastOf(helper, "audio")
	if samples.Length() != 3 || samples.Index(1).Float() != -0.25 {
		t.Fatalf("shim got %v", samples)
	}

	if !MtmdBitmapIsAudio(b) || MtmdBitmapGetNx(b) != 3 || MtmdBitmapGetNy(b) != 1 {
		t.Fatal("bitmap getters are wrong")
	}
	data := MtmdBitmapGetData(b)
	if len(data) != 12 || data[11] != 11 {
		t.Fatalf("MtmdBitmapGetData gave %v", data)
	}

	if err := MtmdBitmapSetId(b, "clip"); err != nil {
		t.Fatal(err)
	}
	if got := lastOf(helper, "id").String(); got != "clip" {
		t.Fatalf("shim got id %q", got)
	}
	if got := MtmdBitmapGetId(b); got != "audio-1" {
		t.Fatalf("MtmdBitmapGetId gave %q", got)
	}

	if _, err := MtmdBitmapInitFromAudio(nil); err == nil {
		t.Fatal("no samples must give an error")
	}
}

func TestMtmdChunkParts(t *testing.T) {
	fakeABI13(t)

	text := MtmdInputChunksGet(1, 0)
	image := MtmdInputChunksGet(1, 1)
	if text != 10 || image != 11 || MtmdInputChunksGet(1, 2) != 0 {
		t.Fatalf("MtmdInputChunksGet gave %d, %d", text, image)
	}

	if MtmdInputChunkGetType(text) != MtmdInputChunkTypeText || MtmdInputChunkGetType(image) != MtmdInputChunkTypeImage {
		t.Fatal("chunk types are wrong")
	}

	tokens := MtmdInputChunkGetTokensText(text)
	if len(tokens) != 3 || tokens[0] != 7 || tokens[2] != 9 {
		t.Fatalf("MtmdInputChunkGetTokensText gave %v", tokens)
	}
	if MtmdInputChunkGetTokensText(image) != nil {
		t.Fatal("an image chunk must give no text tokens")
	}

	it := MtmdInputChunkGetTokensImage(image)
	if it != 20 || MtmdInputChunkGetTokensImage(text) != 0 {
		t.Fatalf("MtmdInputChunkGetTokensImage gave %d", it)
	}

	pos := MtmdImageTokensGetDecoderPos(it, 5, 10)
	if pos != (MtmdDecoderPos{T: 5, X: 2, Y: 1}) {
		t.Fatalf("MtmdImageTokensGetDecoderPos gave %+v", pos)
	}

	embd, err := MtmdGetOutputEmbd(1, 4)
	if err != nil || len(embd) != 4 || embd[3] != 1.5 {
		t.Fatalf("MtmdGetOutputEmbd gave %v, %v", embd, err)
	}

	if got := MtmdBatchAddChunk(1, image); got != MtmdBatchAddTooLarge {
		t.Fatalf("MtmdBatchAddChunk gave %d", got)
	}
}

func TestABI13CallsOnOlderModule(t *testing.T) {
	fakeABI11(t)

	if err := SetAbortCallback(Context(1), true); !errors.Is(err, ErrNoAbort) {
		t.Fatalf("SetAbortCallback gave %v", err)
	}
	if _, err := MtmdBitmapInitFromAudio([]float32{1}); !errors.Is(err, ErrNoMtmdParts) {
		t.Fatalf("MtmdBitmapInitFromAudio gave %v", err)
	}
	if MtmdInputChunksGet(1, 0) != 0 || MtmdGetAudioSampleRate(1) != -1 || Aborted() {
		t.Fatal("an older module must give zero values")
	}
}
