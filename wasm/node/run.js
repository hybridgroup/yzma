// run.js runs the WebAssembly build of yzma in Node, without a browser.
//
// It is the test that CI uses: it loads a small model, makes a fixed number of
// tokens with the greedy sampler, and prints them. The greedy sampler always
// takes the most probable token, so the output of a model does not change, and
// a test can compare it.
//
// Usage:
//   node wasm/node/run.js --dir build/wasm --model ~/models/SmolLM-135M.Q2_K.gguf \
//       --prompt "Are you ready to go?" --tokens 12 [--expect "<text>"] [--mt]
//
// --mt takes the build with more than one thread. Node has SharedArrayBuffer
// without the headers that a browser needs, so this is a way to test that build
// outside a browser.

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
const moduleName = mt ? "yzma_wasm_mt.js" : "yzma_wasm.js";

if (!modelFile) {
  console.error("give a model with --model");
  process.exit(2);
}

// The output of the Go program comes here, because a Node process has no
// postMessage.
const output = [];
globalThis.yzmaOnMessage = (message) => {
  if (message.kind === "token") {
    output.push(message.text);
    process.stdout.write(message.text);
    return;
  }
  console.log("[" + message.kind + "] " + message.text);
};

async function main() {
  // Take the build with one thread. Node has no crossOriginIsolated, so the
  // loader would take it anyway, but this says so.
  globalThis.crossOriginIsolated = mt;
  globalThis.yzmaBase = dir;

  const factory = require(path.join(dir, moduleName));
  const llamaModule = await factory({
    locateFile: (file) => path.join(dir, file),
    print: () => {},
    printErr: () => {},
  });

  globalThis.yzmaModule = llamaModule;
  globalThis.yzmaReady = Promise.resolve(llamaModule);
  globalThis.yzmaThreaded = mt;

  // Put the model in the filesystem of the module. A browser gets it over the
  // network instead, with FetchModelFile.
  llamaModule.FS.mkdirTree("/models");
  llamaModule.FS.writeFile("/models/model.gguf", fs.readFileSync(modelFile));

  require(path.join(dir, "wasm_exec.js"));

  const go = new Go();
  const binary = fs.readFileSync(path.join(dir, "yzma.wasm"));
  const result = await WebAssembly.instantiate(binary, go.importObject);

  // The program blocks at the end of main, so do not wait for this.
  go.run(result.instance);

  await settle();

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

// settle gives the Go program the turns it needs to set its functions.
function settle() {
  return new Promise((resolve) => setTimeout(resolve, 50));
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
