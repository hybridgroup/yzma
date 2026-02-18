// yzma lets you write Go applications that directly integrate llama.cpp (https://github.com/ggml-org/llama.cpp)
// for fully local inference using hardware acceleration.
//
//   - Run the latest Vision Language Models (VLM) and Large/Small/Tiny Language Models (LLM) on Linux, macOS, or Windows.
//   - Use any available hardware acceleration such as CUDA (https://en.wikipedia.org/wiki/CUDA), Metal (https://en.wikipedia.org/wiki/Metal_(API)),
//     or Vulkan (https://en.wikipedia.org/wiki/Vulkan) for maximum performance.
//   - yzma uses the purego (https://github.com/ebitengine/purego) and  ffi (https://github.com/JupiterRider/ffi) packages so CGo is not needed.
//   - Works with the newest llama.cpp releases so you can use the latest features and model support.
package yzma
