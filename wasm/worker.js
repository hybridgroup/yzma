// worker.js runs llama.cpp and the Go program in a Web Worker.
//
// Each call into llama.cpp is synchronous and one token takes milliseconds,
// so this work cannot run on the main thread without freezing the page.
//
// The page sends these messages to the worker:
//
//   { kind: "load", url: "<model URL>" }
//   { kind: "generate", prompt: "<text>", maxTokens: 128 }
//   { kind: "ask", question: "<text>", maxTokens: 256 }
//
// The worker sends back { kind, text } messages, where kind is one of ready,
// status, progress, loaded, token, tool, result, answer, done, or error.

// Emscripten starts each thread with new Worker(_scriptName), which in a
// classic worker is this file. A thread must only load llama.cpp.
const isThread = globalThis.name === "em-pthread";

self.yzmaBase = ".";

// The page picks the backend with a query on this worker's URL, for example
// new Worker("./worker.js?mode=cpu"), and the GPU with gpu=high-performance
// or gpu=low-power. yzma-loader.js lists the values.
const workerQuery = new URLSearchParams((self.location.search || "").slice(1));
if (workerQuery.get("mode")) {
  self.yzmaMode = workerQuery.get("mode");
}
if (workerQuery.get("gpu")) {
  self.yzmaPowerPreference = workerQuery.get("gpu");
}

// Only the multithreaded build has threads, so a thread does not need to look
// for a GPU.
if (isThread) {
  self.yzmaMode = "cpu";
}

// The Go program to run. The chat page uses the default and the image page
// requests its own.
const program = workerQuery.get("program") || "yzma.wasm";

importScripts("./yzma-loader.js");

// A thread stops here.
if (!isThread) {
  run();
}

// run starts llama.cpp and the Go program for the page.
function run() {
  // An unhandled failure must reach the page. Without this the page just sees
  // nothing happen.
  self.onerror = (event) => {
    self.postMessage({ kind: "error", text: String((event && event.message) || event) });
  };
  self.onunhandledrejection = (event) => {
    self.postMessage({ kind: "error", text: String((event && event.reason) || event) });
  };

  importScripts("./wasm_exec.js");

  // The Go program sends a "ready" message once it has set its functions.
  // Wait for that message, because the backend can take a long time to start.
  let programIsReady;
  const programReady = new Promise((resolve) => {
    programIsReady = resolve;
  });

  const sendToPage = self.postMessage.bind(self);
  self.postMessage = (message) => {
    if (message && (message.kind === "ready" || message.kind === "error")) {
      programIsReady();
    }
    sendToPage(message);
  };

  const started = (async () => {
    // llama.cpp must be ready before the Go program calls Load.
    await self.yzmaReady;

    const go = new Go();
    const result = await WebAssembly.instantiateStreaming(fetch("./" + program), go.importObject);

    // The Go program blocks at the end of main and keeps running, so the page
    // can call into it. Do not wait for this promise.
    go.run(result.instance);

    await programReady;

    if (typeof self.yzmaLoadModel !== "function" && typeof self.yzmaDecideLoad !== "function") {
      throw new Error("the Go program did not set its functions");
    }
  })();

  self.onmessage = async (event) => {
    const message = event.data || {};

    try {
      await started;

      switch (message.kind) {
        case "load":
          self.yzmaLoadModel(message.url, message.projector || "");
          break;
        case "generate":
          self.yzmaGenerate(message.prompt, message.maxTokens || 128);
          break;
        case "ask":
          // The tools page sends a question, not a prompt. The Go program
          // renders the model's chat template itself.
          self.yzmaAsk(message.question, message.maxTokens || 256);
          break;
        case "describe":
          // The image arrives as RGBA from a canvas on the page.
          self.yzmaDescribe(
            message.prompt,
            message.width,
            message.height,
            new Uint8Array(message.rgba),
            message.maxTokens || 128,
          );
          break;
        case "decide-load":
          // The decide page loads a model with its config and a readout.
          self.yzmaDecideLoad(message.url, message.config, message.readout, message.manyMode || "");
          break;
        case "decide":
          self.yzmaDecide(message.state, message.question, message.category || "");
          break;
        case "decide-many":
          self.yzmaDecideMany(message.state, message.questions, message.category || "");
          break;
        default:
          self.postMessage({ kind: "error", text: "unknown message: " + message.kind });
      }
    } catch (err) {
      self.postMessage({ kind: "error", text: String(err) });
    }
  };
}
