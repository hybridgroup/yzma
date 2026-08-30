//go:build js && wasm

package llamawasm

import "fmt"

// ModelLoadFromFile loads a model from a file in the filesystem of the
// llama.cpp module.
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

// ModelGetVocab gives the vocabulary of a model.
func ModelGetVocab(model Model) Vocab {
	if !Loaded() {
		return 0
	}
	return Vocab(call("_yzma_model_get_vocab", int(model)))
}

// ModelNEmbd gives the number of values in an embedding of the model.
func ModelNEmbd(model Model) int32 {
	if !Loaded() {
		return 0
	}
	return call("_yzma_model_n_embd", int(model))
}

// ModelNCtxTrain gives the size of the context that the model was trained
// with.
func ModelNCtxTrain(model Model) int32 {
	if !Loaded() {
		return 0
	}
	return call("_yzma_model_n_ctx_train", int(model))
}

// ModelDesc gives a short description of the type of the model.
func ModelDesc(model Model) string {
	return modelString("_yzma_model_desc", model, 256)
}

// ModelChatTemplate gives the chat template that is in the model, or an empty
// string if the model has none.
//
// The name argument is not used. It is here to keep the same shape as
// llama.ModelChatTemplate.
func ModelChatTemplate(model Model, name string) string {
	return modelString("_yzma_model_chat_template", model, 8192)
}

// modelString reads a string that the shim writes into a buffer, and takes a
// larger buffer if the first one is too small.
func modelString(name string, model Model, size int) string {
	if !Loaded() {
		return ""
	}

	for {
		ptr, err := pieceScratch.reserve(size)
		if err != nil {
			return ""
		}

		n := call(name, int(model), ptr, size)
		switch {
		case n < 0 && n == errTooSmall && size < 1<<20:
			size *= 4
			continue
		case n <= 0:
			return ""
		default:
			return string(readBytes(ptr, int(n)))
		}
	}
}

// String gives the handle of the model as text, which helps while debugging.
func (m Model) String() string {
	return fmt.Sprintf("model(%d)", int32(m))
}

// ChatApplyTemplate puts one message into the chat format of the model.
//
// One message is enough for a prompt with a question about an image. A chat with
// turns needs the template of the model, which [ModelChatTemplate] gives.
//
// addAssistant adds the opening of the turn of the assistant, which is what
// makes the model answer instead of carrying on the message.
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
