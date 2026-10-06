# Decide

Uses `yzma` to run a System One model, Jev-Style, JevK5 or decider. The model answers a typed question about a state with calibrated probabilities instead of text. See the [exp/decide](../../exp/decide) package.

## Model

Download the model and its `readout_config.json` from [Jev-Style-0.8B-Decision-v3-GGUF](https://huggingface.co/chaoliangUNSW/Jev-Style-0.8B-Decision-v3-GGUF). The config holds the token ids and calibration temperatures, which are not in the GGUF file.

```shell
yzma model get -u https://huggingface.co/chaoliangUNSW/Jev-Style-0.8B-Decision-v3-GGUF/resolve/main/Jev-Style-0.8B-Decision-v3-Q4_K_M.gguf
curl -L -o ~/models/readout_config.json https://huggingface.co/chaoliangUNSW/Jev-Style-0.8B-Decision-v3-GGUF/resolve/main/readout_config.json
```

## Running

A choice question with described options.

```shell
$ go run ./examples/decide/ -model ~/models/Jev-Style-0.8B-Decision-v3-Q4_K_M.gguf -config ~/models/readout_config.json \
    -state '{"ticket": "I was charged twice for my subscription."}' \
    -question "Which team handles this?" \
    -options '{"billing": "payments, invoices", "technical": "bugs"}' \
    -category theme_routing
{
  "answer": "billing",
  "options": [
    "billing",
    "technical"
  ],
  "probabilities": [
    0.9940535133194868,
    0.0059464866805131154
  ],
  "scores": [
    2.573061943054199,
    -1.931929588317871
  ],
  "temperature": 0.8800546821789332,
  "top_probability": 0.9940535133194868,
  "entropy_concentration": 0.9474797747041599,
  "input_tokens": 61,
  "head_tokens": 44
}
```

A true or false question leaves out `-options`.

```shell
$ go run ./examples/decide/ -model ~/models/Jev-Style-0.8B-Decision-v3-Q4_K_M.gguf -config ~/models/readout_config.json \
    -state "The meeting moved from Tuesday to Thursday at 3pm." \
    -question "The meeting is on Thursday." -category mac_gate
```

A score question takes a list of levels, level 0 first.

```shell
$ go run ./examples/decide/ -model ~/models/Jev-Style-0.8B-Decision-v3-Q4_K_M.gguf -config ~/models/readout_config.json \
    -state "The reply was polite but did not answer the question at all." \
    -question "How helpful was the reply?" -type score \
    -options '["not helpful", "somewhat helpful", "very helpful"]'
```

The `-category` flag picks a fitted calibration temperature, for example `general_sentiment`, `theme_routing`, `intent` or `mac_gate`. Without it the global temperature is used.

## JevK5

The `-readout jevk5` flag runs a [JevK5](https://huggingface.co/alibiserikbay/JevK5-GGUF) model. It scores each option by the logit of its letter and answers questions with up to 256 options. Its temperature is in the `jevk5_config.json` in the unquantized model repo.

```shell
yzma model get -u https://huggingface.co/alibiserikbay/JevK5-GGUF/resolve/main/jevk5-2b-v0.2-Q8_0.gguf
curl -L -o ~/models/jevk5-2b_config.json https://huggingface.co/alibiserikbay/JevK5-2B/resolve/main/jevk5_config.json
```

```shell
$ go run ./examples/decide/ -readout jevk5 -model ~/models/jevk5-2b-v0.2-Q8_0.gguf -config ~/models/jevk5-2b_config.json \
    -state "I was billed twice for order #4411. Please refund the duplicate charge today." \
    -question "Which team should handle this?" \
    -options '{"billing": "Payments and refunds", "tech": "Bugs", "sales": "New purchases"}'
{
  "answer": "billing",
  "options": [
    "billing",
    "tech",
    "sales"
  ],
  "probabilities": [
    0.9929183710323725,
    0.0036116651697149305,
    0.003469963797912578
  ],
  "scores": [
    26.563657760620117,
    18.5882568359375,
    18.531421661376953
  ],
  "temperature": 1.42,
  "top_probability": 0.9929183710323725,
  "entropy_concentration": 0.9572009781782457,
  "input_tokens": 129
}
```

Each GGUF file needs the config for its version. The [JevK5](https://huggingface.co/alibiserikbay/JevK5/resolve/main/jevk5_config.json) config fits `jevk5-4b-v0.3-*.gguf`, and the [JevK5-9B](https://huggingface.co/alibiserikbay/JevK5-9B/resolve/main/jevk5_config.json) config fits `jevk5-9b-v0.3.3-*.gguf`.

## decider

The `-readout decider` flag runs a [decider](https://huggingface.co/Mapika/decider-0.8b) model. It scores each option by the logit of its letter and answers questions with up to 255 options. Its temperature is in the `decider_config.json` in the unquantized model repo.

```shell
yzma model get -u https://huggingface.co/mradermacher/decider-0.8b-GGUF/resolve/main/decider-0.8b.Q4_K_M.gguf
curl -L -o ~/models/decider-0.8b_config.json https://huggingface.co/Mapika/decider-0.8b/resolve/main/decider_config.json
```

```shell
$ go run ./examples/decide/ -readout decider -model ~/models/decider-0.8b.Q4_K_M.gguf -config ~/models/decider-0.8b_config.json \
    -state '{"ticket": "I was charged twice for order A-104. Please refund the duplicate."}' \
    -question "Which team should handle this?" \
    -options '{"billing": "Charges, invoices, refunds", "technical": "Bugs, outages", "other": ""}'
{
  "answer": "billing",
  "options": [
    "billing",
    "technical",
    "other"
  ],
  "probabilities": [
    0.9819664239853243,
    0.00021987808248419126,
    0.017813697932191414
  ],
  "scores": [
    16.9962100982666,
    8.33984375,
    12.8663330078125
  ],
  "temperature": 1.03,
  "top_probability": 0.9819664239853243,
  "entropy_concentration": 0.916738883989616,
  "input_tokens": 64
}
```

A score question judges each level on its own and has no `scores`.

## WebAssembly

The same package runs in a browser. See `examples/wasm/decide` and the typed decisions section of [wasm/README.md](../../wasm/README.md).

## Install

```shell
go install ./examples/decide
```

## gguf

The `-readout gguf` flag runs a model whose GGUF holds its readout, prompt template and temperatures, as converted for the llama-server `/v1/systemone` API. It needs no `-config`. [Lev](https://huggingface.co/ggml-org/lev-GGUF) and [OpenJev](https://huggingface.co/ggml-org/OpenJev-GGUF) are supported.

```shell
yzma model get -u https://huggingface.co/ggml-org/lev-GGUF/resolve/main/lev-Q8_0.gguf
```

```shell
$ go run ./examples/decide/ -readout gguf -model ~/models/lev-Q8_0.gguf \
    -state '{"ticket": "I was charged twice for order A-104. Please refund the duplicate."}' \
    -question "Which team should handle this?" \
    -options '{"billing": "Charges, invoices, refunds", "technical": "Bugs, outages", "other": ""}'
{
  "answer": "billing",
  "options": [
    "billing",
    "technical",
    "other"
  ],
  "probabilities": [
    0.868248902995133,
    0.034811167574575926,
    0.0969399294302912
  ],
  "temperature": 1.789783,
  "top_probability": 0.868248902995133,
  "entropy_concentration": 0.5760304077230565,
  "confidence": 0.8023733544926995,
  "expected": 0,
  "input_tokens": 256
}
```

Lev reads each choice question twice, the second time with the options in reverse order, so `input_tokens` counts both.
