// yzma-loader.js takes the correct WebAssembly build of llama.cpp and makes it
// ready for the pkg/llamawasm Go package.
//
// It works in a Web Worker and in a page. Load it before the Go program runs.
// It sets:
//
//   globalThis.yzmaReady     a promise whose value is the llama.cpp module
//   globalThis.yzmaModule    the module, after the promise is done
//   globalThis.yzmaThreaded  true if the build uses more than one thread
//   globalThis.yzmaBackend   "webgpu", "cpu-threads", or "cpu"
//   globalThis.yzmaAdapter   the name of the GPU, if there is one
//
// Set these before this file to change what it does:
//
//   globalThis.yzmaBase      where the files of llama.cpp are. The default is "."
//   globalThis.yzmaMode      "auto" (the default), "webgpu", or "cpu"
//
// There are three builds of llama.cpp. WebGPU puts the computation on the GPU
// and needs a browser that has both WebGPU and JSPI, which today is Chrome and
// Edge 137 and later. The two builds on the CPU work everywhere, and the one
// that uses more than one thread needs a page that is isolated with the
// Cross-Origin-Opener-Policy and Cross-Origin-Embedder-Policy headers.
//
// In "auto" the loader takes the best build that the browser can run, so a
// browser without WebGPU still works.

(function () {
  const base = globalThis.yzmaBase || ".";
  const mode = globalThis.yzmaMode || "auto";

  // A browser gives SharedArrayBuffer only to a page that is isolated. The
  // build with more than one thread cannot run without it.
  const canThread =
    typeof SharedArrayBuffer !== "undefined" && globalThis.crossOriginIsolated === true;

  // webgpuAdapter gives the name of a GPU that the WebGPU build of llama.cpp
  // can use, or an empty string.
  //
  // The test is not only for navigator.gpu. The backend of llama.cpp needs f16
  // shaders, and it reports no device at all without them, and the build needs
  // JSPI to wait for the GPU.
  async function webgpuAdapter() {
    if (!globalThis.navigator || !navigator.gpu) {
      return "";
    }
    if (typeof WebAssembly.Suspending !== "function") {
      // No JSPI, so the WebGPU build cannot run in this browser.
      return "";
    }

    try {
      const adapter = await navigator.gpu.requestAdapter();
      if (!adapter || !adapter.features.has("shader-f16")) {
        return "";
      }

      // The name of the GPU is only there if the browser gives it.
      const info = adapter.info || {};
      return [info.vendor, info.architecture, info.device, info.description]
        .filter((part) => part)
        .join(" ") || "webgpu";
    } catch {
      return "";
    }
  }

  function loadScript(url) {
    // A classic worker takes importScripts. A page takes a script element.
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

    if (mode !== "cpu") {
      adapter = await webgpuAdapter();
    }

    if (adapter) {
      name = "yzma_wasm_webgpu";
      backend = "webgpu";
    } else if (mode === "webgpu") {
      // The page asked for WebGPU and the browser cannot give it. Say so, and
      // then go on with the CPU, which is what "auto" would do.
      console.warn("yzma: this browser has no WebGPU that llama.cpp can use, using the CPU");
      if (canThread) {
        name = "yzma_wasm_mt";
        backend = "cpu-threads";
      }
    } else if (canThread) {
      name = "yzma_wasm_mt";
      backend = "cpu-threads";
    }

    globalThis.yzmaThreaded = backend === "cpu-threads";
    globalThis.yzmaBackend = backend;
    globalThis.yzmaAdapter = adapter;

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
    });

    // From here on yzmaModule is the instance, which is what the Go code uses.
    globalThis.yzmaModule = instance;
    return instance;
  })();
})();
