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

// Struct-by-value, in the shape of yzma#289. The cif descriptor for this struct
// omits n_b, and the four bytes are reabsorbed by the alignment padding in
// front of ud, so the descriptor is still 24 bytes: a size-only comparison
// cannot see it, and every member from n_c onwards is read four bytes early.
struct fx_params {
    uint32_t n_a;
    uint32_t n_b;
    uint32_t n_c;
    float    scale;
    void *   ud;
};

LLAMA_API struct fx_params fx_params_default(void);

// Struct-by-value where the descriptor is right and the Go struct is wrong:
// same 24 bytes, but the Go side orders the members differently.
struct fx_geom {
    uint32_t w;
    uint32_t h;
    float    s;
    void *   ud;
};

LLAMA_API struct fx_geom fx_geom_default(void);

// Clean struct-by-value control: C struct, cif descriptor and Go struct all
// agree member for member.
LLAMA_API void fx_use_geom(struct fx_geom g);

// Second clean control, in the shape of llama.Batch: the Go struct appends its
// own bookkeeping after the FFI members, past the end of the C struct, so
// libffi never touches it.
LLAMA_API void fx_use_geom_tail(struct fx_geom g);

// The other side of that tolerance: a Go struct shorter than the C struct is a
// read overrun, and must stay a RULE 2 finding.
LLAMA_API void fx_use_geom_short(struct fx_geom g);

// Two members of the same width and ABI class. Swapping them on the Go side
// moves no offset, changes no size and changes no class, so only the member
// names can catch it -- and the C library receives each value in the other's
// place.
struct fx_pair {
    uint32_t alpha_count;
    uint32_t beta_count;
};

LLAMA_API void fx_use_pair(struct fx_pair p);

// Clean control for the alias table: every member here is named differently on
// the Go side by more than case and underscores (n_ dropped, attn expanded,
// ctx expanded), and none of it is a transposition.
struct fx_named {
    uint32_t n_attn_heads;
    uint32_t ctx_size;
};

LLAMA_API void fx_use_named(struct fx_named n);

// RULE 4. LLAMA_FX_LEVEL_LOW is the member "upstream inserted": every member
// after it moved up by one, and the Go constant for HIGH still carries the old
// value. Nothing about that is visible to a compiler on either side.
//
// The members have no initialisers, so this also pins the implicit counter, and
// MAX is initialised from another member of the same enum.
enum llama_fx_level {
    LLAMA_FX_LEVEL_OFF,
    LLAMA_FX_LEVEL_LOW,
    LLAMA_FX_LEVEL_HIGH,
    LLAMA_FX_LEVEL_MAX = LLAMA_FX_LEVEL_HIGH,
};

// Clean controls for the same rule, one per initialiser form the parser has to
// evaluate: a shift, a hex literal, and a negative sentinel.
enum llama_fx_flag {
    LLAMA_FX_FLAG_AUTO = -1,
    LLAMA_FX_FLAG_NONE = 0,
    LLAMA_FX_FLAG_A    = 1 << 2,
    LLAMA_FX_FLAG_B    = 0x10,
};

// RULE 5, the direction where C calls Go. Nothing below is bound with
// lib.Prep: a closure is described either by an ffi.PrepCif descriptor or by
// the Go signature purego.NewCallback is handed, and libffi unpacks the C stack
// through it on every invocation.
//
// llama_fx_progress_callback is the clean control for the descriptor form and
// llama_fx_report_callback the plant: same signature, but the descriptor for it
// declares its float parameter as a 4-byte int, so the closure reinterprets the
// progress value's bit pattern as an integer.
typedef bool (*llama_fx_progress_callback)(float progress, void * user_data);
typedef bool (*llama_fx_report_callback)(float progress, void * user_data);

// llama_fx_log_callback is the plant for the purego form: the Go closure for it
// is one parameter short, so text arrives in level's place and the string
// pointer is read out of a register C never set.
typedef void (*llama_fx_log_callback)(enum llama_fx_level level, const char * text, void * user_data);

// llama_fx_abort_callback is the clean control for that form. C reads one byte
// of the register the closure returns, so a word-sized Go result is right.
typedef bool (*llama_fx_abort_callback)(void * data);

// Integer #defines, including one initialised from another and one carrying the
// unsigned literal suffix llama.h uses on its file magics.
#define LLAMA_FX_MAGIC   0x66780001u
#define LLAMA_FX_VERSION 3
#define LLAMA_FX_MAGIC_ALIAS LLAMA_FX_MAGIC
