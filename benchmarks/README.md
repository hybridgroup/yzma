# Benchmarks

## Latest results

Want to see the latest results? These numbers change with each llama.cpp release.

| File | What is in it |
| --- | --- |
| [linux.md](linux.md) | Linux, amd64 and arm64, CPU and GPU |
| [macos.md](macos.md) | macOS with Metal |
| [windows.md](windows.md) | Windows, CPU and GPU |
| [webassembly.md](webassembly.md) | WebAssembly in Node and in a browser |
| [comparison.md](comparison.md) | yzma against ollama and Docker Model Runner |

## Before running them yourself

Get the library and the models.

```shell
make download-llama.cpp
make download-benchmark-models
```

## Run them

On Linux and macOS.

```shell
./benchmarks/run.sh
```

On Windows, use a PowerShell prompt in the repository directory. A Command
Prompt opens the file in an editor instead of running it. PowerShell won't run a
script unless the machine's execution policy allows it, so this command sets the
policy for one run.

```powershell
powershell -ExecutionPolicy Bypass -File .\benchmarks\run.ps1
```

A full run takes a while. The first go command builds the packages, and each
suite takes many minutes. The script prints each command and each benchmark line
as it arrives.

The script asks llama.cpp which devices the machine has, runs the text and
multimodal benchmarks on each one, and writes each result to the platform's
file. The llama.cpp build tag comes from `yzma-install.json` in the library
directory.

Useful flags.

```shell
./benchmarks/run.sh --backend vulkan        # one backend only
./benchmarks/run.sh --suite text            # one suite only
./benchmarks/run.sh --machine jetson-orin-nano --label "Jetson Orin Nano 8GB"
./benchmarks/run.sh --llamacpp b10964       # when the library came from elsewhere
./benchmarks/run.sh --threads 24            # a custom thread count
./benchmarks/run.sh --threadpool            # pin each thread to a core
./benchmarks/run.sh --dry-run               # print the result, change no file
```

The text benchmark uses the thread count from `llama.ModelThreads`, which is 4
for SmolLM-135M on any machine with 4 or more cores. The model is too small to
use more threads. Each token needs little arithmetic, and the threads wait for
each other after each operation, so more threads make it slower. On an Apple M4
Pro, 4 threads reach more than 900 tokens a second and 10 threads are much
slower. With the same count on every machine, the text tables measure the same
work, and they match the older rows, which used the llama.cpp default of 4.

The multimodal benchmark uses one thread per performance core, as a yzma
program does. Its model does more work per token, so this suite shows what a
large processor can do.

Use `--threads` to try another count for both suites. `--threads 0` uses the
yzma default. For text, that count comes from the model size. For multimodal,
it is one thread per performance core.

`--threadpool` pins each thread to its own performance CPU. Without it the
system moves threads around during the run, which makes a short run read low
and vary from run to run. It works only on Linux. macOS and Windows don't report
which CPUs are performance CPUs, so the benchmark fails with an error there.

The PowerShell script takes the same flags with one dash and a capital
letter, such as `-Machine`, `-Backend` and `-DryRun`. The flags go after the
file name.

```powershell
powershell -ExecutionPolicy Bypass -File .\benchmarks\run.ps1 -Backend vulkan -DryRun
```

The machine name is part of each section key. Use the same name every time,
or the file gets two sections for one machine. The default is the host name.

## Comparing engines

This suite measures yzma against a model server. yzma calls llama.cpp in the
same process. ollama and Docker Model Runner answer over an OpenAI compatible
REST interface. The suite shows what that round trip costs.

It is a separate script because it has different needs. The servers must be
running, and each one must have the model.

```shell
make download-compare-models
docker model pull hf.co/qwen/qwen3-vl-4b-instruct-gguf:q4_k_m
./benchmarks/compare.sh
```

Each engine must read the same GGUF file. So the commands use the Hugging
Face file and not `gemma4:e4b` or `ai/gemma3`, which are vendor conversions. The
script tells you which command fetches a missing model.

Two things about the servers.

