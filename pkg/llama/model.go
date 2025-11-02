package llama

import (
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/utils"
	"github.com/jupiterrider/ffi"
)

var (
	FFITypeModelParams = ffi.NewType(&ffi.TypePointer, &ffi.TypePointer, &ffi.TypeSint32,
		&ffi.TypeSint32, &ffi.TypeSint32,
		&ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer,
		&ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8, &ffi.TypeUint8)
)

var (
	// LLAMA_API struct llama_model_params          llama_model_default_params(void);
	modelDefaultParamsFunc ffi.Fun

	// LLAMA_API struct llama_model * llama_model_load_from_file(
	//                          const char * path_model,
	//           				struct llama_model_params   params);
	modelLoadFromFileFunc ffi.Fun

	// LLAMA_API struct llama_model_params          llama_model_default_params(void);
	modelFreeFunc ffi.Fun

	// LLAMA_API struct llama_context * llama_init_from_model(
	//                  struct llama_model * model,
	//         			struct llama_context_params   params);
	initFromModelFunc ffi.Fun

	// LLAMA_API const char * llama_model_chat_template(const struct llama_model * model, const char * name);
	modelChatTemplateFunc ffi.Fun

	// LLAMA_API bool llama_model_has_encoder(const struct llama_model * model);
	modelHasEncoderFunc ffi.Fun

	// LLAMA_API bool llama_model_has_decoder(const struct llama_model * model);
	modelHasDecoderFunc ffi.Fun

	// LLAMA_API llama_token llama_model_decoder_start_token(const struct llama_model * model);
	modelDecoderStartTokenFunc ffi.Fun

	// LLAMA_API int32_t llama_model_n_ctx_train(const struct llama_model * model);
	modelNCtxTrainFunc ffi.Fun

	// LLAMA_API int32_t llama_model_n_embd     (const struct llama_model * model);
	modelNEmbdFunc ffi.Fun

	// LLAMA_API int32_t llama_model_n_layer    (const struct llama_model * model);
	modelNLayerFunc ffi.Fun

	// LLAMA_API int32_t llama_model_n_head     (const struct llama_model * model);
	modelNHeadFunc ffi.Fun

	// LLAMA_API int32_t llama_model_n_head_kv  (const struct llama_model * model);
	modelNHeadKVFunc ffi.Fun

	// LLAMA_API int32_t llama_model_n_swa      (const struct llama_model * model);
	modelNSWAFunc ffi.Fun

	// LLAMA_API uint32_t llama_model_n_cls_out(const struct llama_model * model);
	modelNClsOutFunc ffi.Fun

	// LLAMA_API const char * llama_model_cls_label(const struct llama_model * model, uint32_t i);
	modelClsLabelFunc ffi.Fun

	// LLAMA_API int32_t llama_model_desc(const struct llama_model * model, char * buf, size_t buf_size);
	modelDescFunc ffi.Fun

	// LLAMA_API uint64_t llama_model_size(const struct llama_model * model);
	modelSizeFunc ffi.Fun

	// LLAMA_API bool llama_model_is_recurrent(const struct llama_model * model);
	modelIsRecurrentFunc ffi.Fun

	// LLAMA_API bool llama_model_is_hybrid(const struct llama_model * model);
	modelIsHybridFunc ffi.Fun

	// LLAMA_API bool llama_model_is_diffusion(const struct llama_model * model);
	modelIsDiffusionFunc ffi.Fun

	// LLAMA_API float llama_model_rope_freq_scale_train(const struct llama_model * model);
	modelRopeFreqScaleTrainFunc ffi.Fun

	// LLAMA_API enum llama_rope_type llama_model_rope_type(const struct llama_model * model);
	modelRopeTypeFunc ffi.Fun

	// Get metadata value as a string by key name
	// LLAMA_API int32_t llama_model_meta_val_str(const struct llama_model * model, const char * key, char * buf, size_t buf_size);
	modelMetaValStrFunc ffi.Fun

	// Get the number of metadata key/value pairs
	// LLAMA_API int32_t llama_model_meta_count(const struct llama_model * model);
	modelMetaCountFunc ffi.Fun

	// Get metadata key name by index
	// LLAMA_API int32_t llama_model_meta_key_by_index(const struct llama_model * model, int32_t i, char * buf, size_t buf_size);
	modelMetaKeyByIndexFunc ffi.Fun

	// Get metadata value as a string by index
	// LLAMA_API int32_t llama_model_meta_val_str_by_index(const struct llama_model * model, int32_t i, char * buf, size_t buf_size);
	modelMetaValStrByIndexFunc ffi.Fun
)

