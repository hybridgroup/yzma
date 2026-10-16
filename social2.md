# Jeyzma announcement posts

## Bluesky

1. Announcing Jeyzma: Jev in your browser written in @golang with yzma and WebAssembly. Ask a System One model a typed question and get a decision. No server, no API key, no data leaves your machine. WebGPU acceleration too!

https://jeyzma.com

#jev #golang #webassembly #tinygo #webgpu

2. Need to make a decision? Try Jeyzma! Jev in your browser written in Go with WebGPU support. Choice, score, or yes and no, each answer with its probability. https://jeyzma.com

3. Jeyzma gives a decision without generating any text. One forward pass, and the answer is a probability for each option. 100% local in your browser with WebAssembly, built with yzma and TinyGo. https://jeyzma.com

4. Which team handles this ticket? Is the customer angry? How urgent is it? Jeyzma answers all three in about 2 seconds with WebGPU, right in your browser. Try it at https://jeyzma.com

5. Jev-Style and JevK5 now run in your browser. Jeyzma is written in Go with yzma, compiled with TinyGo, and uses llama.cpp in WebAssembly. Source code at https://github.com/hybridgroup/jeyzma. Try it at https://jeyzma.com

## LinkedIn

### 1

Announcing Jeyzma: Jev in your browser written in Go with yzma.

Jeyzma runs a System One model completely in your browser. You give it a state and a list of typed questions, and it gives back a probability for each option. There is no server, no API key, and no data leaves your machine.

Try it at https://jeyzma.com

### 2

Need to make a decision? Try Jeyzma! Jev in your browser written in Go with WebGPU support.

A System One model does not write an answer. It reads one forward pass and gives a probability for each option. So you can route a support ticket, flag a comment, or rate urgency and see how sure the model is.

Try it at https://jeyzma.com

### 3

Most AI models answer with text. A System One model answers with numbers.

Jeyzma asks three kinds of questions: a choice between named options, a score on a scale, and a yes or no. Each answer comes with a probability for every option, from one forward pass and with no generated text.

It runs 100% local in your browser, with WebGPU when available and multithreaded CPU otherwise.

https://jeyzma.com

### 4

How it works under the hood.

Jeyzma is written in Go and compiled to WebAssembly with TinyGo. It uses the yzma exp/decide package to drive llama.cpp, which runs as a second WebAssembly module. Everything runs in a Web Worker, so the page stays responsive.

With WebGPU, the Jev-Style 0.8B model answers three questions about a support ticket in about 2 seconds.

Source code at https://github.com/hybridgroup/jeyzma
Try it at https://jeyzma.com

### 5

Private decisions, made on your own machine.

Jeyzma runs Jev-Style and JevK5 System One models directly in your browser. Your data never leaves your device, and there is nothing to install. Load a model, describe the situation, ask your questions, and see the probabilities.

Go with your own intelligence.

https://jeyzma.com
