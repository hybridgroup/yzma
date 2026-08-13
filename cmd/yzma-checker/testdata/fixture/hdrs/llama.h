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
//
// It is also the deprecation-inventory control, wrapped exactly as upstream
// wraps llama_free_model. A deprecated declaration has the same ABI as any
// other, so it must be parsed as if the macro were not there and must appear in
// the inventory without becoming a violation - which is the half that matters,
// and is pinned by fx_clean's place on the clean list.
DEPRECATED(LLAMA_API int32_t fx_clean(struct fx_thing * t, int32_t a, size_t n),
    "use nothing, this is a fixture");

// The variadic pair, for the RULE 1 nfixed check. Both declare two parameters
// before the "...", and on Apple arm64 that boundary decides whether an
// argument travels in a register or on the stack: fx_logf is bound with nfixed
// 1, so every argument from fmt onwards is looked for in the wrong place, while
// fx_printf declares the same shape and is bound correctly.
//
// fx_printf is also wrapped in DEPRECATED(...), as the control for the
// deprecation note on the Go side: its wrapper carries the `Deprecated:`
// paragraph that fx_clean's does not.
LLAMA_API void fx_logf(struct fx_thing * t, const char * fmt, ...);
DEPRECATED(LLAMA_API int32_t fx_printf(struct fx_thing * t, const char * fmt, ...),
    "use nothing, this is a fixture");

// The coverage-inventory control: nothing binds this, which is not a defect.
// It must appear in the unbound inventory and must never become a violation.
LLAMA_API int32_t fx_unbound(struct fx_thing * t, int32_t which);

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

// Pointer-out params, the shape of llama_get_embeddings_seq and
// llama_get_logits_ith. Every one of these is an 8-byte pointer to the ABI, so
// the slot itself cannot carry the defect: only what it points *at* can.
//
// fx_get_scores is the plant, bound to a Go *int32, so C writes IEEE-754 bit
// patterns into memory Go reads as integers. fx_get_logits is the control, bound
// to the *float32 it declares.
LLAMA_API void fx_get_scores(struct fx_thing * t, float * out);
LLAMA_API void fx_get_logits(struct fx_thing * t, float * out);

// The signedness pair, one indirection down. Neither is an ABI break - a 4-byte
// signed and unsigned integer are the same register and the same bytes - so
// fx_get_count must be reported as a signedness finding rather than a violation,
// and fx_get_token, bound to the *int32 that can hold its -1 sentinel, must not
// be reported at all.
LLAMA_API void fx_get_count(struct fx_thing * t, int32_t * out);
LLAMA_API void fx_get_token(struct fx_thing * t, int32_t * out);

// The same idea on a return buffer, which is where a value is interpreted rather
// than forwarded: fx_decode returns negative on error, exactly as llama_decode
// does, so reading it through an unsigned Go buffer turns a failure into
// 4294967295. fx_decode_ok reads the identical declaration through a signed
// buffer of the same width and must stay clean.
LLAMA_API int32_t fx_decode(struct fx_thing * t);
LLAMA_API int32_t fx_decode_ok(struct fx_thing * t);

// The C strings. Every one of these is an 8-byte pointer to a 1-byte char on
// both sides, so nothing about the slot, the width or the pointer target can
// carry the defect: what C depends on is the bytes ending in a NUL, and Go never
// puts one there.
//
// fx_set_name is the plant, fed a []byte conversion of a Go string with no
// terminator appended, so C reads forward off the end of that allocation.
// fx_set_path and fx_set_text are the controls, one per idiom yzma actually
// uses: the appended "\x00" before the []byte conversion, and the string
// terminated in place before unsafe.StringData.
LLAMA_API void fx_set_name(struct fx_thing * t, const char * name);
LLAMA_API void fx_set_path(struct fx_thing * t, const char * path);
LLAMA_API void fx_set_text(struct fx_thing * t, const char * text);

// The enum parameters. All three slots are &ffi.TypeSint32 against a 4-byte C
// enum on the descriptor side and a 4-byte Go int32 on the value side, so every
// rule passes all of them and only the enum a value belongs to separates them.
//
// fx_set_level is the plant, fed the enum that mirrors llama_fx_split_mode, so
// the library is asked for a split mode where a level belongs. fx_set_flag is the
// control, fed the type that mirrors its own enum. fx_set_mode is the
// out-of-scope control: a plain int32 names no enum to compare against, exactly
// as a void * names no pointer target, so it is neither a finding nor a skip.
LLAMA_API void fx_set_level(struct fx_thing * t, enum llama_fx_level level);
LLAMA_API void fx_set_flag(struct fx_thing * t, enum llama_fx_flag flag);
LLAMA_API void fx_set_mode(struct fx_thing * t, enum llama_fx_split_mode mode);

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

// The partially-mirrored-enum control for the inventory. Only NONE has a Go
// constant, so LAYER and TENSOR are members RULE 4 can never compare - it walks
// the Go side. That is not a defect, exactly as an unbound declaration is not:
// it is the signal that this enum has members yzma has not caught up with, which
// is the same event as an insertion shifting the values of the ones it mirrors.
// The assertion that matters is that these two never become violations.
enum llama_fx_split_mode {
    LLAMA_FX_SPLIT_MODE_NONE,
    LLAMA_FX_SPLIT_MODE_LAYER,
    LLAMA_FX_SPLIT_MODE_TENSOR,
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

// RULE 5 through a struct member, the shape of llama_context_params.cb_eval and
// llama_model_params.progress_callback. cb_progress and cb_abort are function
// pointers C jumps *through* rather than values it reads, and to a layout
// comparison each is 8 bytes of pointer class like user_data next to it. So
// neither the signature C calls them with nor the identity of the code stored in
// them is visible to any offset, width or class comparison.
struct fx_hooks {
    uint32_t                   n_slots;
    llama_fx_progress_callback cb_progress;
    llama_fx_abort_callback    cb_abort;
    void *                     user_data;
};

// fx_use_hooks takes the Go struct whose cb_progress is declared as a Go func
// value: 8 bytes of pointer class, and a pointer to a func descriptor rather
// than to code.
LLAMA_API void fx_use_hooks(struct fx_hooks h);

// fx_use_hooks_ok takes the all-uintptr Go struct, one member of which is given
// the code pointer of a closure built for the wrong typedef.
LLAMA_API void fx_use_hooks_ok(struct fx_hooks h);

// fx_use_hooks_clean is the control: both members hold the code pointer of the
// callback that implements them, one per callback form, and must never be
// reported.
LLAMA_API void fx_use_hooks_clean(struct fx_hooks h);

// Integer #defines, including one initialised from another and one carrying the
// unsigned literal suffix llama.h uses on its file magics.
#define LLAMA_FX_MAGIC   0x66780001u
#define LLAMA_FX_VERSION 3
#define LLAMA_FX_MAGIC_ALIAS LLAMA_FX_MAGIC
