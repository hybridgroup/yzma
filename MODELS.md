# Model Usage

`yzma` uses models in the GGUF format supported by `llama.cpp`. You can find many models in GGUF format on Hugging Face (over 161k at last count):

https://huggingface.co/models?library=gguf&sort=trending

Here are just a few examples of `yzma` with various well-known language models.

Don't forget to set your `YZMA_LIB` env variable to the directory with your `llama.cpp` library files!

## Vision Language Models (VLM)

### Qwen3-VL-2B-Instruct

https://huggingface.co/bartowski/Qwen_Qwen3-VL-2B-Instruct-GGUF

#### Download the model and projector

Good quality, default size for most use cases, recommended.

```
yzma model get -u https://huggingface.co/bartowski/Qwen_Qwen3-VL-2B-Instruct-GGUF/resolve/main/Qwen_Qwen3-VL-2B-Instruct-Q4_K_M.gguf
```

Or for a smaller model with decent quality, smaller than Q4_K_S with similar performance, recommended.

```
yzma model get -u https://huggingface.co/bartowski/Qwen_Qwen3-VL-2B-Instruct-GGUF/resolve/main/Qwen_Qwen3-VL-2B-Instruct-IQ4_XS.gguf
```

In either case, you will need the projector:

```
yzma model get -u https://huggingface.co/bartowski/Qwen_Qwen3-VL-2B-Instruct-GGUF/resolve/main/mmproj-Qwen_Qwen3-VL-2B-Instruct-f16.gguf
```

#### Running

```
go run ./examples/vlm/ -model ~/models/Qwen_Qwen3-VL-2B-Instruct-Q4_K_M.gguf -mmproj ~/models/mmproj-Qwen_Qwen3-VL-2B-Instruct-f16.gguf -image ./images/domestic_llama.jpg -p "What is in this picture?"
```

### LFM2.5-VL-1.6B

https://huggingface.co/LiquidAI/LFM2.5-VL-1.6B-GGUF


#### Download the model and projector

```
yzma model get -u https://huggingface.co/LiquidAI/LFM2.5-VL-1.6B-GGUF/resolve/main/LFM2.5-VL-1.6B-Q8_0.gguf
```

Smaller model

```
yzma model get -u https://huggingface.co/LiquidAI/LFM2.5-VL-1.6B-GGUF/resolve/main/LFM2.5-VL-1.6B-Q4_0.gguf
```

In either case, you will need the projector:

```
yzma model get -u https://huggingface.co/LiquidAI/LFM2.5-VL-1.6B-GGUF/resolve/main/mmproj-LFM2.5-VL-1.6b-Q8_0.gguf
```

#### Running

```
go run ./examples/vlm/ -model ~/models/LFM2.5-VL-1.6B-Q4_0.gguf -mmproj ~/models/mmproj-LFM2.5-VL-1.6b-Q8_0.gguf -image ./images/domestic_llama.jpg -p "What is in this picture?"
```

### Qwen2.5-VL-3B-Instruct

https://huggingface.co/ggml-org/Qwen2.5-VL-3B-Instruct-GGUF

#### Download the model and projector

```
yzma model get -u https://huggingface.co/ggml-org/Qwen2.5-VL-3B-Instruct-GGUF/resolve/main/Qwen2.5-VL-3B-Instruct-Q8_0.gguf
yzma model get -u https://huggingface.co/ggml-org/Qwen2.5-VL-3B-Instruct-GGUF/resolve/main/mmproj-Qwen2.5-VL-3B-Instruct-Q8_0.gguf
```

#### Running

```
go run ./examples/vlm/ -model ~/models/Qwen2.5-VL-3B-Instruct-Q8_0.gguf -mmproj ~/models/mmproj-Qwen2.5-VL-3B-Instruct-Q8_0.gguf -image ./images/domestic_llama.jpg -p "What is in this picture?"
```

### moondream2-20250414-GGUF

https://huggingface.co/ggml-org/moondream2-20250414-GGUF

#### Download the model and projector

```
yzma model get -u https://huggingface.co/ggml-org/moondream2-20250414-GGUF/resolve/main/moondream2-text-model-f16_ct-vicuna.gguf
yzma model get -u https://huggingface.co/ggml-org/moondream2-20250414-GGUF/resolve/main/moondream2-mmproj-f16-20250414.gguf
```

#### Running

```
go run ./examples/vlm/ -model ~/models/moondream2-text-model-f16_ct-vicuna.gguf -mmproj ~/models/moondream2-mmproj-f16-20250414.gguf -image ./images/domestic_llama.jpg -p "What is in this picture?"
```

## Text generation models

### Qwen3-4B-GGUF

https://huggingface.co/Qwen/Qwen3-4B-GGUF

#### Download the model

```
yzma model get -u https://huggingface.co/Qwen/Qwen3-4B-GGUF/blob/main/Qwen3-4B-Q4_K_M.gguf
```

#### Running

```
go run ./examples/chat/ -model ~/models/Qwen3-4B-Q4_K_M.gguf -temp=0.6 -n=512
```

### Qwen3-0.6B-GGUF

https://huggingface.co/bartowski/Qwen_Qwen3-0.6B-GGUF

#### Download the model

```
yzma model get -u https://huggingface.co/bartowski/Qwen_Qwen3-0.6B-GGUF/resolve/main/Qwen_Qwen3-0.6B-Q4_K_M.gguf
```

#### Running

```
go run ./examples/chat/ -model ~/models/Qwen_Qwen3-0.6B-Q4_K_M.gguf -temp=0.6 -n=512
```

### qwen2.5-0.5b-instruct

https://huggingface.co/Qwen/Qwen2.5-0.5B-Instruct-GGUF

