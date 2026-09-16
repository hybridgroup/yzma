//go:build js && wasm

package llamawasm

// The calls here describe a loaded model. Each one follows the same call in
// pkg/llama. A module before ABI version 7 has none of them, thus each gives a
// zero value.

// ModelNEmbdInp gives the number of values in an embedding that goes in.
func ModelNEmbdInp(model Model) int32 { return modelInt("_yzma_model_n_embd_inp", model) }

// ModelNEmbdOut gives the number of values in an embedding that comes out.
func ModelNEmbdOut(model Model) int32 { return modelInt("_yzma_model_n_embd_out", model) }

// ModelNLayer gives the number of layers.
func ModelNLayer(model Model) int32 { return modelInt("_yzma_model_n_layer", model) }

// ModelNLayerNextN gives the number of layers that predict more than one token.
func ModelNLayerNextN(model Model) int32 { return modelInt("_yzma_model_n_layer_nextn", model) }

// ModelNHead gives the number of attention heads.
func ModelNHead(model Model) int32 { return modelInt("_yzma_model_n_head", model) }

// ModelNHeadKV gives the number of heads of the keys and the values.
func ModelNHeadKV(model Model) int32 { return modelInt("_yzma_model_n_head_kv", model) }

// ModelNSWA gives the size of the window of the attention that slides. It is 0
// when the attention of the model does not slide.
func ModelNSWA(model Model) int32 { return modelInt("_yzma_model_n_swa", model) }

// ModelNClsOut gives the number of outputs of a classifier model.
func ModelNClsOut(model Model) int32 { return modelInt("_yzma_model_n_cls_out", model) }

// ModelHasEncoder tells if the model has an encoder, which needs Encode.
func ModelHasEncoder(model Model) bool { return modelInt("_yzma_model_has_encoder", model) == 1 }

// ModelHasDecoder tells if the model has a decoder, which needs Decode.
func ModelHasDecoder(model Model) bool { return modelInt("_yzma_model_has_decoder", model) == 1 }

// ModelIsRecurrent tells if the model is recurrent, as Mamba and RWKV are.
func ModelIsRecurrent(model Model) bool { return modelInt("_yzma_model_is_recurrent", model) == 1 }

// ModelIsHybrid tells if the model is hybrid, as Jamba and Granite are.
func ModelIsHybrid(model Model) bool { return modelInt("_yzma_model_is_hybrid", model) == 1 }

// ModelIsDiffusion tells if the model is a diffusion model, as LLaDA is.
func ModelIsDiffusion(model Model) bool { return modelInt("_yzma_model_is_diffusion", model) == 1 }

// ModelFtype gives the kind of the quantization of the model.
func ModelFtype(model Model) Ftype { return Ftype(modelInt("_yzma_model_ftype", model)) }

// ModelRopeType gives how the model scales the positions of RoPE. It gives
// RopeScalingTypeUnspecified when the module has no such call.
func ModelRopeType(model Model) RopeScalingType {
	if !has("_yzma_model_rope_type") {
		return RopeScalingTypeUnspecified
	}

	// A rope type can be -1, thus only the value of a bad handle is a failure.
	rc := call("_yzma_model_rope_type", int(model))
	if rc <= errBadHandle {
		return RopeScalingTypeUnspecified
	}
	return RopeScalingType(rc)
}

// ModelDecoderStartToken gives the token that starts the decoder of a model
// that has an encoder and a decoder. Every other model gives TokenNull.
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

// ModelSize gives the size of every tensor of the model in bytes.
func ModelSize(model Model) uint64 { return modelUint64("_yzma_model_size", model) }

// ModelNParams gives the number of parameters of the model.
func ModelNParams(model Model) uint64 { return modelUint64("_yzma_model_n_params", model) }

// ModelRopeFreqScaleTrain gives the scale of the frequency of RoPE that the
// model was trained with.
func ModelRopeFreqScaleTrain(model Model) float32 {
	if !has("_yzma_model_rope_freq_scale_train") {
		return 0
	}
	return float32(callValue("_yzma_model_rope_freq_scale_train", int(model)).Float())
}

// modelInt reads a whole number of a model. It gives 0 when the module has no
// such call and when the call fails.
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

// modelUint64 reads a number that does not fit an int32. The shim gives it as
// a double, which holds every whole number to 2^53.
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
