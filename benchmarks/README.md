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

On Windows. Use a PowerShell prompt in the directory of the repository. A
Command Prompt opens the file in an editor and does not run it. PowerShell does
not run a script until the policy of the machine permits it. This command gives
the policy for one run.

```powershell
powershell -ExecutionPolicy Bypass -File .\benchmarks\run.ps1
```

A full run is long. The first go command builds the packages, and each suite
takes many minutes. The script shows each command and each line of the
benchmarks when they come.

The script asks llama.cpp which devices the machine has, runs the text
benchmark and the multimodal benchmark for each one, and puts each result in the
file of the platform. The tag of the llama.cpp build comes from
`yzma-install.json` of the library directory.

Useful flags.

```shell
./benchmarks/run.sh --backend vulkan        # one backend only
./benchmarks/run.sh --suite text            # one suite only
./benchmarks/run.sh --machine jetson-orin-nano --label "Jetson Orin Nano 8GB"
./benchmarks/run.sh --llamacpp b10964       # when the library came from elsewhere
./benchmarks/run.sh --dry-run               # print the result, change no file
```

The PowerShell script takes the same names with one dash and a capital, as
`-Machine`, `-Backend`, `-DryRun` and so on. The flags go after the name of the
file.

```powershell
powershell -ExecutionPolicy Bypass -File .\benchmarks\run.ps1 -Backend vulkan -DryRun
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
the result to the tool. The output gives the llama.cpp build, which the page
reads from `yzma-install.json` of the build directory.

```shell
make serve-wasm
go run ./cmd/yzma-bench update --file benchmarks/webassembly.md \
  --suite browser --backend webgpu --arch wasm \
  --machine <name> --label "<machine and browser>" --output run.txt
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

Each section must say which llama.cpp build made it. The tool refuses a result
that has no tag, thus the numbers of a table always compare the same thing.

To delete a result that is not comparable any more, for example after a change
of the model, give its key to `remove`.

```shell
go run ./cmd/yzma-bench remove --file benchmarks/linux.md \
  multimodal/cuda/amd64/i9-13900hx/cuda0
```

## What is measured

Both benchmarks report tokens a second with `b.ReportMetric`, and the table
gives the median of the runs.

| Suite | Code | Model |
| --- | --- | --- |
| text | [pkg/llama/benchmark_test.go](../pkg/llama/benchmark_test.go) | SmolLM-135M Q2_K |
| multimodal | [pkg/mtmd/benchmark_test.go](../pkg/mtmd/benchmark_test.go) | SmolVLM-256M-Instruct Q8_0 with its projector |
| node, browser | [examples/wasm/chat](../examples/wasm/chat) | SmolLM-135M Q2_K |

The WebAssembly numbers come from the generation loop of the example, thus they
are not comparable with the native tables.
