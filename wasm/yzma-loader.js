// yzma-loader.js takes the correct WebAssembly build of llama.cpp and makes it
// ready for the pkg/llamawasm Go package.
//
// It works in a Web Worker and in a page. Load it before the Go program runs.
// It sets:
//
//   globalThis.yzmaReady     a promise whose value is the llama.cpp module
//   globalThis.yzmaModule    the module, after the promise is done
//   globalThis.yzmaThreaded  true if the build uses more than one thread
//
// Set globalThis.yzmaBase before this file if the files of llama.cpp are not
// beside it.

(function () {
  const base = globalThis.yzmaBase || ".";

  // A browser gives SharedArrayBuffer only to a page that is isolated with the
  // Cross-Origin-Opener-Policy and Cross-Origin-Embedder-Policy headers. The
  // build with more than one thread cannot run without it.
  const threaded =
    typeof SharedArrayBuffer !== "undefined" && globalThis.crossOriginIsolated === true;

  const name = threaded ? "yzma_wasm_mt" : "yzma_wasm";

  globalThis.yzmaThreaded = threaded;

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
