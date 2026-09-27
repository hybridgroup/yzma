# wasm/vlm

Uses `yzma` to answer a question about an image in a browser. It has the same structure as [wasm/chat](../chat) and adds the multimodal calls of the [llamawasm](../../../pkg/llamawasm) package.

The page decodes the image on a canvas and sends the pixels to the program, so any format the browser can read works.

The page is [wasm/vlm.html](../../../wasm/vlm.html). See the Images section of [wasm/README.md](../../../wasm/README.md) for how it works.

## Building

Run these commands from the root of the repo.

```shell
make download-llama.cpp-wasm
make wasm-vlm-example
```

This builds `build/wasm/yzma-vlm.wasm` with TinyGo.

## Running

```shell
make serve-wasm
```

Open <http://localhost:8080/vlm.html>. The page downloads [SmolVLM-256M-Instruct Q8_0](https://huggingface.co/ggml-org/SmolVLM-256M-Instruct-GGUF) and its projector by default. Pick an image, type a question, and start.

The projector is much faster with WebGPU than on the CPU. Add `?mode=cpu` or `?mode=webgpu` to the URL to pick the backend.

## Running without a browser

```shell
yzma model get -u https://huggingface.co/ggml-org/SmolVLM-256M-Instruct-GGUF/resolve/main/SmolVLM-256M-Instruct-Q8_0.gguf
yzma model get -u https://huggingface.co/ggml-org/SmolVLM-256M-Instruct-GGUF/resolve/main/mmproj-SmolVLM-256M-Instruct-Q8_0.gguf
make test-wasm-vlm
```

## Page calls

| Function | What it does |
| --- | --- |
| `yzmaLoadModel(modelURL, projectorURL)` | Downloads the model and the projector and creates the contexts. |
| `yzmaOpenModel(maxImageTokens)` | Loads files that are already in the module's filesystem. A value of 0 uses the model's limits. |
| `yzmaDescribe(prompt, width, height, rgba, maxTokens)` | Answers a question about the RGBA pixels of an image. |
