# wasm/decide

Uses `yzma` to run a System One model in a browser with the [exp/decide](../../../exp/decide) package. The model answers a typed question about a state with calibrated probabilities instead of text. It takes Jev-Style, JevK5 and decider models, the same as [examples/decide](../../decide).

The page is [wasm/decide.html](../../../wasm/decide.html). See the Typed decisions section of [wasm/README.md](../../../wasm/README.md) for how it works.

## Building

Run these commands from the root of the repo.

```shell
make download-llama.cpp-wasm
make wasm-decide-example
```

This builds `build/wasm/yzma-decide.wasm` with TinyGo. `make wasm-example-go` builds it with the standard Go toolchain.

## Running

```shell
make serve-wasm
```

Open <http://localhost:8080/decide.html>. Pick a model and the page downloads it with its config file. The config holds the token ids and calibration temperatures, which are not in the GGUF file.

| Model | Config |
| --- | --- |
| [Jev-Style-0.8B-Decision-v3 Q4_K_M](https://huggingface.co/chaoliangUNSW/Jev-Style-0.8B-Decision-v3-GGUF) | `readout_config.json` |
| [JevK5 2B Q8_0](https://huggingface.co/alibiserikbay/JevK5-GGUF) | `jevk5_config.json` |
| [decider-0.8b Q4_K_M](https://huggingface.co/mradermacher/decider-0.8b-GGUF) | `decider_config.json` |

The page asks every question about the state at once with `DecideMany`.

## Running without a browser

```shell
yzma model get -u https://huggingface.co/chaoliangUNSW/Jev-Style-0.8B-Decision-v3-GGUF/resolve/main/Jev-Style-0.8B-Decision-v3-Q4_K_M.gguf
curl -L -o ~/models/readout_config.json https://huggingface.co/chaoliangUNSW/Jev-Style-0.8B-Decision-v3-GGUF/resolve/main/readout_config.json
make test-wasm-decide
```

## Page calls

| Function | What it does |
| --- | --- |
| `yzmaDecideLoad(modelURL, configURL, readout, manyMode)` | Downloads a model and its config. `readout` is `jev`, `jevk5` or `decider`. `manyMode` is `exact` or `batched`. |
| `yzmaDecideOpen(path, configJSON, readout, manyMode)` | Loads a model that is already in the module's filesystem. |
| `yzmaDecide(state, questionJSON, category)` | Scores one question. |
| `yzmaDecideMany(state, questionsJSON, category)` | Scores a JSON list of questions about one state. |
