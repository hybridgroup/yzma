//go:build js && wasm

package llamawasm

import "fmt"

// ModelLoadFromFile loads a model from a file in the llama.cpp module
// filesystem.
//
// The file must be there before this call. Use WriteModelFile or
// FetchModelFile to put it there.
func ModelLoadFromFile(pathModel string, params ModelParams) (Model, error) {
	if !Loaded() {
		return 0, ErrNotLoaded
	}

	ptr, err := textScratch.reserve(len(pathModel) + 1)
	if err != nil {
		return 0, err
	}
	writeString(ptr, pathModel)

	handle, err := callErr("_yzma_model_load", ptr, int(params.NGpuLayers))
	if err != nil {
		return 0, err
	}
	return Model(handle), nil
}

// ModelFree frees a model.
func ModelFree(model Model) error {
	if !Loaded() {
		return ErrNotLoaded
	}
	callVoid("_yzma_model_free", int(model))
	return nil
}

// ModelGetVocab returns the vocabulary of a model.
func ModelGetVocab(model Model) Vocab {
	if !Loaded() {
		return 0
	}
	rc := call("_yzma_model_get_vocab", int(model))
	if rc <= errBadHandle {
		return 0
	}
	return Vocab(rc)
}

// ModelNEmbd returns the embedding size of the model.
func ModelNEmbd(model Model) int32 {
	if !Loaded() {
		return 0
	}
	return call("_yzma_model_n_embd", int(model))
}

// ModelNCtxTrain returns the context size the model was trained with.
func ModelNCtxTrain(model Model) int32 {
	if !Loaded() {
		return 0
	}
	return call("_yzma_model_n_ctx_train", int(model))
}

// ModelDesc returns a short description of the model type.
func ModelDesc(model Model) string {
	return modelString("_yzma_model_desc", model, 256)
}

// ModelChatTemplate returns the chat template stored in the model, or an empty
// string if the model has none. A name such as "systemone" returns the named
// template, which llama.cpp keeps in the tokenizer.chat_template.<name> metadata.
func ModelChatTemplate(model Model, name string) string {
	if name != "" {
		tmpl, _ := ModelMetaValStr(model, "tokenizer.chat_template."+name)
		return tmpl
	}
	return modelString("_yzma_model_chat_template", model, 8192)
}

// modelString reads a model string that the shim writes into a buffer.
func modelString(name string, model Model, size int) string {
	if !Loaded() {
		return ""
	}
	return callString(name, size, int(model))
}

// String returns the model handle as text, for debugging.
func (m Model) String() string {
	return fmt.Sprintf("model(%d)", int32(m))
}

// ChatApplyTemplate formats one message with the model chat template.
//
// One message is enough for a prompt with a question about an image. For a
// multi turn chat or tool calling, render the template from
// [ModelChatTemplate] with the yzma template package.
//
// addAssistant adds the start of the assistant turn. This makes the model
// answer instead of continuing the message.
func ChatApplyTemplate(model Model, role, content string, addAssistant bool) (string, error) {
	if !Loaded() {
		return "", ErrNotLoaded
	}
	if !has("_yzma_chat_apply_template") {
		return "", ErrNoMultimodal
	}

	roleLen := len(role) + 1
	rolePtr, err := textScratch.reserve(roleLen + len(content) + 1)
	if err != nil {
		return "", err
	}
	writeString(rolePtr, role)

	contentPtr := rolePtr + roleLen
	writeString(contentPtr, content)

	size := len(content) + 512
	for {
		outPtr, err := pieceScratch.reserve(size)
		if err != nil {
			return "", err
		}

		n := call("_yzma_chat_apply_template", int(model), rolePtr, contentPtr,
			boolToInt(addAssistant), outPtr, size)
		switch {
		case n == errTooSmall && size < 1<<20:
			size *= 4
			continue
		case n < 0:
			return "", shimError("_yzma_chat_apply_template", n)
		default:
			return string(readBytes(outPtr, int(n))), nil
		}
	}
}
