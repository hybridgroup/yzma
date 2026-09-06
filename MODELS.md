# Model Usage

The model information is now on the yzma website:

**https://yzma.ai/docs/guides/models/**

That page gives the download command and the run command for each model, in three groups.

- Vision Language Models (VLM)
- Text generation models
- Vision Language Action Models (VLA)

## Quick start

`yzma` uses models in the GGUF format that `llama.cpp` supports. There are many of them on Hugging Face:

https://huggingface.co/models?library=gguf&sort=trending

Download one with the `yzma` command:

```shell
yzma model get -u https://huggingface.co/ggml-org/gemma-3-1b-it-GGUF/resolve/main/gemma-3-1b-it-Q4_K_M.gguf
```

## Related pages

| Page | What it covers |
| --- | --- |
| [Models](https://yzma.ai/docs/concepts/models/) | GGUF, quantization, projector files, context size |
| [Download models](https://yzma.ai/getting-started/download-models/) | How to get a model on your machine |
| [Vision](https://yzma.ai/docs/tutorials/vision/) | How to ask a question about an image |
