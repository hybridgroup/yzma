# Get the absolute path of the current Makefile
MAKEFILE_PATH := $(realpath $(lastword $(MAKEFILE_LIST)))
MAKEFILE_DIR := $(dir $(MAKEFILE_PATH))
YZMA_LIB ?= $(MAKEFILE_DIR)lib
MODELS_DIR ?= $(HOME)/models
WASM_DIR ?= $(MAKEFILE_DIR)build/wasm

# make download-models to download the models used in tests and examples.
# make download-models MODELS_DIR=/path/to/models to specify a different directory for the models.
download-models:
	mkdir -p $(MODELS_DIR)
	yzma model get -y --show-progress=false -o $(MODELS_DIR) -u https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf
	yzma model get -y --show-progress=false -o $(MODELS_DIR) -u https://huggingface.co/ggml-org/SmolVLM-256M-Instruct-GGUF/resolve/main/SmolVLM-256M-Instruct-Q8_0.gguf
	yzma model get -y --show-progress=false -o $(MODELS_DIR) -u https://huggingface.co/ggml-org/SmolVLM-256M-Instruct-GGUF/resolve/main/mmproj-SmolVLM-256M-Instruct-Q8_0.gguf
	yzma model get -y --show-progress=false -o $(MODELS_DIR) -u https://huggingface.co/ggml-org/models-moved/resolve/main/jina-reranker-v1-tiny-en/ggml-model-f16.gguf
	yzma model get -y --show-progress=false -o $(MODELS_DIR) -u https://huggingface.co/callgg/t5-base-encoder-f32/resolve/main/t5base-encoder-q4_0.gguf
	yzma model get -y --show-progress=false -o $(MODELS_DIR) -u https://huggingface.co/deadprogram/yzma-tests/resolve/main/Gemma2-Base-F32.gguf
	yzma model get -y --show-progress=false -o $(MODELS_DIR) -u https://huggingface.co/deadprogram/yzma-tests/resolve/main/Gemma2-Lora-F32-LoRA.gguf
	yzma model get -y --show-progress=false -o $(MODELS_DIR) -u https://huggingface.co/ggml-org/models-moved/resolve/main/tinyllamas/split/stories15M-q8_0-00001-of-00003.gguf
	yzma model get -y --show-progress=false -o $(MODELS_DIR) -u https://huggingface.co/ggml-org/models-moved/resolve/main/tinyllamas/split/stories15M-q8_0-00002-of-00003.gguf
	yzma model get -y --show-progress=false -o $(MODELS_DIR) -u https://huggingface.co/ggml-org/models-moved/resolve/main/tinyllamas/split/stories15M-q8_0-00003-of-00003.gguf