func loadModelFuncs(lib ffi.Lib) error {
	var err error

	if modelDefaultParamsFunc, err = lib.Prep("llama_model_default_params", &FFITypeModelParams); err != nil {
		return loadError("llama_model_default_params", err)
	}

	if modelLoadFromFileFunc, err = lib.Prep("llama_model_load_from_file", &ffi.TypePointer, &ffi.TypePointer, &FFITypeModelParams); err != nil {
		return loadError("llama_model_load_from_file", err)
	}

	if modelFreeFunc, err = lib.Prep("llama_model_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("llama_model_free", err)
	}

	if initFromModelFunc, err = lib.Prep("llama_init_from_model", &ffi.TypePointer, &ffi.TypePointer, &FFITypeContextParams); err != nil {
		return loadError("llama_init_from_model", err)
	}

	if modelChatTemplateFunc, err = lib.Prep("llama_model_chat_template", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("llama_model_chat_template", err)
	}

	if modelHasEncoderFunc, err = lib.Prep("llama_model_has_encoder", &ffi.TypeUint8, &ffi.TypePointer); err != nil {
		return loadError("llama_model_has_encoder", err)
	}

	if modelHasDecoderFunc, err = lib.Prep("llama_model_has_decoder", &ffi.TypeUint8, &ffi.TypePointer); err != nil {
		return loadError("llama_model_has_decoder", err)
	}

	if modelDecoderStartTokenFunc, err = lib.Prep("llama_model_decoder_start_token", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("llama_model_decoder_start_token", err)
	}

	if modelNCtxTrainFunc, err = lib.Prep("llama_model_n_ctx_train", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("llama_model_n_ctx_train", err)
	}

	if modelNEmbdFunc, err = lib.Prep("llama_model_n_embd", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("llama_model_n_embd", err)
	}

	if modelNLayerFunc, err = lib.Prep("llama_model_n_layer", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("llama_model_n_embd", err)
	}

	if modelNHeadFunc, err = lib.Prep("llama_model_n_head", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("llama_model_n_head", err)
	}

	if modelNHeadKVFunc, err = lib.Prep("llama_model_n_head_kv", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("llama_model_n_head_kv", err)
	}

	if modelNSWAFunc, err = lib.Prep("llama_model_n_swa", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("llama_model_n_swa", err)
	}

	if modelNClsOutFunc, err = lib.Prep("llama_model_n_cls_out", &ffi.TypeUint32, &ffi.TypePointer); err != nil {
		return loadError("llama_model_n_cls_out", err)
	}

	if modelClsLabelFunc, err = lib.Prep("llama_model_cls_label", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32); err != nil {
		return loadError("llama_model_cls_label", err)
	}

	if modelDescFunc, err = lib.Prep("llama_model_desc", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32); err != nil {
		return loadError("llama_model_desc", err)
	}

	if modelSizeFunc, err = lib.Prep("llama_model_size", &ffi.TypeUint64, &ffi.TypePointer); err != nil {
		return loadError("llama_model_size", err)
	}

	if modelIsRecurrentFunc, err = lib.Prep("llama_model_is_recurrent", &ffi.TypeUint8, &ffi.TypePointer); err != nil {
		return loadError("llama_model_is_recurrent", err)
	}

	if modelIsHybridFunc, err = lib.Prep("llama_model_is_hybrid", &ffi.TypeUint8, &ffi.TypePointer); err != nil {
		return loadError("llama_model_is_hybrid", err)
	}

	if modelIsDiffusionFunc, err = lib.Prep("llama_model_is_diffusion", &ffi.TypeUint8, &ffi.TypePointer); err != nil {
		return loadError("llama_model_is_diffusion", err)
	}

	if modelRopeFreqScaleTrainFunc, err = lib.Prep("llama_model_rope_freq_scale_train", &ffi.TypeFloat, &ffi.TypePointer); err != nil {
		return loadError("llama_model_rope_freq_scale_train", err)
	}

	if modelRopeTypeFunc, err = lib.Prep("llama_model_rope_type", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("llama_model_rope_type", err)
	}

	if modelMetaValStrFunc, err = lib.Prep("llama_model_meta_val_str", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32); err != nil {
		return loadError("llama_model_meta_val_str", err)
	}

	if modelMetaCountFunc, err = lib.Prep("llama_model_meta_count", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("llama_model_meta_count", err)
	}

	if modelMetaKeyByIndexFunc, err = lib.Prep("llama_model_meta_key_by_index", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeUint32); err != nil {
		return loadError("llama_model_meta_key_by_index", err)
	}

	if modelMetaValStrByIndexFunc, err = lib.Prep("llama_model_meta_val_str_by_index", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeUint32); err != nil {
		return loadError("llama_model_meta_val_str_by_index", err)
	}

	return nil
}

