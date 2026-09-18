# Benchmarks

yzma benchmarks, one file for each platform.

| Platform | Numbers |
| --- | --- |
| Linux | [benchmarks/linux.md](benchmarks/linux.md) |
| macOS | [benchmarks/macos.md](benchmarks/macos.md) |
| Windows | [benchmarks/windows.md](benchmarks/windows.md) |
| WebAssembly | [benchmarks/webassembly.md](benchmarks/webassembly.md) |

yzma against a model server, on one machine.

| Comparison | Numbers |
| --- | --- |
| yzma, ollama, Docker Model Runner | [benchmarks/comparison.md](benchmarks/comparison.md) |

Each file has a table for each suite, with the median of five runs, and the
output of each run below the tables.

To make these numbers again on your machine, or to add a machine, see
[how to run the benchmarks](benchmarks/README.md).

```shell
make download-benchmark-models
./benchmarks/run.sh
```

The benchmark code is in
[pkg/llama/benchmark_test.go](pkg/llama/benchmark_test.go),
[pkg/mtmd/benchmark_test.go](pkg/mtmd/benchmark_test.go) and
[benchmarks/compare](benchmarks/compare).
