// dawn.js runs run.js with Dawn, which is the WebGPU of Chrome, thus the
// WebGPU build of llama.cpp runs with no browser.
//
// Usage.
//   YZMA_WEBGPU=build/node/node_modules/webgpu/index.js \
//     node wasm/node/dawn.js <arguments of run.js> --webgpu
//
// YZMA_DAWN gives the flags of Dawn, separated by commas, for example
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
  Object.defineProperty(globalThis.navigator, "gpu", {
    value: create(flags),
    configurable: true,
  });

  require("./run.js");
})().catch((err) => {
  console.error(err);
  process.exit(1);
});
