# yzma-checker

Mechanical verification of the **FFI parameter and return types** in every yzma
binding, and of the **values of the constants** yzma mirrors, against the
llama.cpp headers for the build yzma installs.

yzma binds llama.cpp through purego/libffi, not cgo. There is no `cgocheck`, no
automatic pinning, and no compiler validation of argument widths, so a mismatch
between

1. the C declaration,
2. the libffi type descriptor passed to `ffi.PrepCif`, and
3. the Go variable whose address is passed to `ffi.Fun.Call`

is completely silent at build time and only shows up at runtime as a corrupted
value or corrupted memory. The same is true in the other direction, where C
calls a Go closure back. 265 bindings, 169 mirrored constants, 4 callbacks and
the 5 function-pointer struct members C reaches them through are far too many to
eyeball credibly, which is why this exists.

## Running it

```sh
cd cmd/yzma-checker
go run .          # add -v to dump every binding with its C signature
```

No arguments needed, and nothing outside this repository is assumed. It works
out:

| what | how |
|---|---|
| yzma tree | walks up for the `github.com/hybridgroup/yzma` `go.mod` |
| llama.cpp ref | `tag_name` from `llama-cpp-builder`, the same endpoint yzma's installer uses |
| headers | fetched from `ggml-org/llama.cpp` at that ref, cached under `$XDG_CACHE_HOME` |

The banner prints all three, so a run is self-documenting. Only the first run
for a given ref needs the network; after that the header cache serves it.

It exits non-zero when it reports a violation, so it can be wired into CI as-is.

### Auditing a consumer's yzma

Run it from another module's root and it audits the yzma **that module actually
compiles against**, resolved with `go list -m`:

```
$ cd ../some-project && yzma-checker
yzma:   /Users/me/go/pkg/mod/github.com/hybridgroup/yzma@v1.21.0 (module cache (via go list))
```

This distinction has caused real confusion: a locally-patched yzma checkout that
is not referenced by a `replace` directive or `go.work` is not compiled, so
auditing it describes code that does not run.

### Flags

```
-yzma   <dir>   yzma tree to audit (default: auto-detect, as above)
-ref    <ref>   llama.cpp git ref, or "auto" to resolve the current release
-llama  <dir>   read headers from a llama.cpp git checkout via `git show <ref>:`
                instead of the network; the checkout's own worktree state is
                irrelevant, only the ref matters
-hdrs   <dir>   use pre-extracted headers; the dir must contain llama.h, ggml.h,
                ggml-backend.h, ggml-cpu.h, mtmd.h and mtmd-helper.h
-pkgs   <list>  comma-separated package patterns to audit
-v              dump every binding: cif types, C signature, and every call site,
                plus every constant and the C constant it was matched to, and
                every callback with the C typedef it was linked to, and every
                unbound C declaration and unmirrored C enum member
```

`-llama` and `-hdrs` are the offline paths.

## The five rules

**RULE 1 — the cif must match the C prototype.** Parameter for parameter, in
width and in class. On arm64: `size_t`/`uint64_t`/any pointer = 8 bytes,
`int`/`int32_t`/`enum`/`float` = 4, `double` = 8. A struct passed by value has
its own descriptor and is a common error site, so its members are compared one
by one — see *Struct layouts*, below.

A **variadic** binding is checked here too, and needs one number the others do
not. `lib.PrepVar("printf", 1, ret, args...)` says the first argument is fixed
and everything after it is variadic, and on Apple arm64 those are two different
calling conventions rather than two spellings of one: a variadic argument is
passed on the stack where a fixed one goes in a register. So a `nfixed` that is
off by one misplaces *every* variadic argument while every type in the descriptor
is still correct — which is why it is compared against the number of parameters
the C prototype declares in front of its `...`, separately from the list it
indexes into. The descriptor may be *longer* than `nfixed`, since a `PrepVar` cif
lists the fixed types plus the concrete variadic types of the call it describes,
but never shorter: that would mean libffi is told a fixed argument exists that it
has no type for. And a variadic C function bound with `Prep`, or a fixed one
bound with `PrepVar`, is an ABI break in itself.

Variadic bindings used to be exempt from arity checking altogether — the arity
comparison was skipped for any variadic prototype and nothing else looked at
`nfixed`, so they were the one class of binding RULE 1 never compared against a C
prototype at all. **yzma has no variadic bindings today**, so this check finds
nothing in the tree as this is written; it exists so the first one added is
verified rather than assumed, and the fixture is where it is demonstrated to
work.

**RULE 2 — the Go variable must be exactly as wide as the cif claims.** libffi
reads exactly `ffi.Type.Size` bytes from each pointer it is handed. `&x` where
`x` is a Go `int32` behind an `ffi.TypeUint64` slot makes libffi splice 4 bytes
of adjacent Go memory into the high half of the C argument.

A pointer slot needs one more comparison than a width, because every pointer is
8 bytes: what it points *at* is compared too, and a `float *` fed a `*int32` is a
RULE 2 finding — see *Pointer targets*, below. A `char *` slot needs a third,
because C finds the end of a string in the bytes and Go never puts a terminator
there — see *C string termination*.

For a **struct** passed by value the C-declared size is the authority in one
direction only: since libffi reaches exactly that far, a Go struct that appends
its own fields *past* the end of the C struct — as `llama.Batch` does with
`capTokens`/`capSeq` — is bytes libffi never touches, and is not a finding. A Go
struct that is short of the C size, or whose extra members land inside the
region libffi does reach, still is. Scalars stay exact in both directions.

**RULE 3 — return buffers must match the return *kind*.**
- Integer/pointer returns: the buffer must be **8 bytes** (`ffi.Arg`), because
  libffi always stores a full `ffi_arg`. A narrower variable is written past.
