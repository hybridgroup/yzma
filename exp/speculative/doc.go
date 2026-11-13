// Package speculative provides Go FFI bindings for the NextN hidden state
// functions that MTP (Multi-Token Prediction) speculative decoding needs.
//
// llama.cpp moved these functions from the src/llama-ext.h staging header to
// llama.h, so newer builds export them under their C names.
//
// Older builds declare them outside extern "C" and export them under a mangled
// C++ name only. [Load] tries the C name and then the Itanium C++ ABI name,
// which covers clang and gcc on Linux and macOS. MSVC uses a different scheme,
// so [Available] can report false for an older Windows build of llama.cpp.
//
// Upstream renamed these functions once, from pre_norm to nextn. This package
// binds the current nextn spelling only. It binds llama_set_ctx_state_src and
// falls back to its old llama_set_embd_state_src name.
//
// The device row functions keep the target NextN rows on the GPU so an MTP
// draft context can read them directly. They came later than the host
// functions, so check [DeviceRowsAvailable] separately from [Available], and
// keep the host path as a fallback.
package speculative
