//go:build js && wasm

package llamawasm

import (
	"syscall/js"
	"testing"
)

const fakeInfoSource = `
globalThis.__yzmaInfo = (function () {
	let rope = 2;
	let start = -1;

	return {
		module: {
			HEAPU8: new Uint8Array(64),
			_malloc: () => 0,
			_free: () => {},
			_yzma_abi_version: () => 7,
			_yzma_last_error: () => 0,
			_yzma_model_n_embd_inp: () => 768,
			_yzma_model_n_embd_out: () => 512,
			_yzma_model_n_layer: () => 24,
			_yzma_model_n_layer_nextn: () => 1,
			_yzma_model_n_head: () => 12,
			_yzma_model_n_head_kv: () => 4,
			_yzma_model_n_swa: () => 0,
			_yzma_model_n_cls_out: () => 2,
			_yzma_model_ftype: () => 15,
			_yzma_model_has_encoder: () => 0,
			_yzma_model_has_decoder: () => 1,
			_yzma_model_is_recurrent: () => 0,
			_yzma_model_is_hybrid: () => 0,
			_yzma_model_is_diffusion: () => 0,
			_yzma_model_rope_type: () => rope,
			_yzma_model_decoder_start_token: () => start,
			_yzma_model_size: () => 532517120,
			_yzma_model_n_params: () => 8000000000,
			_yzma_model_rope_freq_scale_train: () => 0.25,
			_yzma_context_n_threads: () => 16,
			_yzma_context_n_threads_batch: () => 8,
			_yzma_context_n_rs_seq: () => 1,
			_yzma_set_n_threads: () => 0,
		},
		setRope: (r) => { rope = r; },
		setStart: (s) => { start = s; },
	};
})();
`

func fakeInfo(t *testing.T) js.Value {
	t.Helper()

	js.Global().Call("eval", fakeInfoSource)
	helper := js.Global().Get("__yzmaInfo")

	previous := mod
	mod = helper.Get("module")
	t.Cleanup(func() { mod = previous })

	return helper
}

func TestModelInfo(t *testing.T) {
	fakeInfo(t)

	m := Model(1)
	for _, test := range []struct {
		name string
		got  int32
		want int32
	}{
		{"ModelNEmbdInp", ModelNEmbdInp(m), 768},
		{"ModelNEmbdOut", ModelNEmbdOut(m), 512},
		{"ModelNLayer", ModelNLayer(m), 24},
		{"ModelNLayerNextN", ModelNLayerNextN(m), 1},
		{"ModelNHead", ModelNHead(m), 12},
		{"ModelNHeadKV", ModelNHeadKV(m), 4},
		{"ModelNSWA", ModelNSWA(m), 0},
		{"ModelNClsOut", ModelNClsOut(m), 2},
	} {
		if test.got != test.want {
			t.Errorf("%s gave %d, want %d", test.name, test.got, test.want)
		}
	}

	if ModelHasEncoder(m) || !ModelHasDecoder(m) {
		t.Error("the encoder and the decoder of the model are not correct")
	}
	if ModelIsRecurrent(m) || ModelIsHybrid(m) || ModelIsDiffusion(m) {
		t.Error("the kind of the model is not correct")
	}
	if got := ModelFtype(m); got != FtypeMostlyQ4_K_M {
		t.Errorf("ModelFtype gave %v, want FtypeMostlyQ4_K_M", got)
	}
}

func TestModelSizeAndParams(t *testing.T) {
	fakeInfo(t)

	// These do not fit an int32, thus the shim gives them as a double.
	if got := ModelSize(Model(1)); got != 532517120 {
		t.Errorf("ModelSize gave %d, want 532517120", got)
	}
	if got := ModelNParams(Model(1)); got != 8000000000 {
		t.Errorf("ModelNParams gave %d, want 8000000000", got)
	}
	if got := ModelRopeFreqScaleTrain(Model(1)); got != 0.25 {
		t.Errorf("ModelRopeFreqScaleTrain gave %v, want 0.25", got)
	}
}

