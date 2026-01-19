package llama

import (
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/utils"
	"github.com/jupiterrider/ffi"
)

var (
	// LLAMA_API int32_t llama_chat_apply_template(
	//                         const char * tmpl,
	//    const struct llama_chat_message * chat,
	//                             size_t   n_msg,
	//                               bool   add_ass,
	//                               char * buf,
	//                            int32_t   length);
	chatApplyTemplateFunc ffi.Fun

	// LLAMA_API int32_t llama_chat_builtin_templates(const char ** output, size_t len);
	chatBuiltinTemplatesFunc ffi.Fun
)

func loadChatFuncs(lib ffi.Lib) error {
	var err error
	if chatApplyTemplateFunc, err = lib.Prep("llama_chat_apply_template", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32,
		&ffi.TypeUint8, &ffi.TypePointer, &ffi.TypeSint32); err != nil {

		return loadError("llama_chat_apply_template", err)
	}

	if chatBuiltinTemplatesFunc, err = lib.Prep("llama_chat_builtin_templates", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeUint32); err != nil {
		return loadError("llama_chat_builtin_templates", err)
	}

	return nil
}

// NewChatMessage creates a new ChatMessage.
func NewChatMessage(role, content string) ChatMessage {
	r, err := utils.BytePtrFromString(role)
	if err != nil {
		return ChatMessage{}
	}

	c, err := utils.BytePtrFromString(content)
	if err != nil {
		return ChatMessage{}
	}

	return ChatMessage{Role: r, Content: c}
}

// ChatApplyTemplate applies a chat template to a slice of [ChatMessage], Set addAssistantPrompt to true to generate the
// assistant prompt, for example on the first message.
func ChatApplyTemplate(template string, chat []ChatMessage, addAssistantPrompt bool, buf []byte) int32 {
	tmpl, err := utils.BytePtrFromString(template)
	if err != nil {
		return 0
	}

	if len(chat) == 0 {
		return 0
	}

	c := unsafe.Pointer(&chat[0])
	nMsg := uint32(len(chat))

	out := unsafe.SliceData(buf)
	len := uint32(len(buf))

	var result ffi.Arg
	chatApplyTemplateFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&tmpl), unsafe.Pointer(&c), &nMsg, &addAssistantPrompt, unsafe.Pointer(&out), &len)
	return int32(result)
}

// ChatBuiltinTemplates returns a list of built-in chat template names.
// The function populates the provided output slice with template names and returns
// the number of templates found. If the output slice is too small, only the first
// len(output) templates will be returned.
func ChatBuiltinTemplates(output []string) int32 {
	if len(output) == 0 {
		return 0
	}

	cStrings := make([]*byte, len(output))
	for i := range cStrings {
		cStrings[i] = nil
	}

	cOutput := unsafe.Pointer(&cStrings[0])
	length := uint32(len(output))

	var result ffi.Arg
	chatBuiltinTemplatesFunc.Call(unsafe.Pointer(&result), &cOutput, &length)

	for i, cStr := range cStrings {
		output[i] = ""
		if cStr != nil {
			output[i] = utils.BytePtrToString(cStr)
		}
	}

	return int32(result)
}
