# hello

The "hello, world" of `yzma`.

## Install yzma

Make sure you have installed `yzma` as described in [INSTALL.md](../../INSTALL.md)

Once installed, follow the instructions from your installer to set the `YZMA_LIB` environment variable.

## Download model

Download the model to the default location on your machine:

```shell
yzma model get -u https://huggingface.co/bartowski/SmolLM2-135M-Instruct-GGUF/resolve/main/SmolLM2-135M-Instruct-Q4_K_M.gguf
```

## Run

```shell
$ go run ./examples/hello/


"Yes, I'm ready to go."
```
