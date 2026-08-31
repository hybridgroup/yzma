// run.js runs the WebAssembly build of yzma in Node, without a browser.
//
// CI uses this test. It loads a small model, makes a fixed number of tokens
// with the greedy sampler, and prints them. The greedy sampler always takes the
// most probable token, thus the output does not change and a test can compare
// it.
//
// Usage.
//   node wasm/node/run.js --dir build/wasm --model ~/models/SmolLM-135M.Q2_K.gguf \
//       --prompt "Are you ready to go?" --tokens 12 [--expect "<text>"] [--mt]
//       [--webgpu]
//
// --mt selects the build with more than one thread. Node gives
// SharedArrayBuffer without the headers that a browser needs, thus this tests
// that build outside a browser.
//
// --webgpu asks for the WebGPU build. Node has no WebGPU, thus this tests the
// fallback. The loader must select a CPU build and the program must make text.

const fs = require("node:fs");
const path = require("node:path");

function option(name, fallback) {
  const i = process.argv.indexOf("--" + name);
  return i >= 0 && i + 1 < process.argv.length ? process.argv[i + 1] : fallback;
}

const dir = path.resolve(option("dir", "build/wasm"));
const modelFile = option("model", "");
const prompt = option("prompt", "Are you ready to go?");
const maxTokens = parseInt(option("tokens", "12"), 10);
const expect = option("expect", "");
const mt = process.argv.includes("--mt");
const webgpu = process.argv.includes("--webgpu");

// This harness replaces yzma-loader.js. It makes the same choice as the loader
// where there is no WebGPU, thus it falls back to the CPU.
let moduleName = "yzma_wasm.js";
if (mt) {
  moduleName = "yzma_wasm_mt.js";
}
if (webgpu) {
  if (globalThis.navigator && globalThis.navigator.gpu) {
    moduleName = "yzma_wasm_webgpu.js";
  } else {
    console.log("[loader] no WebGPU here, using the build on the CPU");
  }
}

if (!modelFile) {
  console.error("give a model with --model");
  process.exit(2);
}

// The output of the Go program comes here, because Node has no postMessage.
const output = [];

let programIsReady;
const programReady = new Promise((resolve) => {
  programIsReady = resolve;
});

globalThis.yzmaOnMessage = (message) => {
  if (message.kind === "ready" || message.kind === "error") {
    programIsReady();
  }
  if (message.kind === "token") {
    output.push(message.text);
    process.stdout.write(message.text);
    return;
  }
  console.log("[" + message.kind + "] " + message.text);
};

async function main() {
  // Select the build with one thread. Node has no crossOriginIsolated, thus
  // the loader makes the same choice.
  globalThis.crossOriginIsolated = mt;
  globalThis.yzmaBase = dir;

  const factory = require(path.join(dir, moduleName));
  // The same values that the JavaScript glue selects. The pool follows the
  // machine and the Go side reads the thread count.
  const threads = mt ? Math.max(1, Math.min(require("node:os").cpus().length, 16)) : 1;

  const llamaModule = await factory({
    locateFile: (file) => path.join(dir, file),
    print: () => {},
    printErr: () => {},
    pthreadPoolSize: threads,
  });

  globalThis.yzmaModule = llamaModule;
  globalThis.yzmaReady = Promise.resolve(llamaModule);
  globalThis.yzmaThreaded = mt;
  globalThis.yzmaThreads = threads;
  globalThis.yzmaBackend = moduleName.includes("webgpu")
    ? "webgpu"
    : mt
      ? "cpu-threads"
      : "cpu";

  // Put the model in the filesystem of the module. A browser instead gets it
  // from the network with FetchModelFile.
  llamaModule.FS.mkdirTree("/models");
  llamaModule.FS.writeFile("/models/model.gguf", fs.readFileSync(modelFile));

  require(path.join(dir, "wasm_exec.js"));

  const go = new Go();
  const binary = fs.readFileSync(path.join(dir, "yzma.wasm"));
  const result = await WebAssembly.instantiate(binary, go.importObject);

  // The program blocks at the end of main, thus do not wait for this.
  go.run(result.instance);

  // Wait for the ready message. The backend needs time to start, and more
  // time with WebGPU.
  await programReady;

  const done = new Promise((resolve) => {
    const previous = globalThis.yzmaOnMessage;
    globalThis.yzmaOnMessage = (message) => {
      previous(message);
      if (message.kind === "loaded") {
        globalThis.yzmaGenerate(prompt, maxTokens);
      }
      if (message.kind === "done" || message.kind === "error") {
        resolve(message);
      }
    };
  });

  globalThis.yzmaOpenModel("/models/model.gguf");

  const last = await done;
  console.log();

  if (last.kind === "error") {
    process.exit(1);
  }

  const text = output.join("");
  if (expect && text !== expect) {
    console.error("the output is not what the test expects");
    console.error("  expected: " + JSON.stringify(expect));
    console.error("  received: " + JSON.stringify(text));
    process.exit(1);
  }
  if (text.length === 0) {
    console.error("no tokens came out");
    process.exit(1);
  }
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
