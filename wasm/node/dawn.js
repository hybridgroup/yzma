// dawn.js runs run.js with Dawn, the WebGPU implementation in Chrome. This runs
// the WebGPU build of llama.cpp without a browser.
//
// Usage.
//   YZMA_WEBGPU=build/node/node_modules/webgpu/index.js \
//     node wasm/node/dawn.js <run.js arguments> --webgpu
//
// YZMA_DAWN sets comma separated Dawn flags, for example
// backend=vulkan,adapter=NVIDIA. It needs Node 25 or later for JSPI.

const path = require("node:path");
const { pathToFileURL } = require("node:url");

(async () => {
  const file = process.env.YZMA_WEBGPU || "webgpu";
  const { create, globals } = await import(
    path.isAbsolute(file) ? pathToFileURL(file).href : file
  );
  Object.assign(globalThis, globals);

  const flags = (process.env.YZMA_DAWN || "").split(",").filter((flag) => flag);
  const gpu = create(flags);

  // Dawn only returns an OpenGL adapter in compatibility mode, which Chrome requests.
  if (flags.some((flag) => /^backend=opengl(es)?$/.test(flag))) {
    const request = gpu.requestAdapter.bind(gpu);
    gpu.requestAdapter = (options) => request({ ...options, featureLevel: "compatibility" });
  }

  Object.defineProperty(globalThis.navigator, "gpu", { value: gpu, configurable: true });

  require("./run.js");
})().catch((err) => {
  console.error(err);
  process.exit(1);
});
