# Get the absolute path of the current Makefile
MAKEFILE_PATH := $(realpath $(lastword $(MAKEFILE_LIST)))
MAKEFILE_DIR := $(dir $(MAKEFILE_PATH))
YZMA_LIB ?= $(MAKEFILE_DIR)lib
MODELS_DIR ?= $(HOME)/models

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
	YZMA_LIB=$(YZMA_LIB) go build -o yzma ./cmd/yzma

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

roadmap:
	@echo "Checked items (have wrapper):"
	@grep -E '^\s*[-*]\s*\[x\]' ROADMAP.md | wc -l
	@echo "Unchecked items (no wrapper):"
	@grep -E '^\s*[-*]\s*\[ \]' ROADMAP.md | wc -l
	@echo "Total checklist items:"
	@grep -E '^\s*[-*]\s*\[(x| )\]' ROADMAP.md | wc -l
