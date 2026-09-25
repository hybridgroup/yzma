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

On Windows. Use a PowerShell prompt in the directory of the repository. A
Command Prompt opens the file in an editor and does not run it. PowerShell does
not run a script until the policy of the machine permits it. This command gives
the policy for one run.

```powershell
powershell -ExecutionPolicy Bypass -File .\benchmarks\run.ps1
```

A full run is long. The first go command builds the packages, and each suite
takes many minutes. The script shows each command and each line of the
benchmarks when they come.

The script asks llama.cpp which devices the machine has, runs the text
benchmark and the multimodal benchmark for each one, and puts each result in the
file of the platform. The tag of the llama.cpp build comes from
`yzma-install.json` of the library directory.

Useful flags.

```shell
./benchmarks/run.sh --backend vulkan        # one backend only
./benchmarks/run.sh --suite text            # one suite only
./benchmarks/run.sh --machine jetson-orin-nano --label "Jetson Orin Nano 8GB"
./benchmarks/run.sh --llamacpp b10964       # when the library came from elsewhere
./benchmarks/run.sh --threads 24            # a thread count of your own
./benchmarks/run.sh --threadpool            # hold each thread to a core
./benchmarks/run.sh --dry-run               # print the result, change no file
```

The text benchmark uses the thread count of `llama.ModelThreads`, which is 4 for
SmolLM-135M on each machine with 4 or more cores. The model is too small to
use more threads. Each token has little arithmetic, and the threads wait for
each other after each operation, thus more threads make it slower. On an Apple
M4 Pro, 4 threads give more than 900 tokens a second and 10 threads give much
less. With the same count on each machine, the text tables measure the same
work, and they agree with the older rows, which used the 4 threads of llama.cpp.

The multimodal benchmark uses one thread for each performance core of the
machine, as a program of yzma does. Its model does more work for each token,
thus this suite shows what a large processor can do.

Use `--threads` to try another count for both suites. `--threads 0` gives the
default of yzma. For text, this count comes from the model size. For
multimodal, it is one thread for each performance core.

`--threadpool` holds each thread to a performance CPU of its own. Without it the
system moves the threads while the work goes on, which makes a short run read
low and gives a different answer each time. It works only on Linux. macOS and
Windows do not say which CPUs are performance CPUs, thus the benchmark stops
with an error there.

The PowerShell script takes the same names with one dash and a capital, as
`-Machine`, `-Backend`, `-DryRun` and so on. The flags go after the name of the
file.

```powershell
powershell -ExecutionPolicy Bypass -File .\benchmarks\run.ps1 -Backend vulkan -DryRun
```

The machine name is part of the key of a section. Give the same name each time,
or the file gets two sections for one machine. The default is the host name.

## Comparing engines to engines

This suite measures yzma against a model server. yzma calls llama.cpp in the
same process. ollama and Docker Model Runner answer over an OpenAI compatible
REST interface. The suite shows what that round trip costs.

It is a separate script, because its needs are different. The servers must run,
and each one must have the model.

```shell
make download-compare-models
docker model pull hf.co/qwen/qwen3-vl-4b-instruct-gguf:q4_k_m
./benchmarks/compare.sh
```

Each engine must read the same GGUF file. Thus the commands take the file of
Hugging Face and not `gemma4:e4b` or `ai/gemma3`, which are the conversions of a
vendor. The script says which command gets a model that is absent.

Two things about the servers.

- The reference of Docker Model Runner must be lower case. An upper case
  reference gives "Invalid model reference".
- ollama 0.34 does not follow the redirect of Hugging Face to its CDN, thus
  `ollama pull hf.co/...` fails. The script gives ollama the file of
  `MODELS_DIR` instead, which makes sure that it reads the same bytes. Mount
  that directory in the container of ollama.

```shell
docker run -d --gpus all -v ollama:/root/.ollama -v "$HOME/models:/models:ro" \
  -p 11434:11434 --name ollama ollama/ollama:latest
```

`OLLAMA_CONTAINER` names that container, and `OLLAMA_MODELS_DIR` says where the
models are inside it. Set `OLLAMA_CONTAINER=` when ollama runs on the host.

Useful flags. They are the flags of `run.sh` where they mean the same thing.

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

`smolvlm-256m` is the model of the other suites. Use it for a quick check of
the engines, because it needs no large download.

### The llama.cpp build of each engine

Each engine brings its own build of llama.cpp, thus the same GGUF file does not
give the same work.

| Engine | Build |
| --- | --- |
| yzma | the library of `lib/`, see `yzma-install.json` |
| Docker Model Runner | pinned in its image, b9879 in version 1.2 |
| ollama | its own fork |

Ask the server which build it has. The answer of Docker Model Runner gives it in
`system_fingerprint`.

```shell
curl -s http://localhost:12434/engines/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -d '{"model":"ai/gemma3:270M-F16","messages":[{"role":"user","content":"hi"}],"max_tokens":1}' |
  sed -n 's/.*"system_fingerprint":"\([^"]*\)".*/\1/p'
```

The runner of Docker Model Runner is a container. `docker model install-runner`
does nothing when that container runs already, thus a new plugin can keep an old
engine. To make it new again, remove it first. The models stay, because they are
in a volume of their own.

