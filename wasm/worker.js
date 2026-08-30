// worker.js runs llama.cpp and the Go program in a Web Worker.
//
// Every call into llama.cpp is synchronous and one token takes milliseconds, so
// this work cannot go on the main thread of the page: it would stop the page.
// The worker sends each piece of text to the page as it comes.
//
// The page sends these messages to the worker:
//
//   { kind: "load", url: "<model URL>" }
//   { kind: "generate", prompt: "<text>", maxTokens: 128 }
//
// The worker sends back { kind, text } messages, where kind is one of ready,
// status, progress, loaded, token, done, or error.

self.yzmaBase = ".";

// The page can choose the backend with a query on the URL of this worker, such
// as new Worker("./worker.js?mode=cpu"). The values are the ones that
// yzma-loader.js takes: auto, webgpu, or cpu.
const workerQuery = new URLSearchParams((self.location.search || "").slice(1));
if (workerQuery.get("mode")) {
  self.yzmaMode = workerQuery.get("mode");
}

// A failure with nobody to catch it must reach the page. Without this the page
// only sees that nothing more happens.
self.onerror = (event) => {
  self.postMessage({ kind: "error", text: String((event && event.message) || event) });
};
self.onunhandledrejection = (event) => {
  self.postMessage({ kind: "error", text: String((event && event.reason) || event) });
};

importScripts("./yzma-loader.js");
importScripts("./wasm_exec.js");

// The Go program says when it is ready by sending its first message, which is
// the one of kind "ready". Waiting for that is the only safe way to know that
// it has set its functions: starting the backend takes a moment with the CPU
// and much longer with WebGPU, where it has to find an adapter and make the
// shaders.
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
  const result = await WebAssembly.instantiateStreaming(fetch("./yzma.wasm"), go.importObject);

  // The Go program blocks at the end of main, so it keeps running and the page
  // can call into it. Do not wait for this promise.
  go.run(result.instance);

  await programReady;

  if (typeof self.yzmaLoadModel !== "function") {
    throw new Error("the Go program did not set its functions");
  }
})();

self.onmessage = async (event) => {
  const message = event.data || {};

  try {
    await started;

    switch (message.kind) {
      case "load":
        self.yzmaLoadModel(message.url);
        break;
      case "generate":
        self.yzmaGenerate(message.prompt, message.maxTokens || 128);
        break;
      default:
        self.postMessage({ kind: "error", text: "unknown message: " + message.kind });
    }
  } catch (err) {
    self.postMessage({ kind: "error", text: String(err) });
  }
};
