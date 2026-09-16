// loader-test.js tests the choice that yzma-loader.js makes, with no browser
// and no llama.cpp. It gives the loader a false document and false builds, so
// the tests cover the paths that need a GPU that computes wrong values.
//
// Run it with: node wasm/node/loader-test.js

const assert = require("node:assert");
const path = require("node:path");
const vm = require("node:vm");
const fs = require("node:fs");

const source = fs.readFileSync(path.join(__dirname, "..", "yzma-loader.js"), "utf8");

// fakeWorld makes the globals that the loader reads. The adapter and the result
// of the self test come from the options, thus each test describes one machine.
function fakeWorld({ mode, gpu, f16, jspi, threads, check, agent }) {
  const loaded = [];
  const world = {
    console: { log() {}, warn() {}, error() {} },
    yzmaMode: mode,
    crossOriginIsolated: threads === true,
    navigator: { userAgent: agent || "Chrome/152", hardwareConcurrency: 8 },
  };

  if (threads === true) {
    world.SharedArrayBuffer = ArrayBuffer;
  }
  if (jspi !== false) {
    world.WebAssembly = { Suspending: function () {}, promising: function () {} };
  } else {
    world.WebAssembly = {};
  }
  if (gpu) {
    world.navigator.gpu = {
      requestAdapter: async () => ({
        features: new Set(f16 === false ? [] : ["shader-f16"]),
        info: { vendor: "test", device: "test gpu" },
      }),
    };
  }

  // The loader appends a script tag. Answer it with a factory that gives a
  // false module of the build that the name asks for.
  world.document = {
    head: {
      appendChild(element) {
        const name = path.basename(element.src).replace(/\.js$/, "");
        loaded.push(name);
        world.yzmaModule = async () => {
          const instance = { name };
          if (name === "yzma_wasm_webgpu" && check !== "missing") {
            instance._yzma_backend_init = async () => {};
            instance._yzma_backend_check = async () =>
              check === "wrong" ? 1 : check === "throws" ? Promise.reject(new Error("boom")) : 0;
          }
          return instance;
        };
        element.onload();
      },
    },
    createElement: () => ({}),
  };

  world.globalThis = world;
  return { world, loaded };
}

async function run({ name, ...options }) {
  const { world, loaded } = fakeWorld(options);
  vm.createContext(world);
  vm.runInContext(source, world);
  const instance = await world.yzmaReady;
  return { name, instance, loaded, world };
}

async function main() {
  let failures = 0;
  const check = (label, got, want) => {
    try {
      assert.deepStrictEqual(got, want);
      console.log("  ok    " + label);
    } catch (e) {
      failures++;
      console.log("  FAIL  " + label + ": got " + JSON.stringify(got) + ", want " + JSON.stringify(want));
    }
  };

  // A GPU that agrees with the CPU keeps the build on the GPU.
  {
    const r = await run({ gpu: true, check: "ok" });
    console.log("a GPU that computes correct values");
    check("uses the GPU build", r.instance.name, "yzma_wasm_webgpu");
    check("reports webgpu", r.world.yzmaBackend, "webgpu");
    check("loads one build", r.loaded.length, 1);
    check("no reason to reject", r.world.yzmaGPUReject, undefined);
  }

  // A GPU that computes wrong values gives the build on the CPU instead.
  {
    const r = await run({ gpu: true, check: "wrong" });
    console.log("a GPU that computes wrong values");
    check("falls back to the CPU build", r.instance.name, "yzma_wasm");
    check("reports cpu", r.world.yzmaBackend, "cpu");
    check("forgets the adapter", r.world.yzmaAdapter, "");
    check("loads both builds", r.loaded, ["yzma_wasm_webgpu", "yzma_wasm"]);
    check("gives the reason", typeof r.world.yzmaGPUReject, "string");
  }

  // The same machine with the headers for threads takes the faster CPU build.
  {
    const r = await run({ gpu: true, check: "wrong", threads: true });
    console.log("a GPU that computes wrong values, with threads");
    check("falls back to the thread build", r.instance.name, "yzma_wasm_mt");
    check("reports cpu-threads", r.world.yzmaBackend, "cpu-threads");
    check("takes the cores of the machine", r.world.yzmaThreads, 8);
    check("says it has threads", r.world.yzmaThreaded, true);
  }

  // A self test that throws counts as a GPU that is not usable.
  {
    const r = await run({ gpu: true, check: "throws" });
    console.log("a self test that throws");
    check("falls back to the CPU build", r.instance.name, "yzma_wasm");
    check("gives the reason", typeof r.world.yzmaGPUReject, "string");
  }

  // A build from before ABI version 8 has no self test, thus it stays.
  {
    const r = await run({ gpu: true, check: "missing" });
    console.log("a build with no self test");
    check("keeps the GPU build", r.instance.name, "yzma_wasm_webgpu");
    check("loads one build", r.loaded.length, 1);
  }

  // A machine with no GPU never runs the self test.
  {
    const r = await run({ gpu: false, check: "wrong" });
    console.log("a machine with no GPU");
    check("uses the CPU build", r.instance.name, "yzma_wasm");
    check("loads one build", r.loaded.length, 1);
  }

  // Mode cpu asks for no GPU at all.
  {
    const r = await run({ mode: "cpu", gpu: true, check: "wrong" });
    console.log("mode cpu");
    check("uses the CPU build", r.instance.name, "yzma_wasm");
    check("loads one build", r.loaded.length, 1);
  }

  console.log(failures === 0 ? "\nall tests pass" : "\n" + failures + " tests fail");
  process.exit(failures === 0 ? 0 : 1);
}

main();
