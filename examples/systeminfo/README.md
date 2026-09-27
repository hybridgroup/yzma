# systeminfo

Shows system information, the devices that llama.cpp finds, and how yzma sees
the CPU cores on this machine.

## Running

```shell
$ go run ./examples/systeminfo/
-- Devices --
Device 0: CUDA0
Device 1: CPU

-- CPU Threads --
Logical CPUs:      32
Inference threads: 8
Performance CPUs:  [0 2 4 6 8 10 12 14]
Thread pool:       8 threads, each pinned to a CPU

-- llama.cpp System Information --
CUDA : ARCHS = 860,890 | USE_GRAPHS = 1 | PEER_MAX_BATCH_SIZE = 128 | CPU : SSE3 = 1 | SSSE3 = 1 | AVX = 1 | AVX_VNNI = 1 | AVX2 = 1 | F16C = 1 | FMA = 1 | BMI2 = 1 | LLAMAFILE = 1 | OPENMP = 1 | REPACK = 1 |
```

`Inference threads` is the thread count a context uses unless the program changes it.
`Thread pool` shows whether yzma can pin each thread to its own CPU, which it
cannot do on macOS.
