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
   --version value, -v value    version of llama.cpp to install (leave empty for latest)
   --lib value, -l value        path to llama.cpp compiled library files [$YZMA_LIB]
   --processor value, -p value  processor to use (cpu, cuda, metal, vulkan) (default: "cpu")
   --upgrade, -u                upgrade existing installation (default: false)
   --quiet, -q                  suppress output during installation (default: false)
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
