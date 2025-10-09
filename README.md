# yzma

[![Linux](https://github.com/hybridgroup/yzma/actions/workflows/linux.yml/badge.svg)](https://github.com/hybridgroup/yzma/actions/workflows/linux.yml) [![macOS](https://github.com/hybridgroup/yzma/actions/workflows/macos.yml/badge.svg)](https://github.com/hybridgroup/yzma/actions/workflows/macos.yml) [![Windows](https://github.com/hybridgroup/yzma/actions/workflows/windows.yml/badge.svg)](https://github.com/hybridgroup/yzma/actions/workflows/windows.yml)

`yzma` lets you use Go to perform local inference with Vision Language Models (VLMs), Large Language Models (LLMs), Small Language Models (SLMs), and Tiny Language Models (TLMs) by using the [`llama.cpp`](https://github.com/ggml-org/llama.cpp) libraries all running on your own hardware.

It uses the [`purego`](https://github.com/ebitengine/purego) and [`ffi`](https://github.com/JupiterRider/ffi) packages so calls can be made directly to `llama.cpp` without CGo.

This example uses the [SmolLM-135M](https://huggingface.co/QuantFactory/SmolLM-135M-GGUF) model:

```go
package main

import (
	"fmt"
	"os"

	"github.com/hybridgroup/yzma/pkg/llama"
)

var (
	modelFile            = "./models/SmolLM-135M.Q2_K.gguf"
	prompt               = "Are you ready to go?"
	libPath              = os.Getenv("YZMA_LIB")
	responseLength int32 = 12
)

func main() {
	llama.Load(libPath)
	llama.Init()

	model := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	lctx := llama.InitFromModel(model, llama.ContextDefaultParams())

	vocab := llama.ModelGetVocab(model)

	// call once to get the size of the tokens from the prompt
	count := llama.Tokenize(vocab, prompt, nil, true, false)

	// now get the actual tokens
	tokens := make([]llama.Token, count)
	llama.Tokenize(vocab, prompt, tokens, true, false)

	batch := llama.BatchGetOne(tokens)

	sampler := llama.SamplerChainInit(llama.SamplerChainDefaultParams())
	llama.SamplerChainAdd(sampler, llama.SamplerInitGreedy())

	for pos := int32(0); pos+batch.NTokens < count+responseLength; pos += batch.NTokens {
		llama.Decode(lctx, batch)
		token := llama.SamplerSample(sampler, lctx, -1)

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

Produces the following output:

```shell
$ go run ./examples/hello/ 2>/dev/null

The first thing you need to do is to get your hands on a computer.
```

What's with the `2>/dev/null` at the end? That is the "easy way" to suppress the logging from `llama.cpp`.

## Installation

You will need to download the `llama.cpp` libraries for your platform. You can obtain them from https://github.com/ggml-org/llama.cpp/releases

Extract the library files into a directory on your local machine.

For Linux, they have the `.so` file extension. For example, `libllama.so`, `libmtmd.so` and so on. When using macOS, they have a `.dylib` file extension. And on Windows, they have a `.dll` file extension. You do not need the other downloaded files to use the `llama.cpp` libraries with `yzma`.

***Important Note***
You currently need to set both the `LD_LIBRARY_PATH` and the `YZMA_LIB` env variable to the directory with your llama.cpp library files. For example:

```shell
export LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:/home/ron/Development/yzma/lib
export YZMA_LIB=/home/ron/Development/yzma/lib
```

## Examples

### Vision Language Model (VLM) multimodal example

This example uses the [`Qwen2.5-VL-3B-Instruct-Q8_0`](https://huggingface.co/ggml-org/Qwen2.5-VL-3B-Instruct-GGUF) VLM model to process both a text prompt and an image, then displays the result.

```shell
$ go run ./examples/vlm/ -model ./models/Qwen2.5-VL-3B-Instruct-Q8_0.gguf -proj ./models/mmproj-Qwen2.5-VL-3B-Instruct-Q8_0.gguf -image ./images/domestic_llama.jpg -prompt "What is in this picture?" 2>/dev/null
Loading model ./models/Qwen2.5-VL-3B-Instruct-Q8_0.gguf
encoding image slice...
image slice encoded in 966 ms
decoding image batch 1/1, n_tokens_batch = 910
image decoded (batch 1/1) in 208 ms

The picture shows a white llama standing in a fenced-in area, possibly a zoo or a wildlife park. The llama is the main focus of the image, and it appears to be looking to the right. The background features a grassy area with trees and a fence, and there are some vehicles visible in the distance.
```

[See the code here](./examples/vlm/main.go).

### Small Language Model (SLM) interactive chat example

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

## More info

`yzma` is still a work in progress but it has support for some basic functionality.

You can already use VLMs and other language models with full hardware acceleration.

Here are some advantages of `yzma` over other Go packages for `llama.cpp`:

- Compile Go programs that use `yzma` with the normal `go build` and `go run` commands. No C compiler needed!
- Use the `llama.cpp` libraries with whatever hardware acceleration is available for your configuration. CUDA, Vulkan, etc.
- Download `llama.cpp` precompiled libraries directly from Github, or include them with your application.
- Update the `llama.cpp` libraries without recompiling your Go program, as long as `llama.cpp` does not make any breaking changes.

The idea is to make it easier for Go developers to use language models as part of "normal" applications without having to use containers or do anything other than the normal `GOOS` and `GOARCH` env variables for cross-complication.

`yzma` borrows definitions from the https://github.com/dianlight/gollama.cpp package then modifies them rather heavily. Thank you!
