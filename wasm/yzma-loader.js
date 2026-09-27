// yzma-loader.js picks the right WebAssembly build of llama.cpp and
// prepares it for the pkg/llamawasm Go package.
//
// It works in a Web Worker and in a page. Load it before the Go program.
// It sets these globals.
//
//   globalThis.yzmaReady     a promise that resolves to the llama.cpp module
//   globalThis.yzmaModule    the module, after the promise completes
//   globalThis.yzmaThreaded  true if the build is multithreaded
//   globalThis.yzmaThreads   the number of threads the build can use
//   globalThis.yzmaBackend   "webgpu", "cpu-threads", or "cpu"
//   globalThis.yzmaAdapter   the name of the GPU, if there is one
//   globalThis.yzmaGPUReject why the GPU build was dropped, if it was
//
// Set these globals before loading this file to change the result.
//
//   globalThis.yzmaBase      the location of the llama.cpp files, default "."
//   globalThis.yzmaMode      "auto" (the default), "webgpu", or "cpu"
//   globalThis.yzmaPowerPreference "high-performance" or "low-power" to pick
//                            the GPU on a machine with two, default lets the browser pick
//
// There are three builds of llama.cpp. The WebGPU build computes on the GPU and
// needs a browser with WebGPU and JSPI, which is Chrome and Edge 137 or later,
// or Firefox 153 or later with dom.webgpu.enabled and dom.webgpu.workers.enabled
// set in about:config. The two CPU builds work in all browsers. The
// multithreaded build needs an isolated page with the
// Cross-Origin-Opener-Policy and Cross-Origin-Embedder-Policy headers.
//
// In "auto" mode the loader picks the best build the browser can run.
// Firefox is the one exception. Its WebGPU gives llama.cpp correct values but
// is far slower than its CPU, so auto mode picks the CPU there. Mode
// "webgpu" still picks the GPU.
//
// Some drivers expose an adapter that llama.cpp accepts but that computes
// wrong values, so the model answers with random tokens. The loader therefore
// tests the GPU build against the CPU before it hands off the module, and
// picks a CPU build if the test fails. The test needs no model, so it only
// takes a few milliseconds.

