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
)

func loadChatFuncs(lib ffi.Lib) error {
	var err error
	if chatApplyTemplateFunc, err = lib.Prep("llama_chat_apply_template", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32,
		&ffi.TypeUint8, &ffi.TypePointer, &ffi.TypeSint32); err != nil {

		return loadError("llama_chat_apply_template", err)
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
