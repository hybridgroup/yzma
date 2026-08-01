# yzma-checker

Mechanical verification of the **FFI parameter and return types** in every yzma
binding, against the llama.cpp headers for the build yzma installs.

yzma binds llama.cpp through purego/libffi, not cgo. There is no `cgocheck`, no
automatic pinning, and no compiler validation of argument widths, so a mismatch
between

1. the C declaration,
2. the libffi type descriptor passed to `ffi.PrepCif`, and
3. the Go variable whose address is passed to `ffi.Fun.Call`

is completely silent at build time and only shows up at runtime as a corrupted
value or corrupted memory. 265 bindings is far too many to eyeball credibly,
which is why this exists.

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
-v              dump every binding: cif types, C signature, and every call site
```

`-llama` and `-hdrs` are the offline paths.

## The three rules

**RULE 1 — the cif must match the C prototype.** Parameter for parameter, in
width and in class. On arm64: `size_t`/`uint64_t`/any pointer = 8 bytes,
`int`/`int32_t`/`enum`/`float` = 4, `double` = 8. A struct passed by value has
its own descriptor and is a common error site.

**RULE 2 — the Go variable must be exactly as wide as the cif claims.** libffi
reads exactly `ffi.Type.Size` bytes from each pointer it is handed. `&x` where
`x` is a Go `int32` behind an `ffi.TypeUint64` slot makes libffi splice 4 bytes
of adjacent Go memory into the high half of the C argument.

**RULE 3 — return buffers must match the return *kind*.**
- Integer/pointer returns: the buffer must be **8 bytes** (`ffi.Arg`), because
  libffi always stores a full `ffi_arg`. A narrower variable is written past.
- Float returns are the **opposite**: they must be `float32`. libffi stores only
  4 bytes and leaves the rest of the buffer alone, so an `ffi.Arg` ends up
  holding the IEEE-754 bit pattern in its low word — and `float32(someArg)` is a
  numeric conversion, not a bit reinterpretation.
- A `TypeVoid` return descriptor makes libffi write **nothing at all**, so a
  supplied buffer silently keeps its zero value.

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
  next to the C struct size, so those checks are visible rather than implied.
- `SUMMARY` — declarations parsed, bindings matched, and checked/clean counts
  per rule.

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
per rule per direction, plus one deliberately clean binding. The gate asserts
that every plant is found *and* that the clean binding is never reported, along
with the accounting counters and run-to-run repeatability.

**Treat a failing `go test` as "do not trust anything this run reports."**

The previous gate was "a run must re-derive these two live defects in yzma."
That is no longer usable, because those defects are fixed — which is exactly the
failure mode the fixture avoids.

## Files

| file | role |
|---|---|
| `main.go` | the three rules, the report, and the CLI |
| `sources.go` | resolving the yzma tree, the llama.cpp ref, and the headers |
| `cheader.go` | C declaration parser: comment blanking that preserves byte offsets, `DEPRECATED(...)` unwrapping, typedef/enum resolution, top-level comma splitting |
| `cstruct.go` | C struct layout for struct-by-value slots; iterates to a fixpoint because `llama.h` structs reference typedefs declared in `ggml-backend.h` |
| `goside.go` | `go/packages` type-checked walk: `lib.Prep`/`PrepVar` → binding spec, `<var>.Call(...)` → the Go type libffi will actually read bytes from |
| `main_test.go` | the correctness gate |

## What it does not cover

- **Callback descriptors.** Four bindings use `ffi.PrepCif` / `PrepClosureLoc` /
  `purego.NewCallback` rather than `lib.Prep`, so they fall outside the 265.
- **Struct field order.** Only total struct sizes are compared, so two structs
  that agree in size but disagree in field order pass.
- **Pointer lifetime.** Go pointers stored in C-visible memory, missing
  `runtime.KeepAlive`, and slices retained over C-owned memory are a separate
  class of hazard entirely.

## Note on the nested module

This directory has its own `go.mod` so the `golang.org/x/tools` dependency never
reaches the root yzma module. `go build ./...`, `go vet ./...` and
`go test ./...` at the repo root skip it entirely.
