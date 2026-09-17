// browser-bench.js measures yzma in a browser. A browser cannot be driven from
// a script here, thus paste this file in the console of the page that
// `make serve-wasm` gives, and put the result in webassembly.md.
//
// Change model and mode below. The modes are auto, cpu and webgpu, which
// yzma-loader.js reads.
//
// The output has the same shape as go test -bench, and it gives the llama.cpp
// build, which the page reads from yzma-install.json. Thus you can give it to
// cmd/yzma-bench:
//
//   go run ./cmd/yzma-bench update --file benchmarks/webassembly.md \
//     --suite browser --backend webgpu --arch wasm --machine <name> \
//     --label "<machine and browser>" --output run.txt
(async () => {
  const model =
    "https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf";
  const mode = "auto";
  const prompt = "Are you ready to go?";
  const maxTokens = 64;
  const count = 5;

  // installTag gives the llama.cpp build of the files that the page serves. A
  // nightly build has no upstream_tag, because its own tag names the assets.
  const installTag = async () => {
    try {
      const record = await (await fetch("./yzma-install.json")).json();
      return record.upstream_tag || record.tag || "";
    } catch {
      return "";
    }
  };

  const worker = new Worker("./worker.js?mode=" + mode);
  let backend = "";
  let waiting = () => {};
  worker.onmessage = (event) => waiting(event.data || {});

  const until = (kind) =>
    new Promise((resolve, reject) => {
      waiting = (message) => {
        if (message.kind === "ready") {
          backend = message.text;
        }
        if (message.kind === "error") {
          reject(new Error(message.text));
        }
        if (message.kind === kind) {
          resolve(message.text);
        }
      };
    });

  await until("ready");
  worker.postMessage({ kind: "load", url: model });
  await until("loaded");

  const generate = () => {
    const done = until("done");
    worker.postMessage({ kind: "generate", prompt, maxTokens });
    return done;
  };

  // The first generation warms the build up and is not in the result.
  await generate();

  const runs = [];
  const started = Date.now();
  for (let i = 0; i < count; i++) {
    const text = await generate();
    runs.push({
      tokens: parseInt(/(\d+) tokens/.exec(text)[1], 10),
      rate: parseFloat(/([\d.]+) tokens\/s/.exec(text)[1]),
    });
  }
  const elapsed = (Date.now() - started) / 1000;

  const threads = navigator.hardwareConcurrency || 1;
  const lines = [
    "$ browser-bench.js mode=" + mode + " tokens=" + maxTokens,
    "goos: js",
    "goarch: wasm",
    "pkg: github.com/hybridgroup/yzma/examples/wasm/chat",
    "cpu: " + navigator.userAgent,
    backend,
  ];
  const tag = await installTag();
  if (tag) {
    lines.push("llama.cpp: " + tag);
  } else {
    console.warn("no yzma-install.json, give the tag with --llamacpp");
  }
  for (const run of runs) {
    lines.push(
      [
        "BenchmarkInference-" + threads,
        "1",
        Math.round((run.tokens / run.rate) * 1e9) + " ns/op",
        run.rate.toFixed(1) + " tokens/s",
      ].join("\t"),
    );
  }
  lines.push("PASS");
  lines.push(
    "ok\tgithub.com/hybridgroup/yzma/examples/wasm/chat\t" + elapsed.toFixed(3) + "s",
  );

  const out = lines.join("\n");
  console.log(out);
  if (typeof copy === "function") {
    copy(out);
    console.log("the result is in the clipboard");
  }
})();
