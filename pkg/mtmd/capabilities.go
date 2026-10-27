package mtmd

import (
	"errors"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/loader"
	"github.com/jupiterrider/ffi"
)

// ErrNoCapabilities is returned when libmtmd does not have mtmd_helper_get_capabilities.
var ErrNoCapabilities = errors.New("mtmd_helper_get_capabilities not found in library")

// HelperCap mirrors the C struct mtmd_helper_cap.
// A capability only means the native library supports it, not that the calling application does.
type HelperCap struct {
	// Chat is false for TTS only models with no chat template.
	Chat bool

	Vision bool
	Audio  bool
	// Video also depends on the MTMD_VIDEO build option.
	Video bool

	// AudioStream is realtime transcription.
	AudioStream bool

	TextGen  bool
	AudioGen bool
	// AudioGenSpkRef means audio generation accepts a speaker reference audio.
	AudioGenSpkRef bool
}

var (
	// ffiTypeHelperCap mirrors struct mtmd_helper_cap
	ffiTypeHelperCap = ffi.NewType(&ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8,
		&ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8)

	// MTMD_API struct mtmd_helper_cap mtmd_helper_get_capabilities(const struct llama_model * model, const struct mtmd_context * mctx);
	helperGetCapabilitiesFunc ffi.Fun

	hasHelperGetCapabilities bool
)

func loadCapabilityFuncs(lib loader.Lib) error {
	var err error

	// Older libraries do not have this symbol, so it is optional.
	helperGetCapabilitiesFunc, err = lib.Prep("mtmd_helper_get_capabilities", &ffiTypeHelperCap, &ffi.TypePointer, &ffi.TypePointer)
	hasHelperGetCapabilities = err == nil

	return nil
}

// HelperGetCapabilities returns the capabilities of the model and its projector.
// ctx can be 0 for a text only model.
func HelperGetCapabilities(model llama.Model, ctx Context) (HelperCap, error) {
	var result HelperCap
	if !hasHelperGetCapabilities {
		return result, ErrNoCapabilities
	}
	if model == 0 {
		return result, nil
	}
	helperGetCapabilitiesFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&model), unsafe.Pointer(&ctx))

	return result, nil
}