- Docker Model Runner model references must be lower case. An upper case
  reference fails with "Invalid model reference".
- ollama 0.34 doesn't follow the Hugging Face redirect to its CDN, so
  `ollama pull hf.co/...` fails. The script gives ollama the file from
  `MODELS_DIR` instead, which makes sure it reads the same bytes. Mount that
  directory in the ollama container.

```shell
docker run -d --gpus all -v ollama:/root/.ollama -v "$HOME/models:/models:ro" \
  -p 11434:11434 --name ollama ollama/ollama:latest
```

`OLLAMA_CONTAINER` names that container, and `OLLAMA_MODELS_DIR` is where the
models are inside it. Set `OLLAMA_CONTAINER=` when ollama runs on the host.

Useful flags. They match the `run.sh` flags where the meaning is the same.

```shell
./benchmarks/compare.sh --engine yzma            # one engine only
./benchmarks/compare.sh --suite embeddings       # one suite only
./benchmarks/compare.sh --model qwen3-vl-4b      # one model only
./benchmarks/compare.sh --dry-run                # print the result, change no file
```

| Model | Short name | GGUF |
| --- | --- | --- |
| bge-small-en-v1.5 Q8_0 | `bge-small` | `ggml-org/bge-small-en-v1.5-Q8_0-GGUF` |
| Qwen3-VL-2B-Instruct Q4_K_M | `qwen3-vl-2b` | `Qwen/Qwen3-VL-2B-Instruct-GGUF` |
| Gemma 4 E2B Q4_K_M | `gemma4-e2b` | `unsloth/gemma-4-E2B-it-GGUF` |
| SmolVLM-256M-Instruct Q8_0 | `smolvlm-256m` | `ggml-org/SmolVLM-256M-Instruct-GGUF` |

`smolvlm-256m` is the model the other suites use. Use it for a quick check of
the engines, since it needs no large download.

### The llama.cpp build of each engine

Each engine ships its own llama.cpp build, so the same GGUF file does not mean
the same work.

| Engine | Build |
| --- | --- |
| yzma | the library of `lib/`, see `yzma-install.json` |
| Docker Model Runner | pinned in its image, b9879 in version 1.2 |
| ollama | its own fork |

Ask the server which build it has. Docker Model Runner returns it in
`system_fingerprint`.

```shell
curl -s http://localhost:12434/engines/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -d '{"model":"ai/gemma3:270M-F16","messages":[{"role":"user","content":"hi"}],"max_tokens":1}' |
  sed -n 's/.*"system_fingerprint":"\([^"]*\)".*/\1/p'
```

The Docker Model Runner runner is a container. `docker model install-runner`
does nothing when that container is already running, so a new plugin can keep an
old engine. To update it, remove it first. The models stay, because they live in
their own volume.

```shell
docker model uninstall-runner
docker model install-runner
```

A new runner can run as a different user than the old one. If the models then
fail with "permission denied", give the volume to the new runner's user.

```shell
docker run --rm -v docker-model-runner-models:/models alpine chown -R 995:995 /models
```

### ollama and chat templates

ollama doesn't use the chat template from the GGUF model it imports. It gives
the model `TEMPLATE {{ .Prompt }}`, which sends the prompt text with no user turn
and no assistant turn. The model then answers differently from the other
engines.

For an architecture ollama knows, such as gemma4, it uses a `RENDERER`
instead, which frames the chat correctly. For one it doesn't know, such as
SmolVLM's idefics3, the template in `modelFiles` in `compare.sh` provides the
framing.

The script refuses a model that has no renderer and no template, so this
fault can't slip into a table unnoticed.

Both models take text and images, so one model covers both suites. The server
addresses come from `OLLAMA_URL` and `DMR_URL`.

This script is for Linux and macOS. There is no PowerShell twin yet.

After a run, read what the script prints. It shows each engine's answer and
prompt token count.

Each run gets an image and a prompt that no engine has seen. A server answers a
repeated request from its cache, 40 times faster in one measurement, while yzma
clears its cache after each generation. Otherwise the benchmark would compare a
warm server with a cold yzma.