// ModelDefaultParams returns default parameters for loading a Model.
func ModelDefaultParams() ModelParams {
	var p ModelParams
	modelDefaultParamsFunc.Call(unsafe.Pointer(&p))
	return p
}

// ModelLoadFromFile loads a Model from a GGUF file.
func ModelLoadFromFile(pathModel string, params ModelParams) Model {
	var model Model
	file := &[]byte(pathModel + "\x00")[0]
	modelLoadFromFileFunc.Call(unsafe.Pointer(&model), unsafe.Pointer(&file), unsafe.Pointer(&params))
	return model
}

// ModelFree frees a previously opened model.
func ModelFree(model Model) {
	modelFreeFunc.Call(nil, unsafe.Pointer(&model))
}

// InitFromModel initializes a previously loaded Model, and then returns a new Context.
func InitFromModel(model Model, params ContextParams) Context {
	var ctx Context
	initFromModelFunc.Call(unsafe.Pointer(&ctx), unsafe.Pointer(&model), unsafe.Pointer(&params))

	return ctx
}

// ModelChatTemplate returns a named chat template for the Model.
func ModelChatTemplate(model Model, name string) string {
	var template *byte
	var n *byte
	if len(name) > 0 {
		n = &[]byte(name + "\x00")[0]
	}
	modelChatTemplateFunc.Call(unsafe.Pointer(&template), unsafe.Pointer(&model), unsafe.Pointer(&n))

	return utils.BytePtrToString(template)
}

// ModelHasEncoder returns if the Model has an encoder.
func ModelHasEncoder(model Model) bool {
	var result ffi.Arg
	modelHasEncoderFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&model))

	return result.Bool()
}

// ModelHasDecoder returns if the Model has an decoder.
func ModelHasDecoder(model Model) bool {
	var result ffi.Arg
	modelHasDecoderFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&model))

	return result.Bool()
}

// ModelDecoderStartToken returns the start Token for the Model's decoder.
func ModelDecoderStartToken(model Model) Token {
	var result ffi.Arg
	modelDecoderStartTokenFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&model))

	return Token(result)
}

func ModelNCtxTrain(model Model) int32 {
	var result ffi.Arg
	modelNCtxTrainFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&model))

	return int32(result)
}

func ModelNEmbd(model Model) int32 {
	var result ffi.Arg
	modelNEmbdFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&model))

	return int32(result)
}

func ModelNLayer(model Model) int32 {
	var result ffi.Arg
	modelNLayerFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&model))

	return int32(result)
}

func ModelNHead(model Model) int32 {
	var result ffi.Arg
	modelNHeadFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&model))

	return int32(result)
}

func ModelNHeadKV(model Model) int32 {
	var result ffi.Arg
	modelNHeadKVFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&model))

	return int32(result)
}

func ModelNSWA(model Model) int32 {
	var result ffi.Arg
	modelNSWAFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&model))

	return int32(result)
}

// ModelNClsOut returns the number of classifier outputs (only valid for classifier models).
func ModelNClsOut(model Model) uint32 {
	var nClsOut ffi.Arg
	modelNClsOutFunc.Call(unsafe.Pointer(&nClsOut), unsafe.Pointer(&model))
	return uint32(nClsOut)
}

// ModelClsLabel returns the label of a classifier output by index.
func ModelClsLabel(model Model, index uint32) string {
	var labelPtr *byte
	modelClsLabelFunc.Call(unsafe.Pointer(&labelPtr), unsafe.Pointer(&model), unsafe.Pointer(&index))

	if labelPtr == nil {
		return ""
	}

	return utils.BytePtrToString(labelPtr)
}

// ModelDesc retrieves a string describing the model type.
func ModelDesc(model Model) string {
	buf := make([]byte, 128)
	b := unsafe.SliceData(buf)
	bLen := int32(len(buf))

	var result ffi.Arg
	modelDescFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&model), unsafe.Pointer(&b), &bLen)

	if int32(result) < 0 {
		return ""
	}

	return string(buf[:int32(result)])
}

// ModelSize returns the total size of all tensors in the model in bytes.
func ModelSize(model Model) uint64 {
	var size ffi.Arg
	modelSizeFunc.Call(unsafe.Pointer(&size), unsafe.Pointer(&model))
	return uint64(size)
}

