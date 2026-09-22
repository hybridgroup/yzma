# systeminfo

Shows system information, the devices that llama.cpp finds and what yzma makes
of the cores of the machine.

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
Thread pool:       8 threads, each held to a CPU

-- llama.cpp System Information --
CUDA : ARCHS = 860,890 | USE_GRAPHS = 1 | PEER_MAX_BATCH_SIZE = 128 | CPU : SSE3 = 1 | SSSE3 = 1 | AVX = 1 | AVX_VNNI = 1 | AVX2 = 1 | F16C = 1 | FMA = 1 | BMI2 = 1 | LLAMAFILE = 1 | OPENMP = 1 | REPACK = 1 |
```

`Inference threads` is what a context takes unless the program changes it.
`Thread pool` says whether yzma can hold each thread to a CPU of its own,
which it cannot do on macOS.
