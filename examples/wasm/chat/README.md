# wasm/chat

Uses `yzma` to generate text in a browser. It has the same structure as [examples/hello](../../hello), but it uses the [llamawasm](../../../pkg/llamawasm) package and sends each piece of text to the page.

The page is [wasm/index.html](../../../wasm/index.html). See [wasm/README.md](../../../wasm/README.md) for how it works.

## Building

Run these commands from the root of the repo.

```shell
make download-llama.cpp-wasm
make wasm-example
```

This builds `build/wasm/yzma.wasm` with TinyGo. `make wasm-example-go` builds it with the standard Go toolchain.

## Running

```shell
make serve-wasm
```

Open <http://localhost:8080>. The page downloads [SmolLM-135M Q2_K](https://huggingface.co/QuantFactory/SmolLM-135M-GGUF) by default. Enter the URL of another GGUF model to use a different one.

Add `?mode=cpu` or `?mode=webgpu` to the URL to pick the backend.

## Running without a browser

```shell
yzma model get -u https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf
make test-wasm
```

## Page calls

| Function | What it does |
| --- | --- |
| `yzmaLoadModel(url)` | Downloads a model and creates a context. |
| `yzmaOpenModel(path)` | Loads a model that is already in the module's filesystem. |
| `yzmaGenerate(prompt, maxTokens)` | Generates text and sends each piece to the page. |
