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