#### Download the model

```
yzma model get -u https://huggingface.co/Qwen/Qwen2.5-0.5B-Instruct-GGUF/resolve/main/qwen2.5-0.5b-instruct-q4_k_m.gguf
```

#### Running

```
go run ./examples/chat/ -model ~/models/qwen2.5-0.5b-instruct-q4_k_m.gguf -temp=0.6 -n=512
```

### tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf

https://huggingface.co/TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF

#### Download the model

```
yzma model get -u https://huggingface.co/TheBloke/TinyLlama-1.1B-Chat-v1.0-GGUF/blob/main/tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf
```

#### Running

```
go run ./examples/chat/ -model ~/models/tinyllama-1.1b-chat-v1.0.Q4_K_M.gguf -c 2048 -temp 0.7 -n 512 -sys "You are a helpful robot companion."
```

### gemma-3-1b-it

https://huggingface.co/ggml-org/gemma-3-1b-it-GGUF

#### Download the model

```
yzma model get -u https://huggingface.co/ggml-org/gemma-3-1b-it-GGUF/resolve/main/gemma-3-1b-it-Q4_K_M.gguf
```

#### Running

```
go run ./examples/chat/ -model ~/models/gemma-3-1b-it-Q4_K_M.gguf
```

### SmolLM2-135M-Instruct

https://huggingface.co/bartowski/SmolLM2-135M-Instruct-GGUF

#### Download the model

```
yzma model get -u https://huggingface.co/bartowski/SmolLM2-135M-Instruct-GGUF/resolve/main/SmolLM2-135M-Instruct-Q4_K_M.gguf
```

#### Running

```
go run ./examples/chat/ -model ~/models/SmolLM2-135M-Instruct-Q4_K_M.gguf -c 2048 -temp 0.8 -n 48 -sys "You are a helpful robot companion."
```

## Visual Language Action (VLA) models


### Pelican1.0-VL-3B

https://huggingface.co/mradermacher/Pelican1.0-VL-3B-i1-GGUF

#### Download the model and projector

Fast, recommended:

```
yzma model get -u https://huggingface.co/mradermacher/Pelican1.0-VL-3B-i1-GGUF/resolve/main/Pelican1.0-VL-3B.i1-Q4_K_M.gguf
```

or alternate with lower quality:

```
yzma model get -u https://huggingface.co/mradermacher/Pelican1.0-VL-3B-i1-GGUF/resolve/main/Pelican1.0-VL-3B.i1-IQ3_XXS.gguf
```

In either case, download the projector file:

```
yzma model get -u https://huggingface.co/mradermacher/Pelican1.0-VL-3B-GGUF/resolve/main/Pelican1.0-VL-3B.mmproj-Q8_0.gguf
```

#### Running

```
$ go run ./examples/vlm/ -model ~/models/Pelican1.0-VL-3B.i1-Q4_K_M.gguf --mmproj ~/models/Pelican1.0-VL-3B.mmproj-Q8_0.gguf -p "What is in this picture? Provide a description, bounding box, and estimated distance for the llama in json format." -sys "You are a helpful robotic drone camera currently in flight." -image ./images/domestic_llama.jpg
```

```json
{
  "description": "The image shows a fluffy white llama standing on a green grassy area with a dirt path nearby. The llama has a curly coat and appears to be in a fenced-in area with trees and some buildings in the background.",
  "bounding_box_2d": [40, 155, 635, 811],
  "estimated_distance": "The llama is approximately 1 meter away from the camera."
}
```

### InternVLA-M1

https://huggingface.co/mradermacher/InternVLA-M1-GGUF

#### Download the model and projector

```
yzma model get -u https://huggingface.co/mradermacher/InternVLA-M1-GGUF/resolve/main/InternVLA-M1.Q8_0.gguf
yzma model get -u https://huggingface.co/mradermacher/InternVLA-M1-GGUF/resolve/main/InternVLA-M1.mmproj-Q8_0.gguf
```

#### Running

```
go run ./examples/vlm/ -model ~/models/InternVLA-M1.Q8_0.gguf --mmproj ~/models/InternVLA-M1.mmproj-Q8_0.gguf -p "What is in this picture? Provide a description, bounding box, and estimated distance for the llama in json format." -sys "You are a helpful robotic drone camera currently in flight." -image ./images/domestic_llama.jpg
```

```json
{"label": "llama", "bbox_2d": [43, 352, 647, 822], "distance": 10.0}
```

### SpaceQwen2.5-VL-3B-Instruct

https://huggingface.co/mradermacher/SpaceQwen2.5-VL-3B-Instruct-GGUF

#### Download the model and projector

```
yzma model get -u https://huggingface.co/mradermacher/SpaceQwen2.5-VL-3B-Instruct-i1-GGUF/resolve/main/SpaceQwen2.5-VL-3B-Instruct.i1-Q4_K_M.gguf
yzma model get -u https://huggingface.co/remyxai/SpaceQwen2.5-VL-3B-Instruct/resolve/main/spaceqwen2.5-vl-3b-instruct-vision.gguf
```

#### Running

```
$ go run ./examples/vlm/ -model ~/models/SpaceQwen2.5-VL-3B-Instruct.i1-Q4_K_M.gguf --mmproj ~/models/spaceqwen2.5-vl-3b-instruct-vision.gguf -p "What is in this picture? Provide a description, bounding box, and estimated distance for the llama in json format." -sys "You are a helpful robotic drone camera currently in flight." -image ./images/domestic_llama.jpg  
```

```json
{
  "bbox_2d": [40, 20, 67, 35],
  "label": "llama",
  "estimated_distance": "1.5 meters"
}
```
