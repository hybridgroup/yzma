// wgpu.cjs runs run.js in Deno, which has wgpu, the WebGPU of Firefox. Thus
// the WebGPU build of llama.cpp runs with no browser.
//
// Usage.
//   deno run -A wasm/node/wgpu.cjs <arguments of run.js> --webgpu

require("./run.js");
