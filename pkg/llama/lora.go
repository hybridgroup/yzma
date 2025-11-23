package llama

import (
	"errors"
	"unsafe"

	"github.com/hybridgroup/yzma/pkg/utils"
	"github.com/jupiterrider/ffi"
)

var (
	// Load a LoRA adapter from file
	// LLAMA_API struct llama_adapter_lora * llama_adapter_lora_init(
	//         struct llama_model * model,
	//         const char * path_lora);
	adapterLoraInitFunc ffi.Fun

	// Manually free a LoRA adapter
	// NOTE: loaded adapters will be free when the associated model is deleted
	// LLAMA_API void llama_adapter_lora_free(struct llama_adapter_lora * adapter);
	adapterLoraFreeFunc ffi.Fun

	// Get metadata value as a string by key name
	// LLAMA_API int32_t llama_adapter_meta_val_str(const struct llama_adapter_lora * adapter, const char * key, char * buf, size_t buf_size);
	adapterMetaValStrFunc ffi.Fun

	// Get the number of metadata key/value pairs
	// LLAMA_API int32_t llama_adapter_meta_count(const struct llama_adapter_lora * adapter);
	adapterMetaCountFunc ffi.Fun

	// Get metadata key name by index
	// LLAMA_API int32_t llama_adapter_meta_key_by_index(const struct llama_adapter_lora * adapter, int32_t i, char * buf, size_t buf_size);
	adapterMetaKeyByIndexFunc ffi.Fun

	// Get metadata value as a string by index
	// LLAMA_API int32_t llama_adapter_meta_val_str_by_index(const struct llama_adapter_lora * adapter, int32_t i, char * buf, size_t buf_size);
	adapterMetaValStrByIndexFunc ffi.Fun

	// Add a loaded LoRA adapter to given context
	// This will not modify model's weight
	// LLAMA_API int32_t llama_set_adapter_lora(
	//         struct llama_context * ctx,
	//         struct llama_adapter_lora * adapter,
	//         float scale);
	setAdapterLoraFunc ffi.Fun

	// Remove a specific LoRA adapter from given context
	// Return -1 if the adapter is not present in the context
	// LLAMA_API int32_t llama_rm_adapter_lora(
	//         struct llama_context * ctx,
	//         struct llama_adapter_lora * adapter);
	rmAdapterLoraFunc ffi.Fun

	// Remove all LoRA adapters from given context
	// LLAMA_API void llama_clear_adapter_lora(struct llama_context * ctx);
	clearAdapterLoraFunc ffi.Fun

	// LLAMA_API uint64_t            llama_adapter_get_alora_n_invocation_tokens(const struct llama_adapter_lora * adapter);
	adapterGetAloraNInvocationTokensFunc ffi.Fun

	// LLAMA_API const llama_token * llama_adapter_get_alora_invocation_tokens  (const struct llama_adapter_lora * adapter);
	adapterGetAloraInvocationTokensFunc ffi.Fun
)

var (
	errInvalidAdapter = errors.New("invalid LoRA adapter")
)

func loadLoraFuncs(lib ffi.Lib) error {
	var err error

	if adapterLoraInitFunc, err = lib.Prep("llama_adapter_lora_init", &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("llama_adapter_lora_init", err)
	}

	if adapterLoraFreeFunc, err = lib.Prep("llama_adapter_lora_free", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("llama_adapter_lora_free", err)
	}

	if adapterMetaValStrFunc, err = lib.Prep("llama_adapter_meta_val_str", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeUint32); err != nil {
		return loadError("llama_adapter_meta_val_str", err)
	}

	if adapterMetaCountFunc, err = lib.Prep("llama_adapter_meta_count", &ffi.TypeSint32, &ffi.TypePointer); err != nil {
		return loadError("llama_adapter_meta_count", err)
	}

	if adapterMetaKeyByIndexFunc, err = lib.Prep("llama_adapter_meta_key_by_index", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeUint32); err != nil {
		return loadError("llama_adapter_meta_key_by_index", err)
	}

	if adapterMetaValStrByIndexFunc, err = lib.Prep("llama_adapter_meta_val_str_by_index", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypeUint32); err != nil {
		return loadError("llama_adapter_meta_val_str_by_index", err)
	}

	if setAdapterLoraFunc, err = lib.Prep("llama_set_adapter_lora", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer, &ffi.TypeFloat); err != nil {
		return loadError("llama_set_adapter_lora", err)
	}

	if rmAdapterLoraFunc, err = lib.Prep("llama_rm_adapter_lora", &ffi.TypeSint32, &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("llama_rm_adapter_lora", err)
	}

	if clearAdapterLoraFunc, err = lib.Prep("llama_clear_adapter_lora", &ffi.TypeVoid, &ffi.TypePointer); err != nil {
		return loadError("llama_clear_adapter_lora", err)
	}

	if adapterGetAloraNInvocationTokensFunc, err = lib.Prep("llama_adapter_get_alora_n_invocation_tokens", &ffi.TypeUint64, &ffi.TypePointer); err != nil {
		return loadError("llama_adapter_get_alora_n_invocation_tokens", err)
	}

	if adapterGetAloraInvocationTokensFunc, err = lib.Prep("llama_adapter_get_alora_invocation_tokens", &ffi.TypePointer, &ffi.TypePointer); err != nil {
		return loadError("llama_adapter_get_alora_invocation_tokens", err)
	}

	return nil
}

