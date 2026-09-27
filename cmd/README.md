# yzma

Command line tool for managing yzma and llama.cpp libraries.

## Installation

```shell
go install github.com/hybridgroup/yzma@latest
```

## Commands

```
NAME:
   yzma - YZMA command line tool

USAGE:
   yzma [global options] command [command options]

COMMANDS:
   install  Install llama.cpp libraries used by yzma
   verify   Check the installed llama.cpp libraries against their published digests
   system   Show llama.cpp system information
   llama    Show most recent llama.cpp version
   model    Manage models
   version  Show yzma version
   info     Show yzma version
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

## Using the `yzma` command to install `llama.cpp`

You can use the `yzma install` command to download the llama.cpp pre-built binaries for the current operating system.

```
NAME:
   yzma install - Install llama.cpp libraries used by yzma

USAGE:
   yzma install [command options]

OPTIONS:
   --version value, -v value    version of llama.cpp to install, optionally as VERSION@sha256:DIGEST to pin the digests (leave empty for the version this yzma release uses)
   --lib value, -l value        path to llama.cpp compiled library files [$YZMA_LIB]
   --processor value, -p value  processor to use (cpu, cuda, metal, openvino, rocm, vulkan) (default: "cpu")
   --upgrade, -u                upgrade existing installation (default: false)
   --quiet, -q                  suppress output during installation (default: false)
   --verify value               how to check the digest of each download (available, require, off) (default: "available") [$YZMA_VERIFY]
   --help, -h                   show help
```

Here are a few examples:

```
# Install with default settings (uses YZMA_LIB env var)
yzma install

# Install to specific path
yzma install --lib /path/to/lib

# Install specific version with CUDA
yzma install --lib /path/to/lib --version b1234 --processor cuda

# Upgrade existing installation
yzma install --lib /path/to/lib --upgrade

# Using short flags
yzma install -l /path/to/lib -v b1234 -p cuda -u

# Install a version and pin the digests it must have
yzma install --lib /path/to/lib --version b1234@sha256:<digest>
```

A version with no digest works exactly as before, so `--version b1234` and a
bare `yzma install` need nothing new.

The digest to pin is the SHA-256 of the release's digest manifest, which
`llama-cpp-builder` publishes with each release. The release notes for the tag print the
full pin, and the version file has it for the newest build:

```console
$ curl -s https://hybridgroup.github.io/llama-cpp-builder/version.json
{"tag_name":"b10816","manifest_sha256":"<digest>","pin":"b10816@sha256:<digest>"}

$ yzma install --lib /path/to/lib --version b10816@sha256:<digest>
```

This is not the digest of a platform archive. The manifest holds those, one per
asset. See [Verify an installation](https://yzma.ai/docs/guides/verifying/)
for the whole chain.

## Other commands

See the `yzma help` command for more information about the other things you can do with the `yzma` CLI tool.

## `yzma-checker`

The `yzma-checker` directory holds a separate developer tool, not a subcommand of the `yzma` CLI. It checks the FFI parameter and return types of each yzma binding, and the values of the constants yzma mirrors, against the llama.cpp headers. It is a nested Go module, so `go build ./...` and `go test ./...` at the repo root do not include it.

```shell
make check-ffi
```

See [yzma-checker/README.md](./yzma-checker/README.md) for what it verifies and how.

## Using the `yzma` command to check an installation

`yzma install` checks the digest of each archive as it downloads it, but deletes the
archive once it is extracted. The `yzma verify` command checks the installed files
against the digests the publisher recorded for the release.

```
NAME:
   yzma verify - Check the installed llama.cpp libraries against their published digests

USAGE:
   yzma verify [command options]

OPTIONS:
   --lib value, -l value      path to llama.cpp compiled library files [$YZMA_LIB]
   --version value, -v value  the llama.cpp version that must be installed, optionally as VERSION@sha256:DIGEST to pin the digests (leave empty to take the installed one)
   --strict                   also fail when a file in the directory is not part of the install (default: false)
   --json                     write the report as JSON (default: false)
   --help, -h                 show help
```

```
$ yzma verify --lib /path/to/lib
llama.cpp b10783 in /path/to/lib
68 verified, 0 changed, 0 missing, 0 not part of this install
ok.
```

If a file was changed or removed, the command exits with status 1:

```
$ yzma verify --lib /path/to/lib
llama.cpp b10783 in /path/to/lib
  changed    libllama.so.0.3.0
  missing    libggml-base.so.0.22.0
66 verified, 1 changed, 1 missing, 0 not part of this install
```

Notes:

- `yzma install` writes `yzma-install.json` next to the libraries to record what it
  installed. `yzma verify` needs this file, so an installation made by an older yzma
  must be installed again first.
- `yzma install` also writes `yzma-manifest.json`, which holds the release digests, so
  `yzma verify` works offline. If an installation has no manifest, the command fetches
  one and saves it, so only the first check needs the network.
- The record sits next to the libraries, so anything that can change the libraries can
  change the record too. Pass `--version` to say which release must be installed. The
  check then resolves that release's assets itself and ignores the tag in the record.
- A directory can hold more than one install, so files that are not part of this one are
  reported but do not fail the check. Add `--strict` to fail on them too.
- Only assets built by `llama-cpp-builder` have digests for the files inside them.
  An install from the `llama.cpp` release page has an archive digest but no file
  digests, so `yzma verify` reports that instead of passing.
- `--version` also accepts `VERSION@sha256:<digest>`, where the digest is the SHA-256 of
  the release's digest manifest. The expected value then comes from wherever you store
  it, not from the site that serves the manifest. A pin makes the check mandatory,
  so it cannot be combined with `--verify off`. See [Verify an installation](https://yzma.ai/docs/guides/verifying/)
  for what a pin does and does not show.
