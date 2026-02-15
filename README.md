![yzma logo](./images/yzma_logo.png)

# yzma

[![Go Reference](https://pkg.go.dev/badge/github.com/hybridgroup/yzma.svg)](https://pkg.go.dev/github.com/hybridgroup/yzma) [![Linux](https://github.com/hybridgroup/yzma/actions/workflows/linux.yml/badge.svg)](https://github.com/hybridgroup/yzma/actions/workflows/linux.yml) [![macOS](https://github.com/hybridgroup/yzma/actions/workflows/macos.yml/badge.svg)](https://github.com/hybridgroup/yzma/actions/workflows/macos.yml) [![Windows](https://github.com/hybridgroup/yzma/actions/workflows/windows.yml/badge.svg)](https://github.com/hybridgroup/yzma/actions/workflows/windows.yml) [![llama.cpp Release](https://img.shields.io/github/v/release/ggml-org/llama.cpp?label=llama.cpp)](https://github.com/ggml-org/llama.cpp/releases)

`yzma` lets you write Go applications that directly integrate [`llama.cpp`](https://github.com/ggml-org/llama.cpp) for fully local inference using hardware acceleration.

Run the latest Vision Language Models and Large/Small/Tiny Language Models on Linux, macOS, or Windows. Use hardware acceleration such as CUDA, Metal, or Vulkan. `yzma` uses the [`purego`](https://github.com/ebitengine/purego) and [`ffi`](https://github.com/JupiterRider/ffi) packages so CGo is not needed. This also means that `yzma` always works with the newest `llama.cpp` releases, so you can use the latest features and models.

This example uses the [SmolLM2-135M-Instruct](https://huggingface.co/bartowski/SmolLM2-135M-Instruct-GGUF) model:

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hybridgroup/yzma/pkg/download"
	"github.com/hybridgroup/yzma/pkg/llama"
)

var (
	modelFile            = "SmolLM2-135M.Q4_K_M.gguf"
	prompt               = "Are you ready to go?"
	libPath              = os.Getenv("YZMA_LIB")
	responseLength int32 = 12
)

func main() {
	llama.Load(libPath)
	llama.LogSet(llama.LogSilent())

	llama.Init()

	model, _ := llama.ModelLoadFromFile(filepath.Join(download.DefaultModelsDir(), modelFile), llama.ModelDefaultParams())
	ctx, _ := llama.InitFromModel(model, llama.ContextDefaultParams())

	vocab := llama.ModelGetVocab(model)

	tokens := llama.Tokenize(vocab, prompt, true, false)

	batch := llama.BatchGetOne(tokens)

	sampler := llama.SamplerChainInit(llama.SamplerChainDefaultParams())
	llama.SamplerChainAdd(sampler, llama.SamplerInitGreedy())

	for pos := int32(0); pos < responseLength; pos += batch.NTokens {
		llama.Decode(ctx, batch)
		token := llama.SamplerSample(sampler, ctx, -1)

		if llama.VocabIsEOG(vocab, token) {
			fmt.Println()
			break
		}

		buf := make([]byte, 36)
		len := llama.TokenToPiece(vocab, token, buf, 0, true)

		fmt.Print(string(buf[:len]))

		batch = llama.BatchGetOne([]llama.Token{token})
	}

	fmt.Println()
}
```

Download the model using the `yzma` command line tool:

```shell
yzma model get -u https://huggingface.co/bartowski/SmolLM2-135M-Instruct-GGUF/resolve/main/SmolLM2-135M-Instruct-Q4_K_M.gguf
```

And run the Go program:

```shell
$ go run ./examples/hello/


"Yes, I'm ready to go."
```

## Installation

You can use the convenient `yzma` command line tool to download the `llama.cpp` prebuilt libraries for your platform, download them manually, or even have your application download them automatically.

See [INSTALL.md](./INSTALL.md) for installation instructions for macOS, Linux, and Windows. We also have specific information on running `yzma` on Raspberry Pi or NVIDIA Jetson Orin.

## Examples

### Vision Language Model (VLM) Multimodal Example

This example uses the [`Qwen2.5-VL-3B-Instruct-Q8_0`](https://huggingface.co/ggml-org/Qwen2.5-VL-3B-Instruct-GGUF) VLM model to process both a text prompt and an image, then displays the result.

```shell
$ go run ./examples/vlm/ -model ~/models/Qwen2.5-VL-3B-Instruct-Q8_0.gguf -mmproj ~/models/mmproj-Qwen2.5-VL-3B-Instruct-Q8_0.gguf -image ./images/domestic_llama.jpg -p "What is in this picture?"

The image features a white llama standing in a fenced-in area, possibly a zoo or a farm. The llama is positioned in the center of the image, with its body facing the right side. The fenced area is surrounded by trees, creating a natural environment for the llama.
```

[See the code here](./examples/vlm/main.go).

### Small Language Model (SLM) Interactive Chat Example

You can use `yzma` to do inference on text language models. This example uses the [`qwen2.5-0.5b-instruct-fp16.gguf `](https://huggingface.co/Qwen/Qwen2.5-0.5B-Instruct-GGUF) model for an interactive chat session.

```shell
$ go run ./examples/chat/ -model ./models/qwen2.5-0.5b-instruct-fp16.gguf
Enter prompt: Are you ready to go?

Yes, I'm ready to go! What would you like to do?

Enter prompt: Let's go to the zoo


Great! Let's go to the zoo. What would you like to see?

Enter prompt: I want to feed the llama


Sure! Let's go to the zoo and feed the llama. What kind of llama are you interested in feeding?
```

[See the code here](./examples/chat/main.go).

### Additional Examples

See the [examples](./examples/) directory for more examples of how to use `yzma`.

## Projects

Who is using yzma? Check out some of the [tools, applications, examples, and blog posts](./PROJECTS.md) here!

## Models

`yzma` uses models in the GGUF format supported by `llama.cpp`. There are many models in GGUF format on Hugging Face (over 161k at last count):

https://huggingface.co/models?library=gguf&sort=trending

You can use the `yzma` command to download models for you! 

For example, this downloads the `gemma-3-1b-it-GGUF` model:

```
yzma model get -u https://huggingface.co/ggml-org/gemma-3-1b-it-GGUF/resolve/main/gemma-3-1b-it-Q4_K_M.gguf

```

Check out the [Model Usage](./MODELS.md) page for more information.

## Support

`yzma` currently has support for over 95% of `llama.cpp` functionality. See [ROADMAP.md](./ROADMAP.md) for the complete list.

You can use multimodal models (image/audio) and text language models with full hardware acceleration on Linux, macOS, and Windows.

| OS      | CPU          | GPU                             |
| ------- | ------------ | ------------------------------- |
| Linux   | amd64, arm64 | CUDA, Vulkan, HIP, ROCm, SYCL   |
| macOS   | arm64        | Metal                           |
| Windows | amd64        | CUDA, Vulkan, HIP, SYCL, OpenCL |

Whenever there is a new release of `llama.cpp`, the tests for `yzma` are run automatically. This helps us stay up to date with the latest code and models.

## Benchmarks

`yzma` is fast because it calls `llama.cpp` in the same process. No external servers needed!

For example, here is the `Qwen3-VL-2B-Instruct` Visual Language Model (VLM) performing multi-modal inference on an image and text prompt running on a Apple M4 Max with 128 GB RAM:

```shell
$ go test -run none -benchtime=10s -count=5 -bench BenchmarkMultimodalInference
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/mtmd
cpu: Apple M4 Max
BenchmarkMultimodalInference-16		10		1577948683 ns/op	788.9 tokens/s
BenchmarkMultimodalInference-16		12		1243692014 ns/op	910.8 tokens/s
BenchmarkMultimodalInference-16		 7		1654741804 ns/op	737.2 tokens/s
BenchmarkMultimodalInference-16		 7		1568106947 ns/op	771.9 tokens/s
BenchmarkMultimodalInference-16		10		1704669371 ns/op	706.1 tokens/s
PASS
ok  	github.com/hybridgroup/yzma/pkg/mtmd	76.644s
```

Want to see more benchmarks? Take a look at the [BENCHMARKS.md](./BENCHMARKS.md) document.

## More Info

`yzma` is now ready to be used to build complete applications that incorporate language models directly into your Golang code.

Here are some advantages of `yzma` with `llama.cpp`:

- Compile Go programs that use `yzma` with the normal `go build` and `go run` commands. No C compiler needed!
- Use the `llama.cpp` libraries with whatever hardware acceleration is available for your configuration. CUDA, Vulkan, etc.
- High performance from making function calls from within the same process. No external model servers!
- Download `llama.cpp` precompiled libraries directly from Github, or include them with your application.
- Update the `llama.cpp` libraries without recompiling your Go program, as long as `llama.cpp` does not make any breaking changes.

The idea is to make it easier for Go developers to use language models as part of "normal" applications without having to use containers or do anything other than the normal `GOOS` and `GOARCH` env variables for cross-complication.

`yzma` originally started with definitions from the https://github.com/dianlight/gollama.cpp package, but then has gone on to modify them rather heavily. Thank you!
