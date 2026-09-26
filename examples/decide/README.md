# Decide

Uses `yzma` to run a System One model, Jev-Style or JevK5. The model answers a typed question about a state with calibrated probabilities instead of text. See the [exp/decide](../../exp/decide) package.

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

The `-readout jevk5` flag runs a [JevK5](https://huggingface.co/alibiserikbay/JevK5-GGUF) model. It reads each option from the logit of its letter and answers questions with up to 256 options. Its temperature is in the `jevk5_config.json` of the unquantized model repo.

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

Each GGUF file needs the config of its own version. The [JevK5](https://huggingface.co/alibiserikbay/JevK5/resolve/main/jevk5_config.json) config fits `jevk5-4b-v0.3-*.gguf`, and the [JevK5-9B](https://huggingface.co/alibiserikbay/JevK5-9B/resolve/main/jevk5_config.json) config fits `jevk5-9b-v0.3.3-*.gguf`.

## Install

```shell
go install ./examples/decide
```
