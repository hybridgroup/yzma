// decide.js runs the decide example of yzma in Node, without a browser.
//
// It loads a System One model, asks every question about the state with
// yzmaDecideMany, then asks each one alone with yzmaDecide, then asks them all
// again, and prints the results as JSON. In the exact mode the first two must
// be the same.
//
// Usage.
//   node wasm/node/decide.js --dir build/wasm --model <gguf> --config <json> \
//       --readout jev|jevk5 [--mode exact|batched] [--state <text>] \
//       [--questions <json list>] [--category <name>] [--expect a,b,c] [--mt]
//
// --expect lists the answer that each question must get.

const fs = require("node:fs");
const path = require("node:path");

function option(name, fallback) {
  const i = process.argv.indexOf("--" + name);
  return i >= 0 && i + 1 < process.argv.length ? process.argv[i + 1] : fallback;
}

const dir = path.resolve(option("dir", "build/wasm"));
const modelFile = option("model", "");
const configFile = option("config", "");
const readout = option("readout", "jev");
const manyMode = option("mode", "exact");
const state = option("state", '{"ticket": "I was charged twice for my subscription and want a refund."}');
const questions = option(
  "questions",
  JSON.stringify([
    { question: "Which team handles this?", options: { billing: "payments, invoices", technical: "bugs" } },
    { question: "The customer wants money back." },
    { type: "score", question: "How urgent is this?", options: ["low", "medium", "high"] },
  ]),
);
const category = option("category", "");
const expect = option("expect", "");
const mt = process.argv.includes("--mt");

if (!modelFile || !configFile) {
  console.error("give a model with --model and its config with --config");
  process.exit(2);
}

let programIsReady;
const programReady = new Promise((resolve) => {
  programIsReady = resolve;
});

let onResult = () => {};
globalThis.yzmaOnMessage = (message) => {
  if (message.kind === "ready" || message.kind === "error") {
    programIsReady();
  }
  if (message.kind === "loaded" || message.kind === "result" || message.kind === "error") {
    onResult(message);
  }
  if (message.kind !== "result") {
    console.log("[" + message.kind + "] " + message.text);
  }
};

// next waits for the next loaded, result or error message.
function next() {
  return new Promise((resolve) => {
    onResult = (message) => {
      onResult = () => {};
      if (message.kind === "error") {
        console.error("error: " + message.text);
        process.exit(1);
      }
      resolve(message);
    };
  });
}

async function main() {
  // This harness replaces yzma-loader.js, with a build on the CPU.
  const moduleName = mt ? "yzma_wasm_mt.js" : "yzma_wasm.js";
  globalThis.crossOriginIsolated = mt;
  globalThis.yzmaBase = dir;

  const factory = require(path.join(dir, moduleName));
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
  globalThis.yzmaBackend = mt ? "cpu-threads" : "cpu";

  llamaModule.FS.mkdirTree("/models");
  llamaModule.FS.writeFile("/models/decide.gguf", fs.readFileSync(modelFile));

  require(path.join(dir, "wasm_exec.js"));

  const go = new Go();
  const binary = fs.readFileSync(path.join(dir, "yzma-decide.wasm"));
  const result = await WebAssembly.instantiate(binary, go.importObject);
  go.run(result.instance);
  await programReady;

  const loaded = next();
  globalThis.yzmaDecideOpen("/models/decide.gguf", fs.readFileSync(configFile, "utf8"), readout, manyMode);
  await loaded;

  const many = next();
  globalThis.yzmaDecideMany(state, questions, category);
  const manyOut = JSON.parse((await many).text);

  const oneOut = [];
  for (const q of JSON.parse(questions)) {
    const one = next();
    globalThis.yzmaDecide(state, JSON.stringify(q), category);
    oneOut.push(JSON.parse((await one).text));
  }

  // A second call reuses the state, and the first one also warms up the module.
  const again = next();
  globalThis.yzmaDecideMany(state, questions, category);
  const againOut = JSON.parse((await again).text);

  console.log(JSON.stringify({ many: manyOut, one: oneOut, again: { ms: againOut.ms } }, null, 2));

  let failed = false;
  manyOut.result.forEach((r, i) => {
    const alone = oneOut[i].result;
    const same = r.probabilities.every((p, k) => p === alone.probabilities[k]);
    if (manyMode === "exact" && !same) {
      console.error("question " + i + ": DecideMany and Decide differ in the exact mode");
      failed = true;
    }
  });
  if (expect) {
    const want = expect.split(",");
    manyOut.result.forEach((r, i) => {
      if (want[i] !== undefined && r.answer !== want[i]) {
        console.error("question " + i + ": answer " + r.answer + ", want " + want[i]);
        failed = true;
      }
    });
  }
  process.exit(failed ? 1 : 0);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
