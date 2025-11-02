# Installer

This program downloads the latest `llama.cpp` pre-built binaries for the current operating system.

## Running

Set the `YZMA_LIB` environment variable to the target installation directory before running.

```shell
export YZMA_LIB="/home/ron/Development/yzma/lib"
```

```shell
$ go run ./examples/installer/ -processor cuda -upgrade
installing llama.cpp version b6924 to /home/ron/Development/yzma/lib
done.
```

### Flags

```shell
  -help string
        show help
  -lib string
        path to llama.cpp compiled library files (leave empty to use YZMA_LIB env var)
  -processor string
        processor to use (cpu, cuda, metal, vulkan) (default "cpu")
  -upgrade
        upgrade existing installation
  -version string
        version of llama.cpp to install (leave empty for latest)
```