clean-llama.cpp:
	rm -rf $(YZMA_LIB)/*

# make download-llama.cpp VERSION=b8080 to download a specific version of llama.cpp
download-llama.cpp:
	yzma install -lib $(YZMA_LIB) $(if $(VERSION),-v $(VERSION))

build:
	YZMA_LIB=$(YZMA_LIB) go build -o yzma .

install:
	go install .

# make test to run the tests. Make sure to run `make download-models` first to download the models used in the tests.
# make test MODELS_DIR=/path/to/models to specify a different directory for the models if you didn't use the default one when downloading the models.
test:
	export YZMA_LIB=$(YZMA_LIB) && \
	export YZMA_TEST_MODEL=$(MODELS_DIR)/SmolLM-135M.Q2_K.gguf && \
	export YZMA_TEST_MMMODEL=$(MODELS_DIR)/SmolVLM-256M-Instruct-Q8_0.gguf && \
	export YZMA_TEST_MMPROJ=$(MODELS_DIR)/mmproj-SmolVLM-256M-Instruct-Q8_0.gguf && \
	export YZMA_TEST_QUANTIZE_MODEL=$(MODELS_DIR)/ggml-model-f16.gguf && \
	export YZMA_TEST_ENCODER_MODEL=$(MODELS_DIR)/t5base-encoder-q4_0.gguf && \
	export YZMA_TEST_LORA_MODEL=$(MODELS_DIR)/Gemma2-Base-F32.gguf && \
	export YZMA_TEST_LORA_ADAPTER=$(MODELS_DIR)/Gemma2-Lora-F32-LoRA.gguf && \
	export YZMA_TEST_SPLIT_MODELS="$(MODELS_DIR)/stories15M-q8_0-00001-of-00003.gguf,$(MODELS_DIR)/stories15M-q8_0-00002-of-00003.gguf,$(MODELS_DIR)/stories15M-q8_0-00003-of-00003.gguf" && \
	go test -count=1 ./...

# make wasm-example to build the browser example with TinyGo. Run
# make download-llama.cpp-wasm first to get the WebAssembly build of llama.cpp.
wasm-example: wasm-assets
	tinygo build -target wasm -o $(WASM_DIR)/yzma.wasm ./examples/wasm/chat
	cp "$(shell tinygo env TINYGOROOT)/targets/wasm_exec.js" $(WASM_DIR)/

# make wasm-example-go to build the same example with the standard toolchain.
# The binary is larger, which helps when TinyGo cannot build a dependency.
wasm-example-go: wasm-assets
	GOOS=js GOARCH=wasm go build -o $(WASM_DIR)/yzma.wasm ./examples/wasm/chat
	cp "$(shell go env GOROOT)/lib/wasm/wasm_exec.js" $(WASM_DIR)/

# wasm-assets copies the JavaScript glue and the page into the build directory.
# The two wasm_exec.js files are not the same, so the target that builds the
# program copies the correct one.
wasm-assets:
	mkdir -p $(WASM_DIR)
	cp wasm/yzma-loader.js wasm/worker.js wasm/index.html $(WASM_DIR)/

# make download-llama.cpp-wasm to get the WebAssembly build of llama.cpp.
download-llama.cpp-wasm:
	mkdir -p $(WASM_DIR)
	yzma install -lib $(WASM_DIR) -os wasm $(if $(VERSION),-v $(VERSION))

# make serve-wasm to serve the example with the headers that a build with more
# than one thread needs.
serve-wasm:
	go run ./wasm/serve -dir $(WASM_DIR)

# make vet-wasm to check the WebAssembly code. It does not run in go test,
# because the tests of the repo build for the machine they run on.
vet-wasm:
	GOOS=js GOARCH=wasm go build -o /dev/null ./examples/wasm/chat
	GOOS=js GOARCH=wasm go vet ./pkg/llamawasm ./examples/wasm/chat

# make test-wasm to run the WebAssembly build in Node, with no browser.
test-wasm:
	node wasm/node/run.js --dir $(WASM_DIR) --model $(MODELS_DIR)/SmolLM-135M.Q2_K.gguf --tokens 12

# make test-wasm-mt to run the build with more than one thread in Node. Node has
# SharedArrayBuffer without the headers that a browser needs.
test-wasm-mt:
	node wasm/node/run.js --dir $(WASM_DIR) --model $(MODELS_DIR)/SmolLM-135M.Q2_K.gguf --tokens 12 --mt

# make test-wasm-webgpu to check what happens where there is no WebGPU. Node has
# none, so this must fall back to a build on the CPU and still make text. Only a
# browser can run the WebGPU build itself.
test-wasm-webgpu:
	node wasm/node/run.js --dir $(WASM_DIR) --model $(MODELS_DIR)/SmolLM-135M.Q2_K.gguf --tokens 12 --webgpu

clean-wasm:
	rm -rf $(WASM_DIR)

# make check-ffi to audit the FFI bindings against the headers of the current
# llama.cpp release. The gate runs first: if it fails, do not use what the
# audit reports.
check-ffi:
	cd cmd/yzma-checker && go test ./... && go run .

roadmap:
	@echo "Checked items (have wrapper):"
	@grep -E '^\s*[-*]\s*\[x\]' ROADMAP.md | wc -l
	@echo "Unchecked items (no wrapper):"
	@grep -E '^\s*[-*]\s*\[ \]' ROADMAP.md | wc -l
	@echo "Total checklist items:"
	@grep -E '^\s*[-*]\s*\[(x| )\]' ROADMAP.md | wc -l
