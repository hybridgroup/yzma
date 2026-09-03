// yzma-loader.js selects the correct WebAssembly build of llama.cpp and
// prepares it for the pkg/llamawasm Go package.
//
// It operates in a Web Worker and in a page. Load it before the Go program.
// It sets these globals.
//
//   globalThis.yzmaReady     a promise that gives the llama.cpp module
//   globalThis.yzmaModule    the module, after the promise completes
//   globalThis.yzmaThreaded  true if the build uses more than one thread
//   globalThis.yzmaThreads   the number of threads that the build can use
//   globalThis.yzmaBackend   "webgpu", "cpu-threads", or "cpu"
//   globalThis.yzmaAdapter   the name of the GPU, if there is one
//
// Set these globals before this file to change the result.
//
//   globalThis.yzmaBase      the location of the llama.cpp files, default "."
//   globalThis.yzmaMode      "auto" (the default), "webgpu", or "cpu"
//
// There are three builds of llama.cpp. The WebGPU build computes on the GPU and
// needs a browser with WebGPU and JSPI, which is Chrome and Edge 137 or later,
// or Firefox 153 or later with dom.webgpu.enabled and dom.webgpu.workers.enabled
// set in about:config. The two CPU builds operate in all browsers. The build
// with more than one thread needs an isolated page with the
// Cross-Origin-Opener-Policy and Cross-Origin-Embedder-Policy headers.
//
// In "auto" mode the loader selects the best build that the browser can run.
// Firefox is the one exception. Its WebGPU gives wrong values to llama.cpp,
// thus auto mode takes the CPU there. Mode "webgpu" still selects the GPU.

(function () {
  const base = globalThis.yzmaBase || ".";
  const mode = globalThis.yzmaMode || "auto";

  // A browser gives SharedArrayBuffer only to an isolated page. The build with
  // more than one thread needs it.
  const canThread =
    typeof SharedArrayBuffer !== "undefined" && globalThis.crossOriginIsolated === true;

  // webgpuAdapter gives the name of a usable GPU, or an empty string, and the
  // reason for an empty result. llama.cpp needs f16 shaders and JSPI, so
  // navigator.gpu alone is not sufficient.
  async function webgpuAdapter() {
    if (!globalThis.navigator || !navigator.gpu) {
      // A worker has its own switch in Firefox, thus name it here. This code
      // runs in the worker that holds llama.cpp.
      return [
        "",
        "this worker has no navigator.gpu, in Firefox set dom.webgpu.enabled" +
          " and dom.webgpu.workers.enabled",
      ];
    }

    // The glue of the WebGPU build uses both parts of JSPI.
    if (
      typeof WebAssembly.Suspending !== "function" ||
      typeof WebAssembly.promising !== "function"
    ) {
      return ["", "no JSPI, which needs Chrome or Edge 137, or Firefox 153, or later"];
    }

    try {
      const adapter = await navigator.gpu.requestAdapter();
      if (!adapter) {
        return ["", "requestAdapter gave no adapter"];
      }
      if (!adapter.features.has("shader-f16")) {
        return [
          "",
          "the adapter has no shader-f16, see the note on NVIDIA in wasm/README.md",
        ];
      }

      // The browser does not always give the name of the GPU.
      const info = adapter.info || {};
      const name = [info.vendor, info.architecture, info.device, info.description]
        .filter((part) => part)
        .join(" ") || "webgpu";
      return [name, ""];
    } catch (err) {
      return ["", "requestAdapter failed: " + err];
    }
  }

  // firefox says if the browser is Firefox. The WebGPU of Firefox uses wgpu,
  // which does not give llama.cpp the same results as Dawn.
  function firefox() {
    return /firefox/i.test(globalThis.navigator?.userAgent || "");
  }

  function loadScript(url) {
    // A classic worker uses importScripts. A page uses a script element.
    if (typeof importScripts === "function") {
      importScripts(url);
      return Promise.resolve();
    }
    return new Promise((resolve, reject) => {
      const element = document.createElement("script");
      element.src = url;
      element.onload = () => resolve();
      element.onerror = () => reject(new Error("cannot load " + url));
      document.head.appendChild(element);
    });
  }

  globalThis.yzmaReady = (async () => {
    let name = "yzma_wasm";
    let backend = "cpu";
    let adapter = "";
    let reason = "";

    // llama.cpp computes wrong values with WebGPU in Firefox, and a model
    // stops before the first token. Auto mode takes the CPU there.
    const skipFirefox = mode !== "webgpu" && mode !== "cpu" && firefox();

    if (skipFirefox) {
      console.warn(
        "yzma: llama.cpp computes wrong values with WebGPU in Firefox, using" +
          " the CPU. Set yzmaMode to webgpu to try the GPU."
      );
    }

    if (mode !== "cpu" && !skipFirefox) {
      [adapter, reason] = await webgpuAdapter();
    }

    if (adapter) {
      name = "yzma_wasm_webgpu";
      backend = "webgpu";
    } else if (mode === "webgpu") {
      // The page asked for WebGPU, but the browser cannot give it. Continue
      // with the CPU, which is the result that "auto" gives.
      console.warn("yzma: no WebGPU that llama.cpp can use, using the CPU: " + reason);
      if (canThread) {
        name = "yzma_wasm_mt";
        backend = "cpu-threads";
      }
    } else if (canThread) {
      name = "yzma_wasm_mt";
      backend = "cpu-threads";
    }

    // llama.cpp uses four threads unless a caller changes it, which is slow on
    // a machine with many cores. The Go side reads this value.
    const cores = Math.max(1, Math.min(globalThis.navigator?.hardwareConcurrency || 4, 16));
    const threads = backend === "cpu-threads" ? cores : 1;

    globalThis.yzmaThreaded = backend === "cpu-threads";
    globalThis.yzmaBackend = backend;
    globalThis.yzmaAdapter = adapter;
    globalThis.yzmaThreads = threads;

    await loadScript(base + "/" + name + ".js");

    // MODULARIZE with EXPORT_NAME=yzmaModule makes yzmaModule a function that
    // gives the instance.
    const factory = globalThis.yzmaModule;
    if (typeof factory !== "function") {
      throw new Error("yzmaModule is not there, check the build of llama.cpp");
    }

    const instance = await factory({
      locateFile: (path) => base + "/" + path,
      print: (text) => console.log(text),
      printErr: (text) => console.warn(text),

      // A pool that is too small makes llama.cpp wait for threads that cannot
      // start, because the thread that starts them is busy.
      pthreadPoolSize: threads,
    });

    // From here yzmaModule is the instance, which the Go code uses.
    globalThis.yzmaModule = instance;
    return instance;
  })();
})();
