# tts

Uses `yzma` with a text to speech model. It speaks the text prompt and writes the audio to a WAV file.

The model uses its first preset voice unless you name one with `-voice`. A model with no preset voices speaks with its own default voice.

Requires `llama.cpp` b12219 or later.

## Running

```shell
go run ./examples/tts/ -model ./models/chatterbox-turbo-Q8_0.gguf -mmproj ./models/mmproj-chatterbox-turbo-F16.gguf -p "Hello from yzma." -o hello.wav
```

The model files are on Hugging Face at [ggml-org/chatterbox-turbo-GGUF](https://huggingface.co/ggml-org/chatterbox-turbo-GGUF).
