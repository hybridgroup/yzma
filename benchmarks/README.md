# How to run the benchmarks

These numbers change with each llama.cpp release. One command makes them again
on any machine and puts each result in the correct file.

| File | What is in it |
| --- | --- |
| [linux.md](linux.md) | Linux, amd64 and arm64, CPU and GPU |
| [macos.md](macos.md) | macOS with Metal |
| [windows.md](windows.md) | Windows, CPU and GPU |
| [webassembly.md](webassembly.md) | WebAssembly in Node and in a browser |

## Before the first run

Get the library and the models.

```shell
make download-llama.cpp
make download-benchmark-models
```

## Run them

On Linux and macOS.

```shell
./benchmarks/run.sh
```

On Windows.

```powershell
.\benchmarks\run.ps1
```

The script asks llama.cpp which devices the machine has, runs the text
benchmark and the multimodal benchmark for each one, and puts each result in the
file of the platform. The tag of the llama.cpp build comes from
`yzma-install.json` of the library directory.

Useful flags. The PowerShell script takes the same names, as `-Machine`,
`-Backend` and so on.

```shell
./benchmarks/run.sh --backend vulkan        # one backend only
./benchmarks/run.sh --suite text            # one suite only
./benchmarks/run.sh --machine jetson-orin-nano --label "Jetson Orin Nano 8GB"
./benchmarks/run.sh --llamacpp b10964       # when the library came from elsewhere
./benchmarks/run.sh --dry-run               # print the result, change no file
```

The machine name is part of the key of a section. Give the same name each time,
or the file gets two sections for one machine. The default is the host name.

## WebAssembly

Node covers the build with one thread and the build with more threads.

```shell
make download-llama.cpp-wasm
make wasm-example
./benchmarks/run.sh --backend wasm
```

WebGPU needs a browser, which no script here can drive. Serve the example, open
the page, paste [browser-bench.js](browser-bench.js) in the console, and give
the result to the tool.

```shell
make serve-wasm
go run ./cmd/yzma-bench update --file benchmarks/webassembly.md \
  --suite browser --backend webgpu --arch wasm \
  --machine <name> --label "<machine and browser>" --llamacpp <tag> --output run.txt
```

## How a file is made

Each result is a section between markers, and the key is
`<suite>/<backend>/<arch>/<machine>[/<device>]`. A new run of the same key
replaces that section only. The table of each suite is a generated block, which
[cmd/yzma-bench](../cmd/yzma-bench) makes again from the sections. Thus a table
and its sections cannot disagree.

Do not edit a table by hand. To check a file after an edit:

```shell
go run ./cmd/yzma-bench check benchmarks/*.md
```

Sections that say `unknown` came from the old BENCHMARKS.md, which did not
record the llama.cpp build or the date. A run on that machine replaces them.

## What is measured

Both benchmarks report tokens a second with `b.ReportMetric`, and the table
gives the median of the runs.

| Suite | Code | Model |
| --- | --- | --- |
| text | [pkg/llama/benchmark_test.go](../pkg/llama/benchmark_test.go) | SmolLM-135M Q2_K |
| multimodal | [pkg/mtmd/benchmark_test.go](../pkg/mtmd/benchmark_test.go) | Qwen3-VL-2B-Instruct Q4_K_M with its projector |
| node, browser | [examples/wasm/chat](../examples/wasm/chat) | SmolLM-135M Q2_K |

The WebAssembly numbers come from the generation loop of the example, thus they
are not comparable with the native tables.
