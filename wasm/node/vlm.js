// vlm.js answers a question about an image in Node, with no browser.
//
// Node has no canvas, so this makes the pixels itself: a plain shape on a plain
// ground. That is enough to prove the whole path works, from the pixels through
// the projector to the tokens. Whether the answer is a good description of the
// shape is a question for a real image in a browser.
//
// Usage:
//   node wasm/node/vlm.js --dir build/wasm --model <model.gguf> \
//       --mmproj <mmproj.gguf> [--tokens 24] [--ppm <file.ppm>] [--side 448]
//       [--mt] [--webgpu]
//
// --ppm takes a real image in the binary PPM format (P6), which needs no
// decoder: convert one with "magick in.jpg out.ppm" or "ffmpeg -i in.jpg out.ppm".

const fs = require("node:fs");
const path = require("node:path");

function option(name, fallback) {
  const i = process.argv.indexOf("--" + name);
  return i >= 0 && i + 1 < process.argv.length ? process.argv[i + 1] : fallback;
}

const dir = path.resolve(option("dir", "build/wasm"));
const modelFile = option("model", "");
const mmprojFile = option("mmproj", "");
const prompt = option("prompt", "Describe this image.");
const maxTokens = parseInt(option("tokens", "24"), 10);
const ppmFile = option("ppm", "");
const webgpu = process.argv.includes("--webgpu");
const mt = process.argv.includes("--mt");
const side = parseInt(option("side", "448"), 10);

if (!modelFile || !mmprojFile) {
  console.error("give a model with --model and a projector with --mmproj");
  process.exit(2);
}

// This harness stands in for yzma-loader.js. Node has no WebGPU, so asking for
// it must fall back the same way the loader does.
let moduleName = "yzma_wasm.js";
if (mt) {
  moduleName = "yzma_wasm_mt.js";
}
if (webgpu) {
  if (globalThis.navigator && globalThis.navigator.gpu) {
    moduleName = "yzma_wasm_webgpu.js";
  } else {
    console.log("[loader] no WebGPU here, using the build on the CPU");
  }
}

// makeImage draws a dark square in the middle of a light ground and gives the
// RGBA that a canvas would give.
function makeImage(width, height) {
  const rgba = new Uint8Array(width * height * 4);
  for (let y = 0; y < height; y++) {
    for (let x = 0; x < width; x++) {
      const inside =
        x > width / 4 && x < (width * 3) / 4 && y > height / 4 && y < (height * 3) / 4;
      const value = inside ? 32 : 224;
      const i = (y * width + x) * 4;
      rgba[i] = value;
      rgba[i + 1] = value;
      rgba[i + 2] = value;
      rgba[i + 3] = 255;
    }
  }
  return { width, height, rgba };
}

// readPPM reads a binary PPM (P6), which is a header of plain text and then RGB.
function readPPM(file) {
  const bytes = fs.readFileSync(file);

  const fields = [];
  let i = 0;
  while (fields.length < 4) {
    while (i < bytes.length && /\s/.test(String.fromCharCode(bytes[i]))) i++;
    if (String.fromCharCode(bytes[i]) === "#") {
      while (i < bytes.length && bytes[i] !== 0x0a) i++;
      continue;
    }
    let start = i;
    while (i < bytes.length && !/\s/.test(String.fromCharCode(bytes[i]))) i++;
    fields.push(bytes.toString("ascii", start, i));
  }
  i++; // the single whitespace after the last field

  const [magic, w, h, maxValue] = fields;
  if (magic !== "P6" || maxValue !== "255") {
    throw new Error("only a binary PPM with 255 levels works here, got " + magic + "/" + maxValue);
  }

  const width = parseInt(w, 10);
  const height = parseInt(h, 10);
  const rgb = bytes.subarray(i, i + width * height * 3);

  // The program takes RGBA, the same as a canvas gives.
  const rgba = new Uint8Array(width * height * 4);
  for (let p = 0; p < width * height; p++) {
    rgba[p * 4] = rgb[p * 3];
    rgba[p * 4 + 1] = rgb[p * 3 + 1];
    rgba[p * 4 + 2] = rgb[p * 3 + 2];
    rgba[p * 4 + 3] = 255;
  }
  return { width, height, rgba };
}

const output = [];

let programIsReady;
const programReady = new Promise((resolve) => {
  programIsReady = resolve;
});

globalThis.yzmaOnMessage = (message) => {
  if (message.kind === "ready" || message.kind === "error") {
    programIsReady();
  }
  if (message.kind === "token") {
    output.push(message.text);
    process.stdout.write(message.text);
    return;
  }
  console.log("[" + message.kind + "] " + message.text);
};

async function main() {
  globalThis.crossOriginIsolated = false;
  globalThis.yzmaBase = dir;

  const factory = require(path.join(dir, moduleName));
  const llamaModule = await factory({
    locateFile: (file) => path.join(dir, file),
    print: () => {},
    printErr: () => {},
  });

  globalThis.yzmaModule = llamaModule;
  globalThis.yzmaReady = Promise.resolve(llamaModule);
  globalThis.yzmaThreaded = mt;
  globalThis.yzmaBackend = moduleName.includes("webgpu") ? "webgpu" : "cpu";

  // A browser gets both files over the network. Here they go straight in.
  llamaModule.FS.mkdirTree("/models");
  llamaModule.FS.writeFile("/models/model.gguf", fs.readFileSync(modelFile));
  llamaModule.FS.writeFile("/models/mmproj.gguf", fs.readFileSync(mmprojFile));

  require(path.join(dir, "wasm_exec.js"));

  const go = new Go();
  const binary = fs.readFileSync(path.join(dir, "yzma-vlm.wasm"));
  const result = await WebAssembly.instantiate(binary, go.importObject);

  go.run(result.instance);
  await programReady;

  const image = ppmFile ? readPPM(ppmFile) : makeImage(side, side);
  console.log("[image] " + image.width + " by " + image.height + (ppmFile ? " from " + ppmFile : " made here"));

  const done = new Promise((resolve) => {
    const previous = globalThis.yzmaOnMessage;
    globalThis.yzmaOnMessage = (message) => {
      previous(message);
      if (message.kind === "loaded") {
        globalThis.yzmaDescribe(prompt, image.width, image.height, image.rgba, maxTokens);
      }
      if (message.kind === "done" || message.kind === "error") {
        resolve(message);
      }
    };
  });

  // The files are in place, so this only opens them. A browser gets them over
  // the network with yzmaLoadModel instead.
  globalThis.yzmaOpenModel();

  const last = await done;
  console.log();

  if (last.kind === "error") {
    process.exit(1);
  }
  if (output.join("").trim().length === 0) {
    console.error("no tokens came out");
    process.exit(1);
  }
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
