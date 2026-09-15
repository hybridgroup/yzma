// Package speculative provides Go FFI bindings for the NextN hidden state
// functions that MTP (Multi-Token Prediction) speculative decoding needs.
//
// These functions come from src/llama-ext.h, which llama.cpp calls a staging
// header. Breaking changes are permitted there, so this API can change or go
// away with a llama.cpp update.
//
// The declarations in that header are outside extern "C", so a shared library
// can export them under a C++ mangled name only. [Load] tries the C name and
// then the Itanium C++ ABI name, which covers clang and gcc on Linux and
// macOS. MSVC uses a different scheme, so [Available] can report false for a
// Windows build of llama.cpp.
//
// Upstream renamed these functions once, from pre_norm to nextn. This package
// binds the current nextn spelling only.
package speculative