```shell
docker model uninstall-runner
docker model install-runner
```

A new runner can have a different user than the old one. When the models then
give "permission denied", give the volume to the user of the new runner.

```shell
docker run --rm -v docker-model-runner-models:/models alpine chown -R 995:995 /models
```

### ollama and chat templates

ollama does not use the chat template from the GGUF model that it imports. It gives
the model `TEMPLATE {{ .Prompt }}`, which sends the words of the prompt with no
User turn and no Assistant turn. The model then answers something else than the
other engines.

For an architecture that ollama knows, such as gemma4, it gives a `RENDERER`
instead, which frames the chat correctly. For an architecture that it does not
know, such as idefics3 of SmolVLM, the template of `modelFiles` in
`compare.sh` gives the frame.

The script refuses a model that has no renderer and no template, thus this
fault cannot reach a table without notice.

Both models take text and images, thus one model covers both suites. The
addresses of the servers come from `OLLAMA_URL` and `DMR_URL`.

This script is for Linux and macOS. There is no PowerShell twin yet.

After a run, read what the script prints. It gives the answer of each engine
and the count of the prompt tokens of each engine.

Each run gets an image and a prompt that no engine has seen. A server answers a
repeat of the same request from its cache, 40 times faster in one measurement,
while yzma empties its cache after each generation. The benchmark would else
compare a warm server with a cold yzma.

The number of a request counts across each `-count` of a run, thus no count
repeats the requests of the one before. The warmup request is not timed and
no timed request repeats it. Before this, ollama answered the images of the
second count from its cache in 190 ms, against 540 ms for a new image.

The count of the prompt tokens is the check that matters. Two engines can give
a similar answer and still do different work. A different count always means a
different prompt, or a different preprocessing of the image. The script says so
and the numbers must not go in a table.

An image gives most of the prompt tokens. A projector that cuts the image into
tiles gives three times the tokens of one that does not, thus it does three
times the work of the vision model. Use `-image-min-tokens` and
`-image-max-tokens` to set that budget for yzma.

Each engine scales an image in its own way. ollama scales a small Qwen3-VL image
up to about 1000 tokens. The llama.cpp of Docker Model Runner scales a Gemma 4
image down to 280 tokens, and a newer build does not. The script therefore
gives each model an image of a size that no engine changes, 1280x960 for
qwen3-vl-2b and 768x576 for gemma4-e2b. `-image-size` sets it.

The multimodal suite makes the GPU of a laptop hot, and a hot GPU lowers its
clock. Before each engine the script waits until the GPU is at 60 degrees or
less. Give `--cool` another limit for a machine that stays warmer when idle.

The image goes before the text in each engine. The count of the prompt tokens
does not show the order, but the model answers another prompt.

The code is in [benchmarks/compare](compare). The conditions that make the
engines equal are in [comparison.md](comparison.md).

## WebAssembly

Node covers the build with one thread and the build with more threads.

```shell
make download-llama.cpp-wasm
make wasm-example
./benchmarks/run.sh --backend wasm
```

WebGPU needs a browser, which no script here can drive. Serve the example, open
the page, paste [browser-bench.js](browser-bench.js) in the console, and give
the result to the tool. The output gives the llama.cpp build, which the page
reads from `yzma-install.json` of the build directory.

```shell
make serve-wasm
go run ./cmd/yzma-bench update --file benchmarks/webassembly.md \
  --suite browser --backend webgpu --arch wasm \
  --machine <name> --label "<machine and browser>" --output run.txt
```

## How results are saved

Each result is a section between markers, and the key is
`<suite>/<backend>/<arch>/<machine>[/<device>]`. A new run of the same key
replaces that section only. The table of each suite is a generated block, which
[cmd/yzma-bench](../cmd/yzma-bench) makes again from the sections. Thus a table
and its sections cannot disagree.

Do not edit a table by hand. To check a file after an edit:

```shell
go run ./cmd/yzma-bench check benchmarks/*.md
```

Each section must say which llama.cpp build made it. The tool refuses a result
that has no tag, thus the numbers of a table always compare the same thing.

To delete a result that is not comparable any more, for example after a change
of the model, give its key to `remove`.

```shell
go run ./cmd/yzma-bench remove --file benchmarks/linux.md \
  multimodal/cuda/amd64/i9-13900hx/cuda0
```

## What is measured

The benchmarks report tokens a second with `b.ReportMetric`, and the table gives
the median of the runs. The comparison suite reports two more metrics, the time
to the first token and the time of a whole request.

| Suite | Code | Model |
| --- | --- | --- |
| text | [pkg/llama/benchmark_test.go](../pkg/llama/benchmark_test.go) | SmolLM-135M Q2_K |
| compare-text, compare-multimodal | [benchmarks/compare](compare) | Qwen3-VL-4B, Gemma 4 E4B |
| multimodal | [pkg/mtmd/benchmark_test.go](../pkg/mtmd/benchmark_test.go) | SmolVLM-256M-Instruct Q8_0 with its projector |
| node, browser | [examples/wasm/chat](../examples/wasm/chat) | SmolLM-135M Q2_K |

The WebAssembly numbers come from the generation loop of the example, thus they
are not comparable with the native tables.
