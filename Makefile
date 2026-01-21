# Get the absolute path of the current Makefile
MAKEFILE_PATH := $(realpath $(lastword $(MAKEFILE_LIST)))
MAKEFILE_DIR := $(dir $(MAKEFILE_PATH))

ifeq ($(shell uname -s),Darwin)
    MAC_FLAG = -p metal
else
    MAC_FLAG =
endif

LLAMA_VER:=b7787
YZMA_LIB:=$(MAKEFILE_DIR)lib

download-models:
	mkdir -p $(MAKEFILE_DIR)models
	$(MAKEFILE_DIR)yzma model get -y --show-progress=false -o $(MAKEFILE_DIR)models -u https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf
	$(MAKEFILE_DIR)yzma model get -y --show-progress=false -o $(MAKEFILE_DIR)models -u https://huggingface.co/ggml-org/SmolVLM-256M-Instruct-GGUF/resolve/main/SmolVLM-256M-Instruct-Q8_0.gguf
	$(MAKEFILE_DIR)yzma model get -y --show-progress=false -o $(MAKEFILE_DIR)models -u https://huggingface.co/ggml-org/SmolVLM-256M-Instruct-GGUF/resolve/main/mmproj-SmolVLM-256M-Instruct-Q8_0.gguf
	$(MAKEFILE_DIR)yzma model get -y --show-progress=false -o $(MAKEFILE_DIR)models -u https://huggingface.co/ggml-org/models-moved/resolve/main/jina-reranker-v1-tiny-en/ggml-model-f16.gguf
	$(MAKEFILE_DIR)yzma model get -y --show-progress=false -o $(MAKEFILE_DIR)models -u https://huggingface.co/callgg/t5-base-encoder-f32/resolve/main/t5base-encoder-q4_0.gguf
	$(MAKEFILE_DIR)yzma model get -y --show-progress=false -o $(MAKEFILE_DIR)models -u https://huggingface.co/deadprogram/yzma-tests/resolve/main/Gemma2-Base-F32.gguf
	$(MAKEFILE_DIR)yzma model get -y --show-progress=false -o $(MAKEFILE_DIR)models -u https://huggingface.co/deadprogram/yzma-tests/resolve/main/Gemma2-Lora-F32-LoRA.gguf
	$(MAKEFILE_DIR)yzma model get -y --show-progress=false -o $(MAKEFILE_DIR)models -u https://huggingface.co/ggml-org/models-moved/resolve/main/tinyllamas/split/stories15M-q8_0-00001-of-00003.gguf
	$(MAKEFILE_DIR)yzma model get -y --show-progress=false -o $(MAKEFILE_DIR)models -u https://huggingface.co/ggml-org/models-moved/resolve/main/tinyllamas/split/stories15M-q8_0-00002-of-00003.gguf
	$(MAKEFILE_DIR)yzma model get -y --show-progress=false -o $(MAKEFILE_DIR)models -u https://huggingface.co/ggml-org/models-moved/resolve/main/tinyllamas/split/stories15M-q8_0-00003-of-00003.gguf

clean-llama.cpp:
	rm -rf $(MAKEFILE_DIR)lib/lib*
	rm -rf $(MAKEFILE_DIR)lib/llama*
	rm $(MAKEFILE_DIR)lib/LICENSE
	rm $(MAKEFILE_DIR)lib/rpc-server

download-llama.cpp:
	$(MAKEFILE_DIR)yzma install -lib $(YZMA_LIB) -version $(LLAMA_VER) $(MAC_FLAG)

build:
	YZMA_LIB=$(YZMA_LIB) go build -o yzma ./cmd/yzma

test:
	export YZMA_LIB=$(YZMA_LIB) && \
	export YZMA_TEST_MODEL=$(MAKEFILE_DIR)models/SmolLM-135M.Q2_K.gguf && \
	export YZMA_TEST_MMMODEL=$(MAKEFILE_DIR)models/SmolVLM-256M-Instruct-Q8_0.gguf && \
	export YZMA_TEST_MMPROJ=$(MAKEFILE_DIR)models/mmproj-SmolVLM-256M-Instruct-Q8_0.gguf && \
	export YZMA_TEST_QUANTIZE_MODEL=$(MAKEFILE_DIR)models/ggml-model-f16.gguf && \
	export YZMA_TEST_ENCODER_MODEL=$(MAKEFILE_DIR)models/t5base-encoder-q4_0.gguf && \
	export YZMA_TEST_LORA_MODEL=$(MAKEFILE_DIR)models/Gemma2-Base-F32.gguf && \
	export YZMA_TEST_LORA_ADAPTER=$(MAKEFILE_DIR)models/Gemma2-Lora-F32-LoRA.gguf && \
	export YZMA_TEST_SPLIT_MODELS="$(MAKEFILE_DIR)models/stories15M-q8_0-00001-of-00003.gguf,$(MAKEFILE_DIR)models/stories15M-q8_0-00002-of-00003.gguf,$(MAKEFILE_DIR)models/stories15M-q8_0-00003-of-00003.gguf" && \
	go test -count=1 ./...