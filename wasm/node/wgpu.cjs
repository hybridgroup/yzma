// wgpu.cjs runs run.js in Deno, which uses wgpu, the WebGPU implementation in
// Firefox. This runs the WebGPU build of llama.cpp without a browser.
//
// Usage.
//   deno run -A wasm/node/wgpu.cjs <run.js arguments> --webgpu

require("./run.js");
