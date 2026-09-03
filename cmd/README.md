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
   --version value, -v value    version of llama.cpp to install (leave empty for the version this yzma release uses)
   --lib value, -l value        path to llama.cpp compiled library files [$YZMA_LIB]
   --processor value, -p value  processor to use (cpu, cuda, metal, vulkan) (default: "cpu")
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
```

## Other commands

See the `yzma help` command for more information about the other things you can do with the `yzma` CLI tool.

## `yzma-checker`

The `yzma-checker` directory holds a separate developer tool, not a subcommand of the `yzma` CLI. It compares the FFI parameter and return types of each yzma binding, and the values of the constants yzma mirrors, with the llama.cpp headers. It is a nested Go module, so `go build ./...` and `go test ./...` at the repo root do not include it.

```shell
make check-ffi
```

See [yzma-checker/README.md](./yzma-checker/README.md) for what it verifies and how.

## Using the `yzma` command to check an installation

`yzma install` checks the digest of each archive as it downloads it, but the archive is
removed as soon as it is extracted. The `yzma verify` command checks the files that are
in place, against the digests that the publisher recorded for the release.

```
NAME:
   yzma verify - Check the installed llama.cpp libraries against their published digests

USAGE:
   yzma verify [command options]

OPTIONS:
   --lib value, -l value      path to llama.cpp compiled library files [$YZMA_LIB]
   --version value, -v value  the llama.cpp version that must be installed (leave empty to take the installed one)
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

A file that was changed or removed makes the command exit with a status of 1:

```
$ yzma verify --lib /path/to/lib
llama.cpp b10783 in /path/to/lib
  changed    libllama.so.0.3.0
  missing    libggml-base.so.0.22.0
66 verified, 1 changed, 1 missing, 0 not part of this install
```

Notes:

- `yzma install` writes `yzma-install.json` beside the libraries to say what it put
  there. `yzma verify` needs it, so an installation made by an older yzma has to be
  installed again first.
- The record is beside the libraries, so anything that can change the libraries can
  change the record. Give `--version` to say which release must be there. The check then
  resolves the assets of that release itself and does not read the tag from the record.
- A directory can hold more than one install, so a file that is not part of this one is
  reported but does not fail the check. Add `--strict` to fail on those too.
- Only the assets that `llama-cpp-builder` builds carry digests for the files in them.
  An install that came from the `llama.cpp` release page has an archive digest but no
  file digests, so `yzma verify` says so instead of passing.