Request numbers keep counting across each `-count` of a run, so no count repeats
the requests of the one before. The warmup request is not timed and no timed
request repeats it. Before this, ollama answered the second count's images from
its cache in 190 ms, against 540 ms for a new image.

The prompt token count is the check that matters. Two engines can give similar
answers and still do different work. A different count always means a different
prompt, or different image preprocessing. The script flags it, and those numbers
must not go in a table.

One gap is known and accepted. ollama counts 5 more tokens per image than yzma
and Docker Model Runner, on every model and at every image size, while the
answers match. That is less than 3 percent of the image suite prompt.

The image makes up most of the prompt tokens. A projector that cuts the image
into tiles produces three times the tokens of one that doesn't, so the vision
model does three times the work. Use `-image-min-tokens` and `-image-max-tokens`
to set that budget for yzma.

Each engine scales images its own way. ollama scales a small Qwen3-VL image up
to about 1000 tokens. The llama.cpp in Docker Model Runner scales a Gemma 4
image down to 280 tokens, and a newer build does not. So the script gives each
model an image size that no engine changes, 1280x960 for qwen3-vl-2b and 768x576
for gemma4-e2b. `-image-size` sets it.

The multimodal suite heats up a laptop GPU, and a hot GPU lowers its clock.
Before each engine the script waits until the GPU is at 60 degrees or less. Pass
`--cool` a different limit for a machine that idles warmer.

Every engine gets the image before the text. The prompt token count doesn't
show the order, but the model sees a different prompt.

The code is in [benchmarks/compare](compare). The conditions that keep the
engines on equal terms are in [comparison.md](comparison.md).

## WebAssembly

Node covers the build with one thread and the build with more threads.

```shell
make download-llama.cpp-wasm
make wasm-example
./benchmarks/run.sh --backend wasm
```

WebGPU needs a browser, which no script here can drive. Serve the example, open
the page, paste [browser-bench.js](browser-bench.js) into the console, and pass
the result to the tool. On a machine with two GPUs, set `gpu` in the script and
`--device` in the tool. The output includes the llama.cpp build, which the page
reads from `yzma-install.json` in the build directory.

```shell
make serve-wasm
go run ./cmd/yzma-bench update --file benchmarks/webassembly.md \
  --suite browser --backend webgpu --arch wasm \
  --machine <name> --device <gpu> --label "<machine and browser>" --output run.txt
```

## How results are saved

Each result is a section between markers, with the key
`<suite>/<backend>/<arch>/<machine>[/<device>]`. A new run with the same key
replaces only that section. Each suite's table is a generated block that
[cmd/yzma-bench](../cmd/yzma-bench) rebuilds from the sections, so a table and
its sections can't disagree.

Do not edit a table by hand. To check a file after an edit:

```shell
go run ./cmd/yzma-bench check benchmarks/*.md
```

Each section must record which llama.cpp build produced it. The tool refuses a
result with no tag, so the numbers in a table always compare the same thing.

To delete a result that is no longer comparable, for example after a model
change, pass its key to `remove`.

```shell
go run ./cmd/yzma-bench remove --file benchmarks/linux.md \
  multimodal/cuda/amd64/i9-13900hx/cuda0
```

## What is measured

The benchmarks report tokens a second with `b.ReportMetric`, and the table
shows the median of the runs. The comparison suite reports two more metrics,
time to first token and total request time.

| Suite | Code | Model |
| --- | --- | --- |
| text | [pkg/llama/benchmark_test.go](../pkg/llama/benchmark_test.go) | SmolLM-135M Q2_K |
| compare-text, compare-multimodal | [benchmarks/compare](compare) | Qwen3-VL-4B, Gemma 4 E4B |
| multimodal | [pkg/mtmd/benchmark_test.go](../pkg/mtmd/benchmark_test.go) | SmolVLM-256M-Instruct Q8_0 with its projector |
| node, browser | [examples/wasm/chat](../examples/wasm/chat) | SmolLM-135M Q2_K |

The WebAssembly numbers come from the example's generation loop, so they are
not comparable with the native tables.
