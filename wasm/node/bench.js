// bench.js measures the WebAssembly build of yzma in Node, without a browser.
//
// It makes text more than one time and prints the result in the shape that
// go test -bench prints, thus cmd/yzma-bench can read it.
//
// Usage.
//   node wasm/node/bench.js --dir build/wasm --model ~/models/SmolLM-135M.Q2_K.gguf \
//       --prompt "Are you ready to go?" --tokens 64 --count 5 [--mt] [--webgpu]
//
// --mt selects the build with more than one thread. --webgpu asks for the
// WebGPU build, which Node does not have, thus it falls back to the CPU.
//
// The tokens a second come from the Go side, which times the loop of decode
// and sample. A browser shows the same value. The output also gives the
// llama.cpp build, which comes from yzma-install.json of the directory.

const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");

function option(name, fallback) {
  const i = process.argv.indexOf("--" + name);
  return i >= 0 && i + 1 < process.argv.length ? process.argv[i + 1] : fallback;
}

const dirOption = option("dir", "build/wasm");
const dir = path.resolve(dirOption);
const modelFile = option("model", "");
const prompt = option("prompt", "Are you ready to go?");
const maxTokens = parseInt(option("tokens", "64"), 10);
const count = parseInt(option("count", "5"), 10);
const mt = process.argv.includes("--mt");
const webgpu = process.argv.includes("--webgpu");

let moduleName = "yzma_wasm.js";
if (mt) {
  moduleName = "yzma_wasm_mt.js";
}
if (webgpu) {
  if (globalThis.navigator && globalThis.navigator.gpu) {
    moduleName = "yzma_wasm_webgpu.js";
  } else {
    console.error("[loader] no WebGPU here, using the build on the CPU");
  }
}

if (!modelFile) {
  console.error("give a model with --model");
  process.exit(2);
}

// installTag gives the llama.cpp build of the WebAssembly install. A nightly
// build has no upstream_tag, because its own tag names the assets.
function installTag(dir) {
  try {
    const record = JSON.parse(
      fs.readFileSync(path.join(dir, "yzma-install.json"), "utf8"),
    );
    return record.upstream_tag || record.tag || "";
  } catch {
    return "";
  }
}

let programIsReady;
const programReady = new Promise((resolve) => {
  programIsReady = resolve;
});

// The messages of the Go program come here. Only the last one of a generation
// matters, because it carries the count of tokens and the rate.
let onDone = () => {};
globalThis.yzmaOnMessage = (message) => {
  if (message.kind === "ready" || message.kind === "error") {
    programIsReady();
  }
  if (message.kind === "done" || message.kind === "error") {
    onDone(message);
  }
};

function generateOnce() {
  return new Promise((resolve, reject) => {
    onDone = (message) => {
      onDone = () => {};
      if (message.kind === "error") {
        reject(new Error(message.text));
        return;
      }
      resolve(message.text);
    };
    globalThis.yzmaGenerate(prompt, maxTokens);
  });
}

// parseDone reads "24 tokens, 63.30 tokens/s" from the Go side.
function parseDone(text) {
  const tokens = /(\d+) tokens/.exec(text);
  const rate = /([\d.]+) tokens\/s/.exec(text);
  if (!tokens || !rate) {
    throw new Error("the message of the Go side has no rate: " + text);
  }

  return { tokens: parseInt(tokens[1], 10), rate: parseFloat(rate[1]) };
}

async function main() {
  globalThis.crossOriginIsolated = mt;
  globalThis.yzmaBase = dir;

  const factory = require(path.join(dir, moduleName));
  const threads = mt ? Math.max(1, Math.min(os.cpus().length, 16)) : 1;

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

  llamaModule.FS.mkdirTree("/models");
  llamaModule.FS.writeFile("/models/model.gguf", fs.readFileSync(modelFile));

  require(path.join(dir, "wasm_exec.js"));

  const go = new Go();
  const binary = fs.readFileSync(path.join(dir, "yzma.wasm"));
  const result = await WebAssembly.instantiate(binary, go.importObject);

  go.run(result.instance);
  await programReady;

  const loaded = new Promise((resolve) => {
    const previous = globalThis.yzmaOnMessage;
    globalThis.yzmaOnMessage = (message) => {
      previous(message);
      if (message.kind === "loaded") {
        resolve();
      }
    };
  });

  globalThis.yzmaOpenModel("/models/model.gguf");
  await loaded;

  // The first generation warms the build up and is not in the result.
  await generateOnce();

  const runs = [];
  const started = Date.now();
  for (let i = 0; i < count; i++) {
    runs.push(parseDone(await generateOnce()));
  }
  const elapsed = (Date.now() - started) / 1000;

  const command =
    "$ node wasm/node/bench.js --dir " +
    dirOption +
    " --model " +
    path.basename(modelFile) +
    " --tokens " +
    maxTokens +
    (mt ? " --mt" : "") +
    (webgpu ? " --webgpu" : "");

  const name = "BenchmarkInference-" + threads;
  const lines = [
    command,
    "goos: js",
    "goarch: wasm",
    "pkg: github.com/hybridgroup/yzma/examples/wasm/chat",
    "cpu: " + os.cpus()[0].model,
    "backend: " + globalThis.yzmaBackend + ", " + threads + " thread(s)",
  ];
  const tag = installTag(dir);
  if (tag) {
    lines.push("llama.cpp: " + tag);
  }
  for (const run of runs) {
    const nsPerOp = Math.round((run.tokens / run.rate) * 1e9);
    lines.push(
      [name, "1", nsPerOp + " ns/op", run.rate.toFixed(1) + " tokens/s"].join("\t"),
    );
  }
  lines.push("PASS");
  lines.push("ok\tgithub.com/hybridgroup/yzma/examples/wasm/chat\t" + elapsed.toFixed(3) + "s");

  console.log(lines.join("\n"));
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