(function () {
  const base = globalThis.yzmaBase || ".";
  const mode = globalThis.yzmaMode || "auto";
  const power = globalThis.yzmaPowerPreference || "";

  // A browser only provides SharedArrayBuffer to an isolated page. The
  // multithreaded build needs it.
  const canThread =
    typeof SharedArrayBuffer !== "undefined" && globalThis.crossOriginIsolated === true;

  // webgpuAdapter returns the name of a usable GPU, or an empty string and the
  // reason why. llama.cpp needs f16 shaders and JSPI, so navigator.gpu alone is
  // not enough.
  async function webgpuAdapter() {
    if (!globalThis.navigator || !navigator.gpu) {
      // A worker has its own switch in Firefox, so name it here. This code
      // runs in the worker that holds llama.cpp.
      return [
        "",
        "this worker has no navigator.gpu, in Firefox set dom.webgpu.enabled" +
          " and dom.webgpu.workers.enabled",
      ];
    }

    // The WebGPU build glue uses both parts of JSPI.
    if (
      typeof WebAssembly.Suspending !== "function" ||
      typeof WebAssembly.promising !== "function"
    ) {
      return ["", "no JSPI, which needs Chrome or Edge 137, or Firefox 153, or later"];
    }

    try {
      const adapter = await navigator.gpu.requestAdapter(
        power ? { powerPreference: power } : undefined
      );
      if (!adapter) {
        return ["", "requestAdapter gave no adapter"];
      }
      if (!adapter.features.has("shader-f16")) {
        return [
          "",
          "the adapter has no shader-f16, see the note on NVIDIA in wasm/README.md",
        ];
      }

      // The browser does not always report the GPU name.
      const info = adapter.info || {};
      const name = [info.vendor, info.architecture, info.device, info.description]
        .filter((part) => part)
        .join(" ") || "webgpu";
      return [name, ""];
    } catch (err) {
      return ["", "requestAdapter failed: " + err];
    }
  }

  // firefox reports whether the browser is Firefox. WebGPU in Firefox uses
  // wgpu, which is much slower with llama.cpp than Dawn.
  function firefox() {
    return /firefox/i.test(globalThis.navigator?.userAgent || "");
  }

  // preferAdapter makes llama.cpp request the same GPU the loader tested.
  // llama.cpp sets no power preference of its own.
  function preferAdapter(preference) {
    const gpu = navigator.gpu;
    const request = gpu.requestAdapter.bind(gpu);
    gpu.requestAdapter = (options) =>
      request({ ...options, powerPreference: options?.powerPreference || preference });
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

    // WebGPU in Firefox generates less than one token per second, and the CPU
    // more than a hundred. Auto mode picks the CPU there.
    const skipFirefox = mode !== "webgpu" && mode !== "cpu" && firefox();

    if (skipFirefox) {
      console.warn(
        "yzma: WebGPU in Firefox is much slower than the CPU, using" +
          " the CPU. Set yzmaMode to webgpu to try the GPU."
      );
    }

    if (mode !== "cpu" && !skipFirefox) {
      [adapter, reason] = await webgpuAdapter();
    }

    if (adapter) {
      name = "yzma_wasm_webgpu";
      backend = "webgpu";
      if (power) {
        preferAdapter(power);
      }
    } else if (mode === "webgpu") {
      // The page asked for WebGPU, but the browser cannot provide it. Fall back
      // to the CPU, the same result "auto" would give.
      console.warn("yzma: no WebGPU that llama.cpp can use, using the CPU: " + reason);
      if (canThread) {
        name = "yzma_wasm_mt";
        backend = "cpu-threads";
      }
    } else if (canThread) {
      name = "yzma_wasm_mt";
      backend = "cpu-threads";
    }

    // llama.cpp uses four threads unless the caller changes it, which is slow
    // on a machine with many cores. The Go side reads this value.
    const cores = Math.max(1, Math.min(globalThis.navigator?.hardwareConcurrency || 4, 16));
    let threads = backend === "cpu-threads" ? cores : 1;

    globalThis.yzmaThreaded = backend === "cpu-threads";
    globalThis.yzmaBackend = backend;
    globalThis.yzmaAdapter = adapter;
    globalThis.yzmaThreads = threads;

    let instance = await instantiate(name, threads);

    // A GPU build is useless if the driver computes wrong values, so check
    // before the page downloads a model. The shim compares one small matrix
    // multiply on the GPU against the same one on the CPU.
    if (backend === "webgpu") {
      const trouble = await badBackend(instance);
      if (trouble) {
        console.warn("yzma: " + trouble + ", using the CPU");
        globalThis.yzmaGPUReject = trouble;

        backend = canThread ? "cpu-threads" : "cpu";
        name = canThread ? "yzma_wasm_mt" : "yzma_wasm";
        adapter = "";
        threads = backend === "cpu-threads" ? cores : 1;

        globalThis.yzmaThreaded = backend === "cpu-threads";
        globalThis.yzmaBackend = backend;
        globalThis.yzmaAdapter = adapter;
        globalThis.yzmaThreads = threads;

        instance = await instantiate(name, threads);
      }
    }

    // From here yzmaModule is the instance, which the Go code uses.
    globalThis.yzmaModule = instance;
    return instance;
  })();

  // instantiate loads one build and returns its module.
  async function instantiate(name, threads) {
    // A second build needs its own factory, not the previous one.
    globalThis.yzmaModule = undefined;
    await loadScript(base + "/" + name + ".js");

    // MODULARIZE with EXPORT_NAME=yzmaModule makes yzmaModule a function that
    // returns the instance.
    const factory = globalThis.yzmaModule;
    if (typeof factory !== "function") {
      throw new Error("yzmaModule is not there, check the build of llama.cpp");
    }

    return factory({
      locateFile: (path) => base + "/" + path,
      print: (text) => console.log(text),
      printErr: (text) => console.warn(text),

      // A pool that is too small makes llama.cpp wait for threads that cannot
      // start, because the thread that starts them is busy.
      pthreadPoolSize: threads,
    });
  }

  // badBackend returns why the module's device is not usable, or an empty
  // string if the device matches the CPU. Builds before ABI 8 have no such
  // test, so they pass.
  async function badBackend(instance) {
    if (typeof instance._yzma_backend_check !== "function") {
      return "";
    }
    try {
      // Both calls reach the GPU, so JSPI makes each one return a promise.
      await instance._yzma_backend_init();
      if ((await instance._yzma_backend_check()) === 1) {
        return "llama.cpp computes wrong values on this GPU";
      }
    } catch (e) {
      return "the test of this GPU failed: " + e;
    }
    return "";
  }
})();