- Float returns are the **opposite**: they must be `float32`. libffi stores only
  4 bytes and leaves the rest of the buffer alone, so an `ffi.Arg` ends up
  holding the IEEE-754 bit pattern in its low word — and `float32(someArg)` is a
  numeric conversion, not a bit reinterpretation.
- A `TypeVoid` return descriptor makes libffi write **nothing at all**, so a
  supplied buffer silently keeps its zero value.

**RULE 4 — a mirrored constant must still hold its C value.** The first three
rules check the *shape* of every call and never look at a single value. But yzma
also transcribes every llama.cpp enum member and every interesting `#define`
into a Go constant, by hand:

```go
// GGML_TYPE_F32
GGMLTypeF32 GGMLType = 0

PoolingTypeMean PoolingType = 1     // no comment; mirrors LLAMA_POOLING_TYPE_MEAN
```

Nothing declares those numbers. `LLAMA_POOLING_TYPE_MEAN` is 1 because of where
it sits in its enum, so when upstream inserts a member every later value shifts
by one and the transcription is quietly wrong. Both sides still compile, every
width still matches, every offset still matches, and the library is simply asked
for a different pooling type than the caller named — the same silent drift as
the struct-field bug in [#289](https://github.com/hybridgroup/yzma/issues/289),
one layer up. It is live in the tree as this is written: `GGMLTypeCOUNT` is 40
while `GGML_TYPE_COUNT` has become 43, because llama.cpp added NVFP4, Q1_0 and
Q2_0.

The C side is read out of the same headers as everything else: enum members with
their implicit counter, explicit decimal, hex and negative initialisers, members
initialised from an earlier member (`LLAMA_ROPE_SCALING_TYPE_MAX_VALUE =
LLAMA_ROPE_SCALING_TYPE_LONGROPE`), shift expressions (`1 << 4`), and integer
`#define`s. Values on the Go side come from `go/types`, so `iota`, `1 << 6` and
`FileMagicGGSN` arrive as the integer the compiler computed rather than as text
to re-parse.

Matching the two is where the check can go wrong, so it is done twice over and
never guessed:

| evidence | example |
|---|---|
| a bare C name in the comment directly above the constant | `// GGML_TYPE_F32` above `GGMLTypeF32` |
| the names, with case, underscores and an optional `llama`/`ggml`/`mtmd` prefix removed | `PoolingTypeMean` ↔ `LLAMA_POOLING_TYPE_MEAN` |

A Go constant matching **no** C name, or matching **several**, is reported as a
skip. Both cases really happen: `llama.h` and `ggml.h` each declare a
quantisation enum whose members normalise onto each other but do not hold the
same values, so `FtypeMostlyQ4_0` alone is ambiguous and `constEnumTags` records
that yzma's `Ftype` is the `llama_ftype` one. Where yzma spells out a word the
header abbreviates there is `constAliases` (`Continue`/`CONT`,
`XTCThold`/`XTC_THRESHOLD`, `UserDef`/`USER_DEFINED`), and where a constant is
yzma's own invention with nothing upstream to compare it to — the `GpuBackend`
enumeration, the `SamplerType` dispatch tags, `MaxToken` — there is
`goOnlyConsts`. All three are closed hand-checked lists for the same reason
`memberAliases` is: a rule loose enough to guess a name is loose enough to
validate a constant against the wrong one, and a new divergence should cost one
reviewed line rather than pass silently.

**RULE 5 — a callback must match the function-pointer typedef C will call it
through.** Rules 1-3 check calls yzma makes *into* C. A callback is the reverse:
C puts the arguments on the stack and libffi — or purego — hands them to a Go
closure. Nothing in that direction goes through `lib.Prep`, so the first three
rules never see it, and the failure mode is worse rather than milder. A wrong
parameter width does not corrupt one argument of one call: it makes the closure
read *C stack memory* instead of the value C passed, on every single invocation
for the life of the process, and a return that is narrower than what C reads
back makes C use a partly-uninitialised register.

There are four such sites, in two forms — plus the struct members C reaches a
callback through, below:

```go
// the descriptor form: llama_progress_callback / mtmd_progress_callback
ffi.PrepCif(progressCallbackCif, ffi.DefaultAbi, 2, &ffi.TypeUint8, &ffi.TypeFloat, &ffi.TypePointer)
ffi.PrepClosureLoc(closure, progressCallbackCif, fn, nil, progressCallbackCode)

// the purego form: ggml_log_callback / ggml_abort_callback
purego.NewCallback(func(level int32, text, data uintptr) uintptr { ... })
```

For the descriptor form the cif is compared against the typedef exactly as RULE
1 compares a binding's cif against a prototype, plus the `nfixed` argument
against the typedef's parameter count — `nfixed` is what libffi *believes* about
the argument count, so it is checked separately from the list it indexes into.
`&ffi.TypeSint32` where the typedef says `float` is the shape of the bug: the
width is right, so libffi reads the correct four bytes, and the closure then
reads a float's bit pattern as an integer.

For the purego form the Go func literal's own signature *is* the descriptor, so
it is read out of `go/types` and compared parameter by parameter in width and
ABI class. purego carries every argument in a register, so a Go parameter it
cannot put in one is a finding too. A parameter *count* mismatch is reported as
loudly as a width mismatch, because a missing parameter shifts every later one:
a two-parameter closure behind `ggml_log_callback` receives the text pointer
where the level belongs. Returns are asymmetric here. A Go result wider than C
reads is harmless — C reads the low bytes of the register the closure filled,
which is why `func(data uintptr) uintptr` is right for a `bool`-returning
`ggml_abort_callback` — and so is a result C ignores entirely, which is how the
`void`-returning log callback is written. A result *narrower* than C reads, or of
the wrong ABI class, is not.

Which typedef a site belongs to is the part that can go wrong, so it is never
guessed. `callbackTags` names it outright where yzma's identifier does not say
(`LogSilent` is named for what it does, not for what it implements); otherwise
the identifier is normalised, dropping the affixes yzma adds around the C name —
`progressCallbackCif` → `progress_callback`, `newAbortCallback` →
`abort_callback` — and matched against every function-pointer typedef with its
`llama_`/`ggml_`/`mtmd_` prefix removed. Exactly one match is a link; two
candidates are broken by the prefix belonging to the Go package, which is what
separates `mtmd_progress_callback` in `pkg/mtmd` from `llama_progress_callback`
in `pkg/llama`. Anything still unresolved, and any typedef the parser could not
take apart, is a **skip** naming the reason — never a silent pass, because a
callback that looks checked and is not is worse than one that was never claimed.

#### Function-pointer struct members

Those two forms are where a callback is *built*. A C struct member of
function-pointer type is where one is **installed**, and it used to be the form
nothing looked at:

```c
struct llama_context_params {
    ggml_backend_sched_eval_callback cb_eval;      // C jumps through this
    ggml_abort_callback              abort_callback;
    void *                           abort_callback_data;   // C reads this
};
```

`cb_eval` and `abort_callback_data` are indistinguishable to `Leaf`: both are 8
bytes of pointer class at some offset. So the signature C will call the member
*with* was compared against nothing, and neither was what gets stored in it.
There are five such members in the structs yzma binds by value —
`llama_model_params.progress_callback`, `llama_context_params.cb_eval` and
`.abort_callback`, and `mtmd_context_params`' two.

Two things were therefore silent. A Go `func` field is also 8 bytes of pointer
class, and it is a pointer to a func **descriptor** rather than to code, so C
jumping through it lands on whatever the descriptor's first word happens to be —
and every offset, width and class comparison passes it. And where the field is
the `uintptr` yzma actually declares, nothing in the Go types says *which*
callback belongs in it: storing the log callback's code pointer into `cb_eval`
type-checks, lays out identically, and makes C unpack a
`ggml_backend_sched_eval_callback`'s arguments through a `ggml_log_callback`
descriptor on every graph node.

So the offset correspondence comes from the same flattening the layout diff uses
— the C member's offset is `CStruct.Offs`, and the Go leaf sitting there is the
field libffi copies the bytes of, so there is one struct walker rather than two —
and the code pointer stored in that field is traced back to the site that
produced it: `p.ProgressCallback = uintptr(progressCallbackCode)` through
`ffi.PrepClosureLoc(closure, progressCallbackCif, ...)` to the cif, and from the
cif to the typedef the link above resolved. That typedef must be the member's
own. Combined with the descriptor comparison, that is what makes the signature C
calls through the member *checked* rather than assumed.

The scope is deliberately narrower than the member count, and narrowed rather
than filled with skips. A member left at zero, or never assigned at all — which
is the case for `cb_eval` and `abort_callback`, since yzma has no setter for
either — has no code pointer to identify, exactly as a `void *` has no pointer
target: outside the claim rather than a skip. Both numbers are on the `SUMMARY`
as `fn-ptr members compared` and `code pointers traced`, so a narrowing of either
is visible rather than silent. A member whose typedef the parser could not take
apart *is* a skip, because that one looks checked and is not.

### Struct layouts

Struct-by-value slots used to be compared by total size alone, which is the
weakest check the tool made and the one that failed in the field
([#289](https://github.com/hybridgroup/yzma/issues/289)). llama.cpp b10375
inserted a `uint32_t` into `llama_context_params`; neither `ContextParams` nor
`ffiTypeContextParams` gained it, and the four missing bytes were reabsorbed by
the alignment padding in front of `cb_eval`. Both sides really were 160 bytes,
while every member after the insertion point was read four bytes early — so
`llama_get_embeddings_seq` returned NULL and graph computation failed, silently,
with the checker green.

Any struct mixing 4-byte and 8-byte members has interior padding, so this is not
a contrived shape: it is the shape of `llama_context_params`,
`llama_model_params` and `llama_batch`, the three structs upstream churns most.
The checker was therefore at its weakest exactly where llama.cpp moves fastest,
and on the failure mode with the worst symptom — a size mismatch is loud, a
displaced member is not.

So each struct-by-value slot now has all three of its representations flattened
to a common member list — byte offset, width and ABI class per primitive member,
arrays expanded and nested structs inlined, padding never emitted — and the
lists are compared element by element:

| representation | comes from | compared against |
|---|---|---|
| C struct | the headers | the cif descriptor (RULE 1) |
| cif descriptor | `ffi.NewType(...)` | both others |
| Go struct | `go/types`, arm64 layout | the cif descriptor (RULES 2 and 3), as a prefix; and the C struct by member name |

Both directions matter, because the Go struct and the descriptor are written by
hand in two separate places: in #289 they had drifted *together*, so only the
comparison against C could see it. A reordered Go struct against a correct
descriptor is the other direction, and the one the README used to list as
uncovered.

### Transposed members

Offsets cannot see everything. Swapping two members of the same width and ABI
class — two `int32_t`s, or an `int32_t` with a `uint32_t` — moves no offset,
changes no size and changes no class, so the layout comparison above passes it
and the call itself succeeds. The C library then simply receives each value in
the other's place, which is a wrong-data bug rather than a memory bug, and just
as silent.

The member *names* can see it, and only the header has names to compare against:
a cif descriptor carries none. So the Go struct's members are matched to the C
struct's by name, and the check keys on a **permutation** — a Go name that
matches a C member at a *different* index is a transposition. A Go name matching
nothing at all is not evidence of anything, so it is reported as `NOT VERIFIED`
rather than as a defect.

Names are compared with case and underscores removed (`NThreadsBatch` and
`n_threads_batch` are the same member), plus a short table in `layout.go` for
the places yzma diverges by more than that: `attn`/`Attention`, `ctx`/`Context`,
`t_p_eval`/`TPromptEval`, `tt_overrides`/`TensorTypes`, and an optional leading
`n_`. That table is deliberately a closed list rather than a fuzzy matcher — an
alias loose enough to guess is loose enough to accept a real mismatch — so a new
divergence appears as a skip that needs one line added, never as a silent hole.

Only the first few differing members are reported: one displaced member shifts
every member after it, and printing all of them buries the insertion point that
is the actual finding. `-v` adds the full member list of both sides to the
`STRUCT-BY-VALUE COMPARISONS` section.

### Pointer targets

Offsets and widths cannot see one indirection down either. `classify()` reduces
every C `T *` to `KindPointer`/8 bytes, which is all the ABI needs, and
`kindCompat` merges a pointer with an 8-byte integer because on arm64 they are
the same register. So a C `float *` slot fed a Go `*int32` matches in width,
matches in ABI class and passes all five rules: the call succeeds, and C then
writes IEEE-754 bit patterns into memory Go reads as integers. That is the
wrong-data class of [#289](https://github.com/hybridgroup/yzma/issues/289) again,
one level of indirection down, and it lands where llama.cpp is busiest, because
its hottest entry points are pointer-out params — `llama_get_embeddings_seq`,
`llama_get_logits_ith`, the `llama_batch` members.

The data to compare is already there: the C parameter keeps its type text
(`const float *`), so stripping one trailing `*` and classifying the remainder
gives what C points at, and the Go type libffi reads the address from is itself a
Go pointer, so its element type gives what Go points at. Width and ABI class are
then compared exactly as RULE 2 compares the slot itself, and a mismatch is a
RULE 2 finding.

The comparison is only made where **both sides name a concrete scalar target**. A
`void *`, an opaque `struct llama_context *`, a Go `uintptr` or a Go
`unsafe.Pointer` names nothing to compare against, so those slots are outside
what this can decide rather than skips — a skip means "this should have been
checked and was not", and a `void *` is not that. How many slots are inside that
scope is on the `SUMMARY` as `pointer targets compared`, so a narrowing of it is
visible rather than silent.

Pointer *members inside a by-value struct* are **not** covered. `Leaf` carries
offset, width and ABI class, and the cif descriptor has no target information at
all — `TypePointer` says nothing about what it points at — so the only comparable
pair would be the C struct against the Go struct, through a target threaded
alongside every leaf on both sides. That is a wider change than the check earns,
so it was left out deliberately rather than half-done.

### C string termination

Widths, classes and pointer targets all agree for a Go byte buffer behind a
`const char *`: the slot is an 8-byte pointer on both sides, and one indirection
down it is a C `char` against a Go `byte`, which is the one target pair the
sections above are happiest with. Nothing there looks at the property C actually
depends on, which is that the bytes end in a **NUL**.

Go never puts one there. A Go string carries its length, so every C string in
this tree is terminated by hand:

```go
p := &[]byte(path + "\x00")[0]        // pkg/llama/ggml.go
text += "\x00"                        // pkg/mtmd/mtmd.go
p := unsafe.StringData(text)
```

Delete the `+ "\x00"` from any of those and every rule above still passes, the
package still compiles, and C reads forward off the end of a Go allocation —
through whatever the allocator happened to put next — until it finds a zero byte.
That is the loudest possible failure with the least possible warning at build
time, and unlike the drift in #289 it does not even need upstream to change:
one edit is enough.

So each `char *` parameter slot is traced back to the buffer behind it and
requires **positive evidence** of a terminator. The trace runs over the enclosing
Go function, because that is where the buffer is built a statement or two before
the `Call` that hands it over: the avalue `unsafe.Pointer(&file)` says nothing by
itself, and every value `file` is assigned does. Three producer forms are
recognised, which are the ones the tree uses:

| producer | terminated by |
|---|---|
| `&[]byte(s)[0]`, `unsafe.SliceData([]byte(s))` | a `"\x00"` at the end of `s` |
| `unsafe.StringData(s)` | the same, including `s += "\x00"` earlier in the function |
| `&[]byte{...}[0]` | a composite literal whose last element is `0` |

Only the *last* operand of a concatenation can terminate it, so `path + "\x00"`
counts and `"\x00" + path` does not. `nulHelpers` in `goside.go` is there for a
helper function that returns a terminated buffer, and is a closed hand-checked
list for the same reason `constAliases` and `callbackTags` are — a rule loose
enough to guess which helper terminates its result is loose enough to accept one
that does not. It is empty today, because all six live sites append the
terminator inline.

This is a **RULE 2** finding rather than a class of its own, because the failure
is RULE 2's failure — C reading bytes the Go side never meant to hand over — and
because, unlike a signedness difference, nothing about it is permitted: it
belongs in the exit code. What differs from the rest of RULE 2 is the *evidence*,
dataflow rather than a width, which is why the scope is counted separately as
`C string buffers checked` on the `SUMMARY`.

That scope is deliberately narrow, and narrowed rather than filled with skips. A
`char *` fed a `*byte` that arrived as a function parameter, was written by C, or
came out of `make` for C to fill into is not traceable to a Go string, and is not
an unterminated buffer either: it names nothing to decide about, exactly as a
`void *` names no pointer target. Those slots are therefore outside the claim
rather than skips — a skip means "this should have been checked and was not" —
and a clean run still reports 0 skips. What the count guarantees is that a
narrowing of the traced forms shows up as a smaller number rather than as
silence. `char **`, `uint8_t *` and `int8_t *` are outside it too: the buffer is
one indirection further away than this traces in the first case, and a counted
byte buffer whose length travels in another parameter has no terminator to
require in the other two.

### Signedness

`kindCompat` merges `KindSint` and `KindUint`, and it is right to: a same-width
signed and unsigned integer occupy the same register and the same bytes, so
nothing about the ABI is broken. What breaks is the *meaning*. `llama_decode`
returns negative on error and `llama_token` uses `-1` as a sentinel, so reading
either through an unsigned Go type turns a failure into 4294967295 while every
width, offset and class still matches.

So it is reported as its own class, in its own `SIGNEDNESS` section, and it
**never affects the exit code** — calling it an ABI violation would be wrong, and
a CI failure for something the ABI permits is how a checker gets ignored. It is
raised only where the value is actually *interpreted* rather than forwarded: a
return buffer, and a pointer target from the section above.

Two exemptions keep it from being noise. `ffi.Arg` is skipped, and has to be:
RULE 3 *requires* it for an integer return, it is unsigned by construction, so it
is a container rather than a meaning — what a caller converts it to afterwards is
one expression further on than this tool looks, which is also why yzma reports
zero of these today. And a **1-byte** target is skipped, because plain `char` is
signed in the type table while its signedness is implementation-defined in C, and
`const char *` fed a Go `*byte` is the string idiom of every binding in the tree.

`jupiterrider/ffi` states both halves of RULE 3 in its `Fun.Call` doc comment
(`fun.go:19-20`): *"You cannot use integer types smaller than 8 bytes here
(float32 and structs are not affected)."*

## Output

Only mismatches are printed, plus five accounting sections that matter as much
as the findings, because they are what makes *"these are the only ones"* a
measurement rather than an assertion:

- `NOT VERIFIED (skips)` — every slot the tool could not decide, and why. A
  clean run reports **0 skips**; anything else means the coverage claim is
  weaker than the violation list suggests.
- `SIGNEDNESS` — values read through a Go type of the other sign. Same bytes,
  same register, other meaning, so these are **not** violations and do not affect
  the exit code; see *Signedness*, above. yzma reports 0 of them today.
- `STRUCT-BY-VALUE COMPARISONS` — each struct slot with the cif descriptor size
  and member count next to the C struct's, so those checks are visible rather
  than implied, plus a note for each Go struct carrying members past the end of
  the C struct — legitimate, but the reason a size comparison alone cannot be
  exact, so it stays on the page. `layout?` in place of a member count means the layout could not
  be resolved, and the matching `NOT VERIFIED` entry says which side.
- `C DECLS WITH NO BINDING` — the reverse of `BINDINGS WITH NO C DECL`, and
  deliberately **not** a check. yzma binds a chosen subset of llama.cpp, so an
  unbound declaration is not a defect: this section never produces a violation
  and never affects the exit code. 862 declarations are parsed and 265 bound, and
  without this the other 597 are simply unaccounted for. It prints one line per
  header, because that is the unit a maintainer can act on:

  ```
  ggml.h:              0 of 362 bound (362 unbound)
  llama.h:           196 of 239 bound (43 unbound)
  mtmd-helper.h:      10 of  24 bound (14 unbound)
  ```

  The count is always printed; `-v` adds every unbound name with its header line.
  The value is drift over time rather than any single run: a function yzma
  *should* bind appearing upstream, or the neighbours of a bound function
  changing, is invisible unless the unbound set is written down where a diff can
  see it. It makes the coverage claim measurable; it does not claim the coverage
  is right.

- `PARTIALLY MIRRORED C ENUMS` — the same idea one layer down, and equally not a
  check. RULE 4 walks the **Go** side, so a C enum member yzma never transcribed
  is invisible to it: nothing compares `LLAMA_SPLIT_MODE_TENSOR` against
  anything, because there is no `SplitModeTensor`. Mirroring a subset is
  deliberate, exactly as binding a subset is, so this never produces a violation
  and never affects the exit code. What it is, is the signal for the one event
  RULE 4 exists to catch: an enum gaining a member upstream is *both* a new
  unmirrored name here *and* a value shift in every mirrored member declared
  after it, and the first of those is visible before anybody's numbering is
  wrong. `GGML_TYPE_COUNT` is the same event caught after the fact — the three
  types that moved it from 40 to 43 appear in this section as the members yzma
  has not caught up with, and as a RULE 4 violation on `GGMLTypeCOUNT`.

  Only enums with **at least one** mirrored member are listed. An enum with none
  is not a partial mirror, it is simply unused, and the several hundred members
  of the ggml enums yzma never touches would drown the report:

  ```
  ggml_backend_dev_type:               4 of   5 members mirrored (1 not mirrored)
  ggml_type:                          33 of  36 members mirrored (3 not mirrored)
  llama_load_mode:                     5 of   6 members mirrored (1 not mirrored)
  llama_split_mode:                    3 of   4 members mirrored (1 not mirrored)
  mtmd_input_chunk_type:               3 of   4 members mirrored (1 not mirrored)
  ```

  As with the unbound declarations the counts are always printed and `-v` adds
  every member with its value and header line, and the totals are on the
  `SUMMARY`.

- `SUMMARY` — declarations parsed, bindings matched, and checked/clean counts
  per rule, plus struct layouts, pointer targets and C string buffers compared.
  Those three counts are separate because a
  layout comparison that quietly stopped resolving would still leave the rule it
  belongs to reporting zero violations. The RULE 4 line carries three more
  numbers for the same reason:

  ```
  Rule4 constants checked:   169 / 168 clean (C constants parsed: 421, unevaluable: 8; yzma-local: 27)
  ```

  *unevaluable* is the rule's coverage limit — C constants whose value the
  expression evaluator could not pin down, which are mostly non-integer macros
  like `GGML_API` plus `GGML_MEM_ALIGN`, defined differently in three `#if`
  branches this tool does not evaluate. None of them can be compared against, and
  a Go constant mirroring one becomes a skip rather than a pass. *yzma-local* is
  how many constants `goOnlyConsts` excluded; if that number climbs without the
  list being edited, something is being excluded that should not be.

  `-v` adds a `CONST` line per comparison, naming the C constant each Go constant
  was matched to and the header line it came from, so the mapping is auditable
  rather than trusted.

  The RULE 5 line carries the same kind of coverage limit:

  ```
  Rule5 callbacks checked:   4 / 4 clean (C function-pointer typedefs parsed: 27, unparseable: 0)
  ```

  and one for the members:

  ```
  fn-ptr members compared:   5 / 5 clean (code pointers traced: 2)
  ```

  *unparseable* counts function-pointer typedefs whose signature this parser
  could not take apart. Most of the 27 belong to ggml internals yzma never
  implements, so the number being non-zero is not by itself a hole — but a
  callback site linking to one of them is a skip, so the hole can never be
  silent. `-v` adds a `CALLBACK` line per site, naming the typedef it was linked
  to and how the link was made.

Argument types come from `go/types`, not regex, so a named type like
`llama.SeqId` or an alias like `ffiTypeSize` resolves to a real width instead of
being guessed at.

## Correctness gate

A tool that reports nothing is indistinguishable from a tool that checks
nothing, so the checker has to demonstrate it can still find things:

```sh
go test ./...
```

`testdata/fixture` is a miniature binding package carrying one planted defect
per rule per direction — including both struct-layout directions, each planted
in a 24-byte struct that stays 24 bytes on every side, so a size-only comparison
passes it, and a transposition plant that is byte-for-byte identical to its C
struct with two members exchanged — plus four deliberately clean bindings: one
scalar, one struct-by-value, one in the `llama.Batch` shape whose Go struct
appends bookkeeping past the end of the C struct (pinning that the tolerance for
that tail does not also swallow a short Go struct), and one whose members are
spelled differently on each side (pinning that the alias table does not report
those as transposed).

For RULE 4 the fixture header grows an enum whose second member is the one
"upstream inserted", so the Go constant for the member after it still holds the
pre-insertion value; the other three members of that enum, a second enum
covering each initialiser form the C evaluator has to handle (a negative
sentinel, a shift, a hex literal) and three `#define`s are the controls, matched
by comment and by normalisation respectively.

For the `nfixed` check the header grows two variadic declarations of the same
shape — two parameters before the `...` — bound as one plant with `nfixed` 1 and
one control with `nfixed` 2. Since yzma has no variadic binding of its own, the
fixture is the only place this check is exercised at all, which makes the control
half of it the part that matters: a rule that reported a correct `PrepVar` would
be worse than no rule.

The coverage inventory has one fixture declaration nothing binds, and the enum
inventory a fixture enum with one member mirrored and two not — next to two enums
mirrored in full, which must therefore *not* appear. The assertion that the
partial one appears is the small half; the assertion that its unmirrored members
are **not** violations, and do not move the violation count, is the one that pins
the framing.

For the pointer-target comparison the header grows four pointer-out
declarations of the same shape, two per half. `fx_get_scores` declares a `float *`
and is bound to a Go `*int32` — the plant, whose slot is a pointer of the same
width on both sides, so every other rule passes it — against `fx_get_logits`,
which is the same declaration bound to the `*float32` it names. `fx_get_count`
declares an `int32_t *` read through a `*uint32` and `fx_decode` returns a
negative error code into an unsigned 8-byte buffer: those two are the signedness
plants, and the gate asserts both that they *are* reported in the `SIGNEDNESS`
list and that they are **not** violations and do not move the violation count,
which is the assertion that pins the framing. `fx_get_token` and `fx_decode_ok`
are their controls, same declarations read through the signed types, and the gate
also pins that no 1-byte `char` target is ever reported — without that exemption
the check would flag every string parameter in yzma and nothing else.

For the NUL-termination check the header grows three `const char *` declarations
of the same shape. `fx_set_name` is fed `&[]byte(name)[0]` with nothing appended —
the plant, whose slot is a pointer to a `char` read through a `*byte` on both
sides, so every width, class and pointer-target comparison passes it — against
`fx_set_path` and `fx_set_text`, one control per idiom the tree actually uses:
`&[]byte(path + "\x00")[0]`, and `text += "\x00"` before `unsafe.StringData`.
The control half is the one that matters, because those two are written exactly
as all six real string sites are: a rule that reported either would report every
string yzma passes. The gate also pins the scope count at three, so the two
`char *` slots that are *out* of scope — `fx_desc`'s output buffer from `make`
and `fx_mode_from_str`'s `*byte` parameter — stay out of it and stay out of the
skips.

For RULE 5 the fixture header grows four function-pointer typedefs and the
fixture package four callback sites, one plant and one clean control per form: a
descriptor declaring `&ffi.TypeSint32` where the typedef says `float`, a purego
closure one parameter short of its typedef, and — the half that matters as much
— a descriptor that models its typedef exactly and a closure whose parameters and
word-sized result match, both written the way all four real yzma sites are, so a
rule that reported them would report every one of those too. The gate asserts
that every plant is found *and* that the clean binding is never reported, along
with the accounting counters and run-to-run repeatability.

For the function-pointer members the fixture header grows `struct fx_hooks`, two
of whose four members are function pointers, and three by-value declarations that
take it. Three Go structs mirror it, all three byte-for-byte identical to it and
to a single correct cif descriptor, so every layout comparison passes all three
and the plants are only in what the members *hold*: `HooksFunc` declares
`cb_progress` as a Go `func` value, and `Hooks.SetHooks` writes the closure built
for `llama_fx_report_callback` into the `cb_abort` C calls as
`llama_fx_abort_callback`. `HooksClean` is the control, and the half that matters
— one member per callback form, each holding the code pointer of the callback
that actually implements it, written exactly as yzma's two live
`SetProgressCallback` methods write theirs. It cannot be pinned by name, because
all three mirror the same C struct and a member finding is named for the C member,
so the gate pins it by asserting no violation mentions that Go struct at all. The
counters pin the scope at six members compared with four code pointers traced, so
the two members with nothing stored in them stay out of both the traced count and
the skips.

**Treat a failing `go test` as "do not trust anything this run reports."**

The previous gate was "a run must re-derive these two live defects in yzma."
That is no longer usable, because those defects are fixed — which is exactly the
failure mode the fixture avoids.

## Files

| file | role |
|---|---|
| `main.go` | the five rules, the report, and the CLI |
| `sources.go` | resolving the yzma tree, the llama.cpp ref, and the headers |
| `cheader.go` | C declaration parser: comment blanking that preserves byte offsets, `DEPRECATED(...)` unwrapping, typedef/enum resolution, top-level comma splitting |
| `cstruct.go` | C struct layout for struct-by-value slots; iterates to a fixpoint because `llama.h` structs reference typedefs declared in `ggml-backend.h` |
| `cenum.go` | RULE 4: C enum members and integer `#define`s with a small C constant-expression evaluator, the Go constants from `go/types`, the name matching between them, and the partially-mirrored-enum inventory |
| `ccallback.go` | RULE 5: C function-pointer typedefs, the link from a Go callback site to the typedef it implements, the comparison for both callback forms, and the function-pointer struct members C reaches a callback through |
| `layout.go` | flattens a C struct, a cif descriptor and a Go struct to a common member list, diffs them, and matches members by name to find transpositions |
| `goside.go` | `go/packages` type-checked walk: `lib.Prep`/`PrepVar` → binding spec, `<var>.Call(...)` → the Go type libffi will actually read bytes from and, for a C string, the buffer it was built from, `ffi.PrepCif`/`purego.NewCallback` → callback site, `ffi.PrepClosureLoc` and the assignments that install a code pointer in a struct field |
| `main_test.go` | the correctness gate |

## Assumptions

Four things this tool takes as given rather than verifies. None of them is
currently false, but each one is load-bearing, and a reader deciding how much a
green run is worth cannot see them from the output.

**Sizes are computed for arm64 only.** `goside.go` lays out every Go type with
`types.SizesFor("gc", "arm64")`, so every width in this report — 8-byte pointer,
8-byte `size_t`, 8-byte `ffi.Arg`, and every struct offset derived from them — is
that target's. For the types that actually cross this boundary amd64 agrees
member for member, so the audit is valid there too, but nothing in the tool
checks that: it prints arm64 numbers whatever `GOARCH` it runs on. A 32-bit
target is where it would break rather than merely mislead, because a 4-byte
pointer changes every interior offset in `llama_context_params` and makes RULE
3's "integer returns need 8 bytes" the wrong number. yzma has no 32-bit target to
break on: `pkg/download/arch.go` declares exactly two architectures, `amd64` and
`arm64`, and `MustParseArch` panics on anything else, so an unsupported `GOARCH`
cannot get as far as loading a library. `jupiterrider/ffi` narrows it further —
its build tag is `((freebsd || linux || windows || darwin) && (amd64 || arm64))
|| (linux && riscv64)`, all 64-bit — and CI builds on `ubuntu-latest`,
`macos-latest` and `windows-2025`. So this is an assumption about a platform yzma
does not support, and the thing to watch is a 32-bit architecture being added to
that table, not the audit being wrong today.

**Every C enum is assumed to be 4 bytes.** `classify()` in `cheader.go` maps
`enum X` to `KindSint, 4`, which is what every ABI here picks for an enum whose
values fit in `int` — and it is why `llama_pooling_type` in a parameter slot
compares equal to `&ffi.TypeSint32`. A C enum's width is implementation-defined,
though: one whose values exceed `int32`, or one declared with a fixed underlying
type (`enum e : uint64_t`), would be wider, and the tool would then call an
8-byte C slot 4 bytes and report a correct binding as broken — or pass a broken
one. Neither shape exists in these headers. No enum in b10219 declares an
underlying type, and the values RULE 4 parses run from `-1` to `512`
(`LLAMA_TOKEN_ATTR_SINGLE_WORD`, `GGML_SCALE_FLAG_ANTIALIAS`), so every one of
them fits `int` with five orders of magnitude to spare. The risk is real in C and
theoretical here, and RULE 4 is where it would first become visible, since it is
already reading the values.

**Struct returns are exempt from RULE 3's 8-byte rule on the strength of a doc
comment.** RULE 3 requires an integer return buffer to be 8 bytes because libffi
always stores a full `ffi_arg`, and deliberately does *not* require it of a
struct return: it only requires the buffer to be at least as large as the
descriptor. That exemption rests entirely on `jupiterrider/ffi`'s `Fun.Call` doc
comment (`fun.go`), which says of the return pointer:

> You cannot use integer types smaller than 8 bytes here (float32 and structs are
> not affected). Use [Arg] instead and typecast afterwards.

yzma has a **1-byte** struct return — `ffiSamplerChainParams`, a single `bool`,
in `pkg/llama/sampling.go` — so this is not a hypothetical exemption. If libffi
in fact wrote 8 bytes there, it would clobber 7 bytes of adjacent Go memory on
every call, and this tool would report the site as clean. It is a claim the tool
depends on and does not test: nothing here executes a call or examines what
libffi writes, and the same sentence is the sole authority for the float32 half
of the rule.

**Struct-by-value calling conventions are libffi's problem, not this tool's.**
What is compared is the *layout* of the three representations — offset, width and
ABI class per member. How arm64 then passes such a struct (in registers up to 16
bytes, indirectly through memory beyond that, `HFA` rules for all-float structs)
is not modelled at all, because libffi derives it from the descriptor. A correct
descriptor is therefore the whole of what these rules can ask for; a libffi bug
in that derivation would be invisible here.

## What it does not cover

- **What a callback does with its arguments.** RULE 5 covers the four callback
  sites and the five function-pointer struct members now, but only their
  signatures: that the closure is handed the values C
  passed, in the widths C passed them. It cannot say the body then reads `arg[0]`
  as the type the descriptor declares — `*(*float32)(arg[0])` behind an
  `&ffi.TypeFloat` slot is right and `*(*int32)(arg[0])` is not, and both
  type-check. Nor does it cover the closure's lifetime: a `PrepClosureLoc` code
  pointer handed to C outlives the Go call that made it, and keeping it alive is
  the pointer-lifetime problem below.
- **Members yzma renames past the alias table.** A transposition is found by
  matching member names, so a Go field whose name resolves to no C member is
  reported as `NOT VERIFIED` for that check rather than compared. Today the
  alias table covers every such name, so a clean run has 0 skips; a future
  rename shows up as a skip until it is added.
- **Constants whose C value the evaluator cannot reach.** RULE 4 does not run
  the preprocessor, so a name defined differently in several `#if` branches
  (`GGML_MEM_ALIGN`) has no single value to compare against, and neither does a
  macro that is not an integer at all. Those are counted on the RULE 4 summary
  line, and a Go constant mirroring one is a skip naming the reason rather than a
  silent pass. None of the ~30 constants yzma actually mirrors is in that set
  today.
- **Constants yzma renames past the alias table, and its own.** A Go constant
  that matches no C name is reported as a skip until it is either given a
  `// C_NAME` comment or listed in `goOnlyConsts`. That means the rule is only
  as complete as those lists — but never quietly so, because the incomplete case
  is a skip and a clean run has none.
- **What a constant is used *for*.** RULE 4 says `PoolingTypeMean` is the same
  integer as `LLAMA_POOLING_TYPE_MEAN`; it cannot say that yzma passes it to a
  parameter that takes a `llama_pooling_type` rather than to some other enum of
  the same width. Nothing in the C ABI distinguishes those.
- **Whether the unbound declarations *should* be bound.** The coverage inventory
  counts them and names them; it cannot say which of them yzma is missing. That
  judgement is the maintainer's, and the inventory exists so it can be made
  against a number that moves rather than against a guess.
- **Whether an unmirrored enum member *should* be mirrored.** Same limit as the
  unbound declarations, and the same answer: the enum inventory says which
  members have no Go constant, not which of them yzma is missing. `GGML_TYPE_Q2_0`
  and `LLAMA_SPLIT_MODE_TENSOR` are both in that list, and only a maintainer can
  say the first matters and the second does not.
- **What a *data* pointer inside a by-value struct points at.** The comparison in
  *Pointer targets* runs on parameter and return slots only. A `float *` member
  of a struct is one 8-byte leaf on all three sides, and the cif descriptor
  carries no target information for it at all, so the drift is invisible there.
  Deliberately left out rather than half-done. A **function**-pointer member is
  the exception, and is covered: what it points at is a signature the header
  names, so there is something to compare.
- **Which callback a struct member holds, where nothing stores one.** The member
  check follows a code pointer from the field it is assigned to back to the
  callback site that built it, inside the package that does the assigning. A
  member yzma never writes — `cb_eval` and `abort_callback` today — carries no
  code pointer to identify, and neither does one whose value arrives from a
  caller or through a helper this does not trace. Those are outside the claim
  rather than skips, and `code pointers traced` on the `SUMMARY` is what makes
  the difference between the two scopes a number. Nor does it cover the member's
  *lifetime*: `progressCallbackCode` is a package-level variable that a second
  `SetProgressCallback` overwrites, which is the pointer-lifetime problem below.
- **What a caller converts a return buffer to.** RULE 3 requires an integer
  return to be read into an `ffi.Arg`, which is unsigned by construction, so the
  signedness check cannot say anything about the `int32(result)` on the next line.
  That is why yzma has 0 signedness findings today while the fixture demonstrates
  both halves of the check.
- **A C string whose buffer is not built in the same function.** The termination
  check traces the enclosing Go function, so a `char *` that arrives as a
  parameter, is written by C, or is handed over as a struct *member* — which is
  how `mtmd.InputText.Text` reaches C, one level of indirection past every
  parameter slot the rules look at — is outside its scope rather than a skip.
  `NewInputText` does append its `"\x00"`, and nothing here proves it. The
  `C string buffers checked` count is what makes that scope a number.
- **Pointer lifetime.** Go pointers stored in C-visible memory, missing
  `runtime.KeepAlive`, and slices retained over C-owned memory are a separate
  class of hazard entirely.

## Note on the nested module

This directory has its own `go.mod` so the `golang.org/x/tools` dependency never
reaches the root yzma module. `go build ./...`, `go vet ./...` and
`go test ./...` at the repo root skip it entirely.
