# Benchmarks

These benchmarks all use the [SmolLM-135M-GGUF](https://huggingface.co/QuantFactory/SmolLM-135M-GGUF/resolve/main/SmolLM-135M.Q2_K.gguf) model to perform simple text generation.

See https://github.com/hybridgroup/yzma/blob/main/pkg/llama/benchmark_test.go

## Linux

### CPU

```
$ go test -bench=. -benchmem -benchtime=10s -count=5 -run=^$ ./pkg/llama
goos: linux                                    
goarch: amd64                               
pkg: github.com/hybridgroup/yzma/pkg/llama

cpu: 13th Gen Intel(R) Core(TM) i9-13900HX   
BenchmarkInference-32                 97         122892912 ns/op           16280 B/op        722 allocs/op
BenchmarkInference-32                 94         124091072 ns/op           16328 B/op        724 allocs/op
BenchmarkInference-32                 97         119683489 ns/op           16289 B/op        723 allocs/op
BenchmarkInference-32                 96         121095505 ns/op           16301 B/op        723 allocs/op
BenchmarkInference-32                 97         121729456 ns/op           16285 B/op        723 allocs/op
PASS            
ok      github.com/hybridgroup/yzma/pkg/llama   88.329s
```

### CUDA

```
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 580.95.05              Driver Version: 580.95.05      CUDA Version: 13.0     |
+-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  NVIDIA GeForce RTX 4070 ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   38C    P0            590W /  115W |      15MiB /   8188MiB |     17%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+
```

```
$ go test -bench=. -benchmem -benchtime=10s -count=5 -run=^$ ./pkg/llama
goos: linux                                                                                               
goarch: amd64                                                                                             
pkg: github.com/hybridgroup/yzma/pkg/llama                                                                
cpu: 13th Gen Intel(R) Core(TM) i9-13900HX
BenchmarkInference-32                344          34496552 ns/op           15320 B/op        695 allocs/op
BenchmarkInference-32                336          34741437 ns/op           15328 B/op        695 allocs/op
BenchmarkInference-32                340          34682374 ns/op           15327 B/op        695 allocs/op
BenchmarkInference-32                331          34697196 ns/op           15337 B/op        695 allocs/op
BenchmarkInference-32                340          34658749 ns/op           15328 B/op        695 allocs/op
PASS                                    
ok      github.com/hybridgroup/yzma/pkg/llama   81.720s
```

### Vulkan

Coming soon...

## macOS

### Metal

```
$ go test -bench=. -benchmem -benchtime=10s -count=5 -run=^$ -v ./pkg/llama
goos: darwin
goarch: arm64
pkg: github.com/hybridgroup/yzma/pkg/llama
cpu: Apple M4 Max
BenchmarkInference
BenchmarkInference-16                207          57066276 ns/op           15490 B/op        701 allocs/op
BenchmarkInference-16                207          57458342 ns/op           15495 B/op        701 allocs/op
BenchmarkInference-16                207          57461009 ns/op           15492 B/op        701 allocs/op
BenchmarkInference-16                208          56645718 ns/op           15490 B/op        701 allocs/op
BenchmarkInference-16                211          56484412 ns/op           15485 B/op        700 allocs/op
PASS
ok      github.com/hybridgroup/yzma/pkg/llama   86.878s
```

## Windows

### CPU

Coming soon...

### CUDA

Coming soon...

### Vulkan

Coming soon...
