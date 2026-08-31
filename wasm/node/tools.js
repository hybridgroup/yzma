// tools.js runs the tool calling build of yzma in Node, without a browser.
//
// CI uses this test. It loads a model, asks a question, and checks that the
// round trip works: the chat template renders, tokens come out, and the program
// finishes. A model that is trained for tool calls also calls one, which
// --expect-tool checks.
//
// Usage.
//   node wasm/node/tools.js --dir build/wasm --model ~/models/model.gguf \
//       --question "What is the weather in Paris?" --tokens 96 \
//       [--expect-tool get_weather] [--mt] [--webgpu]
//
// --mt selects the build with more than one thread, and --webgpu asks for the
// WebGPU build. See run.js for what those mean in Node.

const fs = require("node:fs");
const path = require("node:path");

function option(name, fallback) {
  const i = process.argv.indexOf("--" + name);
  return i >= 0 && i + 1 < process.argv.length ? process.argv[i + 1] : fallback;
}

const dir = path.resolve(option("dir", "build/wasm"));
const modelFile = option("model", "");
const question = option("question", "What is the weather in Paris?");
const maxTokens = parseInt(option("tokens", "96"), 10);
const expectTool = option("expect-tool", "");
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
const toolCalls = [];

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
  if (message.kind === "tool") {
    toolCalls.push(message.text);
  }
  console.log("\n[" + message.kind + "] " + message.text);
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
  const binary = fs.readFileSync(path.join(dir, "yzma-tools.wasm"));
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
        globalThis.yzmaAsk(question, maxTokens);
      }
      if (message.kind === "done" || message.kind === "error") {
        resolve(message);
      }
    };
  });

  // The name of the file of the model selects the format of the tool calls.
  globalThis.yzmaOpenModel("/models/model.gguf", path.basename(modelFile));

  const last = await done;
  console.log();

  if (last.kind === "error") {
    process.exit(1);
  }

  if (output.join("").length === 0) {
    console.error("no tokens came out");
    process.exit(1);
  }

  // Only a model that is trained for tool calls makes one, thus this is a
  // choice of the caller and not the default.
  if (expectTool && !toolCalls.some((call) => call.startsWith(expectTool))) {
    console.error("the model did not call " + expectTool);
    console.error("  calls: " + JSON.stringify(toolCalls));
    process.exit(1);
  }
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