// LoraAdapterInit loads a LoRA adapter from file and applies it to the model.
func AdapterLoraInit(model Model, pathLora string) (AdapterLora, error) {
	var adapter AdapterLora
	if model == 0 {
		return adapter, errors.New("invalid model")
	}

	file := &[]byte(pathLora + "\x00")[0]

	adapterLoraInitFunc.Call(&adapter, unsafe.Pointer(&model), unsafe.Pointer(&file))
	return adapter, nil
}

// AdapterLoraFree manually frees a LoRA adapter.
// Note that loaded adapters will be freed when the associated model is deleted.
func AdapterLoraFree(adapter AdapterLora) error {
	if adapter == 0 {
		return errInvalidAdapter
	}

	adapterLoraFreeFunc.Call(nil, unsafe.Pointer(&adapter))
	return nil
}

// AdapterMetaValStr gets metadata value as a string by key name.
func AdapterMetaValStr(adapter AdapterLora, key string) (string, bool) {
	if adapter == 0 {
		return "", false
	}
	buf := make([]byte, 8192)
	b := unsafe.SliceData(buf)
	bLen := int32(len(buf))

	keyPtr, _ := utils.BytePtrFromString(key)
	var result ffi.Arg
	adapterMetaValStrFunc.Call(
		unsafe.Pointer(&result),
		unsafe.Pointer(&adapter),
		unsafe.Pointer(&keyPtr),
		unsafe.Pointer(&b),
		&bLen,
	)
	if int32(result) < 0 {
		return "", false
	}

	// copy to a new slice to avoid retaining the entire buffer
	value := make([]byte, int32(result))
	copy(value, buf[:int32(result)])

	return string(value), true
}

// AdapterMetaCount gets the number of metadata key/value pairs.
func AdapterMetaCount(adapter AdapterLora) int32 {
	if adapter == 0 {
		return 0
	}
	var result ffi.Arg
	adapterMetaCountFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&adapter))
	return int32(result)
}

// AdapterMetaKeyByIndex gets metadata key name by index.
func AdapterMetaKeyByIndex(adapter AdapterLora, i int32) (string, bool) {
	if adapter == 0 {
		return "", false
	}
	buf := make([]byte, 128)
	b := unsafe.SliceData(buf)
	bLen := int32(len(buf))

	var result ffi.Arg
	adapterMetaKeyByIndexFunc.Call(
		unsafe.Pointer(&result),
		unsafe.Pointer(&adapter),
		&i,
		unsafe.Pointer(&b),
		&bLen)
	if int32(result) < 0 {
		return "", false
	}

	// copy to a new slice to avoid retaining the entire buffer
	value := make([]byte, int32(result))
	copy(value, buf[:int32(result)])

	return string(value), true
}

// AdapterMetaValStrByIndex gets metadata value as a string by index.
func AdapterMetaValStrByIndex(adapter AdapterLora, i int32) (string, bool) {
	if adapter == 0 {
		return "", false
	}
	buf := make([]byte, 8192)
	b := unsafe.SliceData(buf)
	bLen := int32(len(buf))

	var result ffi.Arg
	adapterMetaValStrByIndexFunc.Call(
		unsafe.Pointer(&result),
		unsafe.Pointer(&adapter),
		&i,
		unsafe.Pointer(&b),
		&bLen)
	if int32(result) < 0 {
		return "", false
	}

	// copy to a new slice to avoid retaining the entire buffer
	value := make([]byte, int32(result))
	copy(value, buf[:int32(result)])

	return string(value), true
}

// SetAdapterLora adds a loaded LoRA adapter to the given context.
// Returns 0 on success, or a negative value on failure.
func SetAdapterLora(ctx Context, adapter AdapterLora, scale float32) int32 {
	if ctx == 0 {
		return -1
	}
	if adapter == 0 {
		return -1
	}

	var result ffi.Arg
	setAdapterLoraFunc.Call(
		unsafe.Pointer(&result),
		unsafe.Pointer(&ctx),
		unsafe.Pointer(&adapter),
		unsafe.Pointer(&scale),
	)
	return int32(result)
}

// RmAdapterLora removes a specific LoRA adapter from the given context.
// Returns 0 on success, or -1 if the adapter is not present.
func RmAdapterLora(ctx Context, adapter AdapterLora) int32 {
	if ctx == 0 {
		return -1
	}
	if adapter == 0 {
		return -1
	}

	var result ffi.Arg
	rmAdapterLoraFunc.Call(
		unsafe.Pointer(&result),
		unsafe.Pointer(&ctx),
		unsafe.Pointer(&adapter),
	)
	return int32(result)
}

// ClearAdapterLora removes all LoRA adapters from the given context.
func ClearAdapterLora(ctx Context) {
	if ctx == 0 {
		return
	}
	clearAdapterLoraFunc.Call(nil, unsafe.Pointer(&ctx))
}

// AdapterGetAloraNInvocationTokens returns the number of invocation tokens for the adapter.
func AdapterGetAloraNInvocationTokens(adapter AdapterLora) uint64 {
	if adapter == 0 {
		return 0
	}
	var result ffi.Arg
	adapterGetAloraNInvocationTokensFunc.Call(unsafe.Pointer(&result), unsafe.Pointer(&adapter))
	return uint64(result)
}

// AdapterGetAloraInvocationTokens returns a slice of invocation tokens for the adapter.
func AdapterGetAloraInvocationTokens(adapter AdapterLora) []Token {
	n := AdapterGetAloraNInvocationTokens(adapter)
	if n == 0 {
		return nil
	}

	var ptr *Token
	adapterGetAloraInvocationTokensFunc.Call(unsafe.Pointer(&ptr), unsafe.Pointer(&adapter))
	if ptr == nil {
		return nil
	}

	return unsafe.Slice(ptr, n)
}
