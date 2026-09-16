#!/usr/bin/env bash
# run.sh runs the benchmarks of this machine and puts each result in the
# markdown file of the platform. See benchmarks/README.md.
set -euo pipefail

root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$root"

machine=""
label=""
backends="all"
suites="text multimodal"
llamacpp=""
nctx=""
count=5
benchtime=10s
dry_run=""
tokens=64

usage() {
  cat <<'EOF'
usage: benchmarks/run.sh [flags]

  --machine NAME     short name of this machine, default the host name
  --label TEXT       name of the machine to show, default the CPU
  --backend LIST     all, cpu, cuda, rocm, vulkan, metal, wasm
  --suite LIST       text, multimodal, or both
  --llamacpp TAG     tag of the llama.cpp build, default from yzma-install.json
  --nctx N           context tokens, default 8192 on the CPU and 32000 on a GPU
  --count N          runs of each benchmark, default 5
  --benchtime D      time of each run, default 10s
  --tokens N         tokens of each run, WebAssembly only, default 64
  --dry-run          print the result and change no file
EOF
}

while [ $# -gt 0 ]; do
  case "$1" in
    --machine) machine=$2; shift 2 ;;
    --label) label=$2; shift 2 ;;
    --backend) backends=$2; shift 2 ;;
    --suite) suites=$2; shift 2 ;;
    --llamacpp) llamacpp=$2; shift 2 ;;
    --nctx) nctx=$2; shift 2 ;;
    --count) count=$2; shift 2 ;;
    --benchtime) benchtime=$2; shift 2 ;;
    --tokens) tokens=$2; shift 2 ;;
    --dry-run) dry_run="--dry-run"; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "unknown flag: $1" >&2; usage; exit 2 ;;
  esac
done

MODELS_DIR=${MODELS_DIR:-$HOME/models}
export YZMA_LIB=${YZMA_LIB:-$root/lib}
export YZMA_BENCHMARK_MODEL=${YZMA_BENCHMARK_MODEL:-$MODELS_DIR/SmolLM-135M.Q2_K.gguf}
export YZMA_BENCHMARK_MMMODEL=${YZMA_BENCHMARK_MMMODEL:-$MODELS_DIR/Qwen3-VL-2B-Instruct.Q4_K_M.gguf}
export YZMA_BENCHMARK_MMPROJ=${YZMA_BENCHMARK_MMPROJ:-$MODELS_DIR/Qwen3-VL-2B-Instruct.mmproj-Q8_0.gguf}

WASM_DIR=${WASM_DIR:-$root/build/wasm}
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

case "$(uname -s)" in
  Linux) file=benchmarks/linux.md ;;
  Darwin) file=benchmarks/macos.md ;;
  *) echo "this script is for Linux and macOS, use run.ps1 on Windows" >&2; exit 2 ;;
esac

# jsonValue takes one value of a flat JSON file, thus the script needs no jq.
jsonValue() {
  sed -n 's/.*"'"$2"'"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$1" 2>/dev/null | head -1
}

if [ -z "$machine" ]; then
  machine=${YZMA_BENCH_MACHINE:-$(hostname -s 2>/dev/null || hostname)}
fi
machine=$(echo "$machine" | tr '[:upper:] ' '[:lower:]-' | tr -cd 'a-z0-9._-')

has() { command -v "$1" >/dev/null 2>&1; }

# deviceInfo collects what the machine says about the device of a backend.
deviceInfo() {
  case "$1" in
    cuda)
      if has nvidia-smi; then nvidia-smi; fi
      ;;
    rocm)
      if has amdgpu_top; then amdgpu_top -d
      elif has rocm-smi; then rocm-smi
      fi
      ;;
    vulkan)
      if has vulkaninfo; then vulkaninfo --summary; fi
      ;;
    metal)
      if has system_profiler; then system_profiler SPDisplaysDataType; fi
      ;;
  esac

  return 0
}

update() {
  # update SUITE BACKEND DEVICE OUTPUT DEVICEINFO
  local args=(go run ./cmd/yzma-bench update --file "$file" --suite "$1"
    --backend "$2" --device "$3" --machine "$machine" --label "$label"
    --llamacpp "$llamacpp" --yzma "$yzma_version" --output "$4")
  if [ -s "$5" ]; then
    args+=(--device-info "$5")
  fi
  if [ -n "$dry_run" ]; then
    args+=("$dry_run")
  fi
  "${args[@]}"
}

yzma_version=$(sed -n 's/.*currentVersion = "\(.*\)".*/\1/p' version.go)