func TestModelRopeTypeAndStartToken(t *testing.T) {
	helper := fakeInfo(t)

	if got := ModelRopeType(Model(1)); got != RopeScalingTypeYARN {
		t.Errorf("ModelRopeType gave %v, want RopeScalingTypeYARN", got)
	}
	// A model with no decoder to start gives -1, which is a value.
	if got := ModelDecoderStartToken(Model(1)); got != TokenNull {
		t.Errorf("ModelDecoderStartToken gave %v, want TokenNull", got)
	}

	helper.Call("setRope", -1)
	helper.Call("setStart", 7)
	if got := ModelRopeType(Model(1)); got != RopeScalingTypeUnspecified {
		t.Errorf("ModelRopeType gave %v, want RopeScalingTypeUnspecified", got)
	}
	if got := ModelDecoderStartToken(Model(1)); got != Token(7) {
		t.Errorf("ModelDecoderStartToken gave %v, want 7", got)
	}
}

func TestThreadCalls(t *testing.T) {
	fakeInfo(t)

	if got := NThreads(Context(1)); got != 16 {
		t.Errorf("NThreads gave %d, want 16", got)
	}
	if got := NThreadsBatch(Context(1)); got != 8 {
		t.Errorf("NThreadsBatch gave %d, want 8", got)
	}
	if got := NRsSeq(Context(1)); got != 1 {
		t.Errorf("NRsSeq gave %d, want 1", got)
	}
	if err := SetNThreads(Context(1), 4, 4); err != nil {
		t.Errorf("SetNThreads gave %v, want no error", err)
	}
}

func TestModelInfoOldModule(t *testing.T) {
	// A module before ABI version 7 has none of the calls.
	js.Global().Call("eval", "globalThis.__yzmaOldInfo = { HEAPU8: new Uint8Array(8), _malloc: () => 0 };")

	previous := mod
	mod = js.Global().Get("__yzmaOldInfo")
	t.Cleanup(func() { mod = previous })

	if ModelNLayer(Model(1)) != 0 || ModelSize(Model(1)) != 0 || NThreads(Context(1)) != 0 {
		t.Error("a module of an earlier version must give a zero value")
	}
	if got := ModelRopeType(Model(1)); got != RopeScalingTypeUnspecified {
		t.Errorf("ModelRopeType gave %v, want RopeScalingTypeUnspecified", got)
	}
	if got := ModelDecoderStartToken(Model(1)); got != TokenNull {
		t.Errorf("ModelDecoderStartToken gave %v, want TokenNull", got)
	}
}

func TestModelMeta(t *testing.T) {
	fakeABI9(t)

	m := Model(1)
	if got := ModelMetaCount(m); got != 2 {
		t.Fatalf("ModelMetaCount gave %d, want 2", got)
	}
	if got, ok := ModelMetaKeyByIndex(m, 0); !ok || got != "general.name" {
		t.Errorf("ModelMetaKeyByIndex gave %q and %v", got, ok)
	}
	// This value is larger than the first buffer, thus it needs a second call.
	if got, ok := ModelMetaValStrByIndex(m, 1); !ok || len(got) != 40000 {
		t.Errorf("ModelMetaValStrByIndex gave %d bytes and %v, want 40000", len(got), ok)
	}
	if got, ok := ModelMetaValStr(m, "general.name"); !ok || got != "Tiny" {
		t.Errorf("ModelMetaValStr gave %q and %v", got, ok)
	}
	if _, ok := ModelMetaValStr(m, "no.such.key"); ok {
		t.Error("ModelMetaValStr found a key that is not there")
	}
	if _, ok := ModelMetaKeyByIndex(m, 5); ok {
		t.Error("ModelMetaKeyByIndex found an index that is not there")
	}
	if got := ModelMetaKeyStr(ModelMetaKeySamplingTopK); got != "general.sampling.top_k" {
		t.Errorf("ModelMetaKeyStr gave %q", got)
	}
	if got := ModelClsLabel(m, 1); got != "label1" {
		t.Errorf("ModelClsLabel gave %q", got)
	}
}

func TestModelMetaOldModule(t *testing.T) {
	fakeOld(t)

	if _, ok := ModelMetaValStr(Model(1), "general.name"); ok || ModelMetaCount(Model(1)) != 0 {
		t.Error("a module of an earlier version must give no metadata")
	}
}
