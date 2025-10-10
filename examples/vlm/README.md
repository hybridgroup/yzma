# vlm

Uses `yzma` with a Visual Language Model (VLM) multimodal model. It passes in both an image and a text prompt, and then receives back the response from the model.

## Running

```shell
go run ./examples/vlm/ -model ./models/Qwen2.5-VL-3B-Instruct-Q8_0.gguf -proj ./models/mmproj-Qwen2.5-VL-3B-Instruct-Q8_0.gguf -image ./images/domestic_llama.jpg -prompt "What is in this picture?" 2>/dev/null
```
