# wasm/tools

Uses `yzma` to call tools in a browser. It renders the model's chat template with the [template](../../../pkg/template) package, offers the tools to the model, parses the tool calls with the [message](../../../pkg/message) package, runs them, and passes the results back to the model for a final answer.

The tools are in [tools.go](tools.go). `get_weather` returns the weather of a city and `calculate` does arithmetic on two numbers. Their answers depend only on the arguments.

The page is [wasm/tools.html](../../../wasm/tools.html). See the Tool calling section of [wasm/README.md](../../../wasm/README.md) for how it works.

## Building

Run these commands from the root of the repo.

```shell
make download-llama.cpp-wasm
make wasm-tools-example
```

This builds `build/wasm/yzma-tools.wasm` with TinyGo. `make wasm-example-go` builds it with the standard Go toolchain.

## Running

```shell
make serve-wasm
```

Open <http://localhost:8080/tools.html>. The page downloads [Qwen2.5-0.5B-Instruct Q8_0](https://huggingface.co/Qwen/Qwen2.5-0.5B-Instruct-GGUF) by default. This is about the smallest model that makes tool calls.

The model can call tools at most 3 times before it must answer.

## Running without a browser

```shell
yzma model get -u https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf
make test-wasm-tools
```

This model is too small to make a tool call, so the test only checks the round trip.

## Page calls

| Function | What it does |
| --- | --- |
| `yzmaLoadModel(url)` | Downloads a model and creates a context. |
| `yzmaOpenModel(path)` | Loads a model that is already in the module's filesystem. |
| `yzmaAsk(question, maxTokens)` | Answers a question and calls tools if it needs them. |