# runNative BACKEND DEVICE INFO [RECORD]. DEVICE goes to go test, RECORD goes
# in the table. The CPU gets no device in the table, because a machine has one.
runNative() {
  local backend=$1 device=$2 info=$3 record=${4-$2}
  local ctx=$nctx
  if [ -z "$ctx" ]; then
    if [ "$backend" = cpu ]; then ctx=8192; else ctx=32000; fi
  fi

  local flag=""
  if [ -n "$device" ]; then
    flag="-device=$device"
  fi

  for suite in $suites; do
    local pkg bench
    case "$suite" in
      text) pkg=pkg/llama; bench=BenchmarkInference ;;
      multimodal) pkg=pkg/mtmd; bench=BenchmarkMultimodalInference ;;
      *) echo "unknown suite: $suite" >&2; exit 2 ;;
    esac
    if [ "$suite" = multimodal ] && [ ! -f "$YZMA_BENCHMARK_MMMODEL" ]; then
      echo "no multimodal model at $YZMA_BENCHMARK_MMMODEL, the suite is not run" >&2
      continue
    fi

    local out=$work/$suite-$backend-${device:-cpu}.txt
    echo "==> $suite, $backend ${device:+($device)}"
    # The benchmarks take -nctx and -device, which go test only gives to the
    # test binary of the package that it runs in.
    if ! {
      echo "\$ cd $pkg && go test -benchtime=$benchtime -count=$count -run=nada -bench $bench -nctx=$ctx ${flag:+$flag}"
      (cd "$root/$pkg" && go test -benchtime="$benchtime" -count="$count" -run=nada \
        -bench "$bench" -nctx="$ctx" ${flag:+"$flag"} 2>&1)
    } | tee "$out"; then
      echo "the run of $suite on $backend failed, the file keeps the numbers it has" >&2
      continue
    fi

    update "$suite" "$backend" "$record" "$out" "$info"
  done
}

runWasm() {
  file=benchmarks/webassembly.md
  # The WebAssembly modules come from their own install, thus the tag of the
  # native library does not apply here.
  if [ -z "$llamacpp" ] && [ -f "$WASM_DIR/yzma-install.json" ]; then
    llamacpp=$(jsonValue "$WASM_DIR/yzma-install.json" upstream_tag)
  fi
  if [ ! -f "$WASM_DIR/yzma.wasm" ]; then
    echo "no build in $WASM_DIR, run make download-llama.cpp-wasm and make wasm-example" >&2
    exit 2
  fi

  local empty=$work/empty.txt
  : > "$empty"

  for mode in cpu cpu-threads; do
    local out=$work/wasm-$mode.txt
    local flag=""
    if [ "$mode" = cpu-threads ]; then
      flag="--mt"
    fi
    echo "==> WebAssembly in Node, $mode"
    if ! node wasm/node/bench.js --dir "${WASM_DIR#"$root"/}" --model "$YZMA_BENCHMARK_MODEL" \
      --tokens "$tokens" --count "$count" $flag | tee "$out"; then
      echo "the run of $mode failed, the file keeps the numbers it has" >&2
      continue
    fi

    update node "$mode" "" "$out" "$empty"
  done
}

if [ "$backends" = wasm ]; then
  if [ -z "$label" ]; then
    label=$machine
  fi
  runWasm
  exit 0
fi

if [ -z "$llamacpp" ] && [ -f "$YZMA_LIB/yzma-install.json" ]; then
  llamacpp=$(jsonValue "$YZMA_LIB/yzma-install.json" upstream_tag)
fi

if [ ! -f "$YZMA_BENCHMARK_MODEL" ]; then
  echo "no model at $YZMA_BENCHMARK_MODEL, run make download-benchmark-models" >&2
  exit 2
fi

# The device list of llama.cpp says which backends this machine has.
devices=$work/devices.txt
go run . system -lib "$YZMA_LIB" > "$devices"

if [ -z "$label" ]; then
  label=$(awk '/Backend:[[:space:]]*CPU/{found=1} found && /Description:/{sub(/.*Description:[[:space:]]*/,""); print; exit}' "$devices")
  if [ -z "$label" ]; then
    label=$machine
  fi
fi

wanted() {
  if [ "$backends" = all ]; then
    return 0
  fi

  case " $(echo "$backends" | tr ',' ' ') " in
    *" $1 "*) return 0 ;;
  esac

  return 1
}

ran_cpu=""
device=""
backend=""
while IFS= read -r line; do
  case "$line" in
    "Device "*) device=${line#*: } ;;
    *"Backend:"*)
      backend=$(echo "${line#*Backend:}" | tr -d ' ' | tr '[:upper:]' '[:lower:]')
      info=$work/$backend-$device.txt
      if [ "$backend" = cpu ]; then
        if wanted cpu && [ -z "$ran_cpu" ]; then
          : > "$info"
          runNative cpu "$device" "$info" ""
          ran_cpu=1
        fi
      elif wanted "$backend"; then
        deviceInfo "$backend" > "$info" 2>/dev/null || true
        runNative "$backend" "$device" "$info"
      fi
      ;;
  esac
done < "$devices"
