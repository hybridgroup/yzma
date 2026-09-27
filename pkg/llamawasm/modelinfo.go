//go:build js && wasm

package llamawasm

// The calls here describe a loaded model. Each one follows the same call in
// pkg/llama. A module before ABI version 7 has none of them, so each returns a
// zero value.

// ModelNEmbdInp returns the input embedding size.
func ModelNEmbdInp(model Model) int32 { return modelInt("_yzma_model_n_embd_inp", model) }

// ModelNEmbdOut returns the output embedding size.
func ModelNEmbdOut(model Model) int32 { return modelInt("_yzma_model_n_embd_out", model) }

// ModelNLayer returns the number of layers.
func ModelNLayer(model Model) int32 { return modelInt("_yzma_model_n_layer", model) }

// ModelNLayerNextN returns the number of layers that predict extra tokens.
func ModelNLayerNextN(model Model) int32 { return modelInt("_yzma_model_n_layer_nextn", model) }

// ModelNHead returns the number of attention heads.
func ModelNHead(model Model) int32 { return modelInt("_yzma_model_n_head", model) }

// ModelNHeadKV returns the number of key and value heads.
func ModelNHeadKV(model Model) int32 { return modelInt("_yzma_model_n_head_kv", model) }

// ModelNSWA returns the sliding window attention size. It is 0 when the model
// does not use sliding window attention.
func ModelNSWA(model Model) int32 { return modelInt("_yzma_model_n_swa", model) }

// ModelNClsOut returns the number of outputs of a classifier model.
func ModelNClsOut(model Model) int32 { return modelInt("_yzma_model_n_cls_out", model) }

// ModelHasEncoder reports whether the model has an encoder, which needs Encode.
func ModelHasEncoder(model Model) bool { return modelInt("_yzma_model_has_encoder", model) == 1 }

// ModelHasDecoder reports whether the model has a decoder, which needs Decode.
func ModelHasDecoder(model Model) bool { return modelInt("_yzma_model_has_decoder", model) == 1 }

// ModelIsRecurrent reports whether the model is recurrent, like Mamba and RWKV.
func ModelIsRecurrent(model Model) bool { return modelInt("_yzma_model_is_recurrent", model) == 1 }

// ModelIsHybrid reports whether the model is hybrid, like Jamba and Granite.
func ModelIsHybrid(model Model) bool { return modelInt("_yzma_model_is_hybrid", model) == 1 }

// ModelIsDiffusion reports whether the model is a diffusion model, like LLaDA.
func ModelIsDiffusion(model Model) bool { return modelInt("_yzma_model_is_diffusion", model) == 1 }

// ModelFtype returns the quantization type of the model.
func ModelFtype(model Model) Ftype { return Ftype(modelInt("_yzma_model_ftype", model)) }

// ModelRopeType returns how the model scales RoPE positions. It returns
// RopeScalingTypeUnspecified when the module has no such call.
func ModelRopeType(model Model) RopeScalingType {
	if !has("_yzma_model_rope_type") {
		return RopeScalingTypeUnspecified
	}

	// A rope type can be -1, so only the bad handle value is a failure.
	rc := call("_yzma_model_rope_type", int(model))
	if rc <= errBadHandle {
		return RopeScalingTypeUnspecified
	}
	return RopeScalingType(rc)
}

// ModelDecoderStartToken returns the decoder start token of an encoder decoder
// model. Every other model returns TokenNull.
func ModelDecoderStartToken(model Model) Token {
	if !has("_yzma_model_decoder_start_token") {
		return TokenNull
	}

	rc := call("_yzma_model_decoder_start_token", int(model))
	if rc <= errBadHandle {
		return TokenNull
	}
	return Token(rc)
}

// ModelSize returns the total size of all model tensors in bytes.
func ModelSize(model Model) uint64 { return modelUint64("_yzma_model_size", model) }

// ModelNParams returns the number of model parameters.
func ModelNParams(model Model) uint64 { return modelUint64("_yzma_model_n_params", model) }

// ModelRopeFreqScaleTrain returns the RoPE frequency scale the model was
// trained with.
func ModelRopeFreqScaleTrain(model Model) float32 {
	if !has("_yzma_model_rope_freq_scale_train") {
		return 0
	}
	return float32(callValue("_yzma_model_rope_freq_scale_train", int(model)).Float())
}

// modelInt reads an integer model property. It returns 0 when the module has
// no such call or when the call fails.
func modelInt(name string, model Model) int32 {
	if !has(name) {
		return 0
	}
	n := call(name, int(model))
	if n < 0 {
		return 0
	}
	return n
}

// modelUint64 reads a number that does not fit in an int32. The shim returns it
// as a double, which holds every integer up to 2^53.
func modelUint64(name string, model Model) uint64 {
	if !has(name) {
		return 0
	}
	v := callValue(name, int(model)).Float()
	if v < 0 {
		return 0
	}
	return uint64(v)
}

// ModelMetaCount returns the number of key value pairs in the model metadata.
func ModelMetaCount(model Model) int32 { return modelInt("_yzma_model_meta_count", model) }

// ModelMetaKeyByIndex returns the key of metadata pair i.
func ModelMetaKeyByIndex(model Model, i int32) (string, bool) {
	return metaString("_yzma_model_meta_key_by_index", 128, int(model), int(i))
}

// ModelMetaValStrByIndex returns the value of metadata pair i as a string.
func ModelMetaValStrByIndex(model Model, i int32) (string, bool) {
	return metaString("_yzma_model_meta_val_str_by_index", 32768, int(model), int(i))
}

// ModelMetaValStr returns the value of a metadata key as a string.
func ModelMetaValStr(model Model, key string) (string, bool) {
	if !has("_yzma_model_meta_val_str") {
		return "", false
	}

	keyPtr, freeKey, err := allocString(key)
	if err != nil {
		return "", false
	}
	defer freeKey()

	return metaString("_yzma_model_meta_val_str", 32768, int(model), keyPtr)
}

// ModelMetaKeyStr returns the name of a metadata key, or an empty string if
// the key is unknown.
func ModelMetaKeyStr(key ModelMetaKey) string {
	return callString("_yzma_model_meta_key_str", 64, int(key))
}

// ModelClsLabel returns the label of output i of a classifier model, or an
// empty string if the model has none.
func ModelClsLabel(model Model, i uint32) string {
	return callString("_yzma_model_cls_label", 256, int(model), int(i))
}

// metaString runs a shim meta call. These return the full value length like
// snprintf does, so a length at or above the buffer size needs a second call
// with a larger buffer.
func metaString(name string, size int, args ...any) (string, bool) {
	if !has(name) {
		return "", false
	}

	for range 2 {
		ptr, err := pieceScratch.reserve(size)
		if err != nil {
			return "", false
		}

		n := int(call(name, append(args, ptr, size)...))
		if n < 0 {
			return "", false
		}
		if n < size {
			return string(readBytes(ptr, n)), true
		}
		size = n + 1
	}
	return "", false
}
