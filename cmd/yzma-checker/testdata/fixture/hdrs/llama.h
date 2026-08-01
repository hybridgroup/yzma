// Fixture header for the checker's self-test. Not a real llama.cpp header.
//
// Each declaration below exists so that a matching Go binding in
// ../bindings/bindings.go either violates exactly one rule or is clean.

// LLAMA_API is deliberately left undefined: the parser scans for the macro
// token, so a #define line would itself be picked up as a declaration.

// Pointer return. Binding it with a void return descriptor is RULE 1.
LLAMA_API struct fx_thing * fx_get_thing(void);

// size_t buf_size. Passing the address of a 4-byte Go value is RULE 2.
LLAMA_API int32_t fx_desc(struct fx_thing * t, char * buf, size_t buf_size);

// float return. Reading it through an ffi.Arg is RULE 3.
LLAMA_API float fx_score(struct fx_thing * t, int32_t token);

// int return. A 4-byte return buffer is RULE 3 in the other direction.
LLAMA_API int fx_mode_from_str(const char * s);

// Fully clean control: every width matches on both sides.
LLAMA_API int32_t fx_clean(struct fx_thing * t, int32_t a, size_t n);