// ModelIsRecurrent returns true if the model is recurrent.
func ModelIsRecurrent(model Model) bool {
	var result ffi.Arg
	modelIsRecurrentFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&model))
	return result.Bool()
}

// ModelIsHybrid returns true if the model is hybrid.
func ModelIsHybrid(model Model) bool {
	var result ffi.Arg
	modelIsHybridFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&model))
	return result.Bool()
}

// ModelIsDiffusion returns true if the model is diffusion-based.
func ModelIsDiffusion(model Model) bool {
	var result ffi.Arg
	modelIsDiffusionFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&model))
	return result.Bool()
}

// ModelRopeFreqScaleTrain retrieves the model's RoPE frequency scaling factor.
func ModelRopeFreqScaleTrain(model Model) float32 {
	var freqScale ffi.Arg
	modelRopeFreqScaleTrainFunc.Call(unsafe.Pointer(&freqScale), unsafe.Pointer(&model))
	return float32(freqScale)
}

// ModelRopeType retrieves the RoPE type of the model.
func ModelRopeType(model Model) RopeScalingType {
	var ropeType ffi.Arg
	modelRopeTypeFunc.Call(unsafe.Pointer(&ropeType), unsafe.Pointer(&model))
	return RopeScalingType(int32(ropeType))
}

// Warmup is to warm-up a model.
func Warmup(lctx Context, model Model) {
	vocab := ModelGetVocab(model)

	SetWarmup(lctx, true)

	tokens := make([]Token, 0)
	bos := VocabBOS(vocab)
	eos := VocabEOS(vocab)

	if bos != TokenNull {
		tokens = append(tokens, bos)
	}
	if eos != TokenNull {
		tokens = append(tokens, eos)
	}
	if len(tokens) == 0 {
		tokens = append(tokens, 0)
	}

	if ModelHasEncoder(model) {
		batch := BatchGetOne(tokens)
		Encode(lctx, batch)

		start := ModelDecoderStartToken(model)
		if start == TokenNull {
			start = bos
		}
		tokens = append([]Token{}, start)
	}

	if ModelHasDecoder(model) {
		batch := BatchGetOne(tokens)
		Decode(lctx, batch)
	}

	mem := GetMemory(lctx)
	MemoryClear(mem, true)

	Synchronize(lctx)

	// llama_perf_context_reset(lctx);
	SetWarmup(lctx, false)
}

// ModelMetaValStr gets metadata value as a string by key name.
// Returns the string and true on success, or "" and false on failure.
func ModelMetaValStr(model Model, key string) (string, bool) {
	buf := make([]byte, 128)
	b := unsafe.SliceData(buf)
	bLen := int32(len(buf))
	keyPtr, _ := utils.BytePtrFromString(key)
	var result ffi.Arg
	modelMetaValStrFunc.Call(
		unsafe.Pointer(&result),
		unsafe.Pointer(&model),
		unsafe.Pointer(&keyPtr),
		unsafe.Pointer(&b),
		&bLen,
	)
	if int32(result) < 0 {
		return "", false
	}
	return string(buf[:int32(result)]), true
}

// ModelMetaCount gets the number of metadata key/value pairs.
func ModelMetaCount(model Model) int32 {
	var result ffi.Arg
	modelMetaCountFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&model))
	return int32(result)
}

// ModelMetaKeyByIndex gets metadata key name by index.
// Returns the string and true on success, or "" and false on failure.
func ModelMetaKeyByIndex(model Model, i int32) (string, bool) {
	buf := make([]byte, 128)
	b := unsafe.SliceData(buf)
	bLen := int32(len(buf))

	var result ffi.Arg
	modelMetaKeyByIndexFunc.Call(
		unsafe.Pointer(&result),
		unsafe.Pointer(&model),
		&i,
		unsafe.Pointer(&b),
		&bLen)
	if int32(result) < 0 {
		return "", false
	}
	return string(buf[:int32(result)]), true
}

// ModelMetaValStrByIndex gets metadata value as a string by index.
// Returns the string and true on success, or "" and false on failure.
func ModelMetaValStrByIndex(model Model, i int32) (string, bool) {
	buf := make([]byte, 128)
	b := unsafe.SliceData(buf)
	bLen := int32(len(buf))

	var result ffi.Arg
	modelMetaValStrByIndexFunc.Call(
		unsafe.Pointer(&result),
		unsafe.Pointer(&model),
		&i,
		unsafe.Pointer(&b),
		&bLen)
	if int32(result) < 0 {
		return "", false
	}
	return string(buf[:int32(result)]), true
}
