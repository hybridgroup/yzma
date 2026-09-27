# Benchmarks

yzma benchmarks, one file per platform.

| Platform | Numbers |
| --- | --- |
| Linux | [benchmarks/linux.md](benchmarks/linux.md) |
| macOS | [benchmarks/macos.md](benchmarks/macos.md) |
| Windows | [benchmarks/windows.md](benchmarks/windows.md) |
| WebAssembly | [benchmarks/webassembly.md](benchmarks/webassembly.md) |

yzma compared with model servers on the same machine.

| Comparison | Numbers |
| --- | --- |
| yzma, ollama, Docker Model Runner | [benchmarks/comparison.md](benchmarks/comparison.md) |

Each file has a table per suite with the median of five runs, and the raw
output of each run below the tables.

To reproduce these numbers on your machine, or to add a new machine, see
[how to run the benchmarks](benchmarks/README.md).

```shell
make download-benchmark-models
./benchmarks/run.sh
```

The benchmark code is in
[pkg/llama/benchmark_test.go](pkg/llama/benchmark_test.go),
[pkg/mtmd/benchmark_test.go](pkg/mtmd/benchmark_test.go) and
[benchmarks/compare](benchmarks/compare).
