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
calls a Go closure back. 265 bindings, 169 mirrored constants and 4 callbacks
are far too many to eyeball credibly, which is why this exists.

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
                every callback with the C typedef it was linked to
```

`-llama` and `-hdrs` are the offline paths.

## The five rules

**RULE 1 — the cif must match the C prototype.** Parameter for parameter, in
width and in class. On arm64: `size_t`/`uint64_t`/any pointer = 8 bytes,
`int`/`int32_t`/`enum`/`float` = 4, `double` = 8. A struct passed by value has
its own descriptor and is a common error site, so its members are compared one
by one — see *Struct layouts*, below.

**RULE 2 — the Go variable must be exactly as wide as the cif claims.** libffi
reads exactly `ffi.Type.Size` bytes from each pointer it is handed. `&x` where
`x` is a Go `int32` behind an `ffi.TypeUint64` slot makes libffi splice 4 bytes
of adjacent Go memory into the high half of the C argument.

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

There are four such sites, in two forms:

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

`jupiterrider/ffi` states both halves of RULE 3 in its `Fun.Call` doc comment
(`fun.go:19-20`): *"You cannot use integer types smaller than 8 bytes here
(float32 and structs are not affected)."*

## Output

Only mismatches are printed, plus three accounting sections that matter as much
as the findings, because they are what makes *"these are the only ones"* a
measurement rather than an assertion:

- `NOT VERIFIED (skips)` — every slot the tool could not decide, and why. A
  clean run reports **0 skips**; anything else means the coverage claim is
  weaker than the violation list suggests.
- `STRUCT-BY-VALUE COMPARISONS` — each struct slot with the cif descriptor size
  and member count next to the C struct's, so those checks are visible rather
  than implied, plus a note for each Go struct carrying members past the end of
  the C struct — legitimate, but the reason a size comparison alone cannot be
  exact, so it stays on the page. `layout?` in place of a member count means the layout could not
  be resolved, and the matching `NOT VERIFIED` entry says which side.
- `SUMMARY` — declarations parsed, bindings matched, and checked/clean counts
  per rule, plus struct layouts compared. The layout count is separate because a
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

For RULE 5 the fixture header grows four function-pointer typedefs and the
fixture package four callback sites, one plant and one clean control per form: a
descriptor declaring `&ffi.TypeSint32` where the typedef says `float`, a purego
closure one parameter short of its typedef, and — the half that matters as much
— a descriptor that models its typedef exactly and a closure whose parameters and
word-sized result match, both written the way all four real yzma sites are, so a
rule that reported them would report every one of those too. The gate asserts
that every plant is found *and* that the clean binding is never reported, along
with the accounting counters and run-to-run repeatability.

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
| `cenum.go` | RULE 4: C enum members and integer `#define`s with a small C constant-expression evaluator, the Go constants from `go/types`, and the name matching between them |
| `ccallback.go` | RULE 5: C function-pointer typedefs, the link from a Go callback site to the typedef it implements, and the comparison for both callback forms |
| `layout.go` | flattens a C struct, a cif descriptor and a Go struct to a common member list, diffs them, and matches members by name to find transpositions |
| `goside.go` | `go/packages` type-checked walk: `lib.Prep`/`PrepVar` → binding spec, `<var>.Call(...)` → the Go type libffi will actually read bytes from, `ffi.PrepCif`/`purego.NewCallback` → callback site |
| `main_test.go` | the correctness gate |

## What it does not cover

- **What a callback does with its arguments.** RULE 5 covers the four callback
  sites now, but only their signatures: that the closure is handed the values C
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
- **Pointer lifetime.** Go pointers stored in C-visible memory, missing
  `runtime.KeepAlive`, and slices retained over C-owned memory are a separate
  class of hazard entirely.

## Note on the nested module

This directory has its own `go.mod` so the `golang.org/x/tools` dependency never
reaches the root yzma module. `go build ./...`, `go vet ./...` and
`go test ./...` at the repo root skip it entirely.
