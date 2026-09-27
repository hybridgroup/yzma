#!/usr/bin/env bash
# compare.sh measures yzma against a model server and writes each result to
# benchmarks/comparison.md. See benchmarks/README.md.
set -euo pipefail

root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$root"

machine=""
label=""
engines="yzma ollama dmr"
suites="text multimodal embeddings"
models="qwen3-vl-2b gemma4-e2b"
embed_model="bge-small"
llamacpp=""
nctx=8192
count=5
benchtime=20x
tokens=16
cool=60
dry_run=""

usage() {
  cat <<'EOF'
usage: benchmarks/compare.sh [flags]

  --machine NAME     short name of this machine, default the host name
  --label TEXT       display name of the machine, default the short name
  --engine LIST      yzma, ollama, dmr, or more than one
  --suite LIST       text, multimodal, embeddings, or more than one
  --model LIST       qwen3-vl-2b, gemma4-e2b, or more than one
  --device NAME      device for yzma, default CUDA0 when the machine has one
  --llamacpp TAG     tag of the llama.cpp build, default from yzma-install.json
  --nctx N           context tokens, default 8192
  --tokens N         tokens per request, default 16
  --count N          runs per benchmark, default 5
  --benchtime D      time or count per run, default 20x
  --cool C           wait until the GPU is at C degrees or less, default 60
  --dry-run          print the result and change no file

The servers must be running and must have the model. The script tells you
which command fetches a missing model.
EOF
}

device=""
device_given=""

while [ $# -gt 0 ]; do
  case "$1" in
    --machine) machine=$2; shift 2 ;;
    --label) label=$2; shift 2 ;;
    --engine) engines=$(echo "$2" | tr ',' ' '); shift 2 ;;
    --suite) suites=$(echo "$2" | tr ',' ' '); shift 2 ;;
    --model) models=$(echo "$2" | tr ',' ' '); shift 2 ;;
    --device) device=$2; device_given=1; shift 2 ;;
    --llamacpp) llamacpp=$2; shift 2 ;;
    --nctx) nctx=$2; shift 2 ;;
    --tokens) tokens=$2; shift 2 ;;
    --count) count=$2; shift 2 ;;
    --benchtime) benchtime=$2; shift 2 ;;
    --cool) cool=$2; shift 2 ;;
    --dry-run) dry_run="--dry-run"; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "unknown flag: $1" >&2; usage; exit 2 ;;
  esac
done

MODELS_DIR=${MODELS_DIR:-$HOME/models}
export YZMA_LIB=${YZMA_LIB:-$root/lib}

OLLAMA_URL=${OLLAMA_URL:-http://localhost:11434}
DMR_URL=${DMR_URL:-http://localhost:12434}

# ollama 0.34 doesn't follow the Hugging Face redirect to its CDN, so pulling
# hf.co/... fails. The script gives it the file from MODELS_DIR instead, which
# also makes sure every engine reads the same bytes.
# OLLAMA_CONTAINER is the ollama container, empty to use the host CLI.
# OLLAMA_MODELS_DIR is where MODELS_DIR is inside that container.
OLLAMA_CONTAINER=${OLLAMA_CONTAINER-ollama}
OLLAMA_MODELS_DIR=${OLLAMA_MODELS_DIR:-/models}
if [ -z "$OLLAMA_CONTAINER" ]; then
  OLLAMA_MODELS_DIR=$MODELS_DIR
fi

file=benchmarks/comparison.md
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

case "$(uname -s)" in
  Linux|Darwin) ;;
  *) echo "this script is for Linux and macOS" >&2; exit 2 ;;
esac

has() { command -v "$1" >/dev/null 2>&1; }

jsonValue() {
  sed -n 's/.*"'"$2"'"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' "$1" 2>/dev/null | head -1
}

installTag() {
  local record=$1/yzma-install.json
  if [ ! -f "$record" ]; then
    return 0
  fi

  local tag
  tag=$(jsonValue "$record" upstream_tag)
  if [ -z "$tag" ]; then
    tag=$(jsonValue "$record" tag)
  fi
  echo "$tag"
}

if [ -z "$machine" ]; then
  machine=${YZMA_BENCH_MACHINE:-$(hostname -s 2>/dev/null || hostname)}
fi
machine=$(echo "$machine" | tr '[:upper:] ' '[:lower:]-' | tr -cd 'a-z0-9._-')
if [ -z "$label" ]; then
  label=$machine
fi

yzma_version=$(sed -n 's/.*currentVersion = "\(.*\)".*/\1/p' version.go)

if [ -z "$llamacpp" ]; then
  llamacpp=$(installTag "$YZMA_LIB")
fi

# Every server engine uses the GPU, so yzma must use it too.
if [ -z "$device_given" ] && has nvidia-smi && nvidia-smi -L >/dev/null 2>&1; then
  device=CUDA0
fi

# modelFiles sets gguf, mmproj and ref for one short name. Every engine reads
# the same Hugging Face file. ref is the Docker Model Runner name, and v1.2
# rejects it unless it is lower case.
gguf=""
mmproj=""
ref=""
template=""
image_size=""
modelFiles() {
  template=""
  image_size=""
  case "$1" in
    qwen3-vl-2b)
      gguf=$MODELS_DIR/Qwen3-VL-2B-Instruct.Q4_K_M.gguf
      mmproj=$MODELS_DIR/Qwen3-VL-2B-Instruct.mmproj-Q8_0.gguf
      ref=hf.co/qwen/qwen3-vl-2b-instruct-gguf:q4_k_m
      # ollama has no renderer for qwen3vl, so it needs the template here.
      template='<|im_start|>user
{{ .Prompt }}<|im_end|>
<|im_start|>assistant
'
      # ollama scales a smaller image up to about 1000 tokens.
      image_size=1280x960
      ;;
    gemma4-e2b)
      gguf=$MODELS_DIR/gemma-4-E2B-it-Q4_K_M.gguf
      # The unsloth projector name doesn't include the model. Keep only one
      # gemma4 in MODELS_DIR, or their projectors collide.
      mmproj=$MODELS_DIR/mmproj-F16.gguf
      ref=hf.co/unsloth/gemma-4-e2b-it-gguf:q4_k_m
      # ollama has a renderer for gemma4, so it needs no template here.
      # This image stays under the 280 token budget of Docker Model Runner's
      # llama.cpp. A larger one gets resized differently by each engine.
      image_size=768x576
      ;;
    bge-small)
      # The embeddings suite. An embedding returns no tokens, so almost all of
      # the request cost is the round trip.
      gguf=$MODELS_DIR/bge-small-en-v1.5-q8_0.gguf
      mmproj=""
      ref=hf.co/ggml-org/bge-small-en-v1.5-q8_0-gguf:q8_0
      ;;
    smolvlm-256m)
      # The small model the other suites use. It gives a quick check of the
      # engines without downloading a large model.
      gguf=$MODELS_DIR/SmolVLM-256M-Instruct-Q8_0.gguf
      mmproj=$MODELS_DIR/mmproj-SmolVLM-256M-Instruct-Q8_0.gguf
      ref=hf.co/ggml-org/smolvlm-256m-instruct-gguf:q8_0
      # ollama has no renderer for idefics3, so it needs the template here.
      template='<|im_start|>User: {{ .Prompt }}<end_of_utterance>
Assistant:'
      ;;
    *)
      echo "unknown model: $1, see benchmarks/README.md" >&2
      exit 2
      ;;
  esac
}

# dmrContextIsSet reports whether this run's model has the context size the
# benchmark needs. Only that model counts, so the check matches on its name.
dmrContextIsSet() {
  docker model configure show 2>/dev/null |
    python3 -c "
import sys, json
want = int(sys.argv[1])
name = sys.argv[2].split('/', 1)[-1].lower()
try:
    rows = json.load(sys.stdin)
except Exception:
    sys.exit(1)
for row in rows:
    if name in row.get('Model', '').lower():
        sys.exit(0 if row.get('Config', {}).get('context-size') == want else 1)
sys.exit(1)
" "$nctx" "$ref"
}

# freeGPU unloads every model. A server keeps its model loaded, so the next
# engine would find no free GPU memory and fail to load.
freeGPU() {
  docker model unload --all >/dev/null 2>&1 || true

  # ollama ps has no --format, so awk extracts each model name.
  if [ -n "$OLLAMA_CONTAINER" ] || command -v ollama >/dev/null 2>&1; then
    ollamaRun sh -c "ollama ps 2>/dev/null | awk 'NR>1 && \$1!=\"\" {print \$1}' |
      while read -r m; do ollama stop \"\$m\"; done" >/dev/null 2>&1 || true
  fi

  # GPU memory frees up a moment after the process releases it.
  sleep 3

  coolGPU
}

# coolGPU waits until the GPU is at $cool degrees or less. A hot laptop GPU
# lowers its clock, so an engine that runs right after another is slower.
coolGPU() {
  if ! has nvidia-smi; then
    return 0
  fi

  local temp waited=0
  while temp=$(nvidia-smi --query-gpu=temperature.gpu --format=csv,noheader,nounits 2>/dev/null | head -1) &&
    [ -n "$temp" ] && [ "$temp" -gt "$cool" ]; do
    if [ "$waited" -ge 600 ]; then
      echo "the GPU is still at $temp degrees after 10 minutes, continuing" >&2
      return 0
    fi
    if [ "$waited" -eq 0 ]; then
      echo "==> the GPU is at $temp degrees, waiting until it is at $cool"
    fi
    sleep 10
    waited=$((waited + 10))
  done
}

# ollamaRun runs one ollama command, in the container or on the host.
ollamaRun() {
  if [ -n "$OLLAMA_CONTAINER" ]; then
    docker exec -i "$OLLAMA_CONTAINER" "$@"
  else
    "$@"
  fi
}

# ollamaImport creates the ollama model from the file in MODELS_DIR. It does
# nothing when ollama already has it.
ollamaImport() {
  local name=$1 wantChat=${2-yes}
  if ollamaRun ollama list 2>/dev/null | grep -q "^$name"; then
    return 0
  fi

  local modelfile="FROM $OLLAMA_MODELS_DIR/$(basename "$gguf")"
  if [ -f "$mmproj" ]; then
    modelfile="$modelfile
FROM $OLLAMA_MODELS_DIR/$(basename "$mmproj")"
  fi
  if [ -n "$template" ]; then
    modelfile="$modelfile
TEMPLATE \"\"\"$template\"\"\""
  fi
  # ollama defaults to 4096 context. Every engine must use the same size, or
  # they don't give the GPU the same work.
  modelfile="$modelfile
PARAMETER num_ctx $nctx"

  echo "==> ollama is importing $name from $(basename "$gguf")"
  printf '%s\n' "$modelfile" |
    ollamaRun sh -c "cat > /tmp/Modelfile.$name && ollama create $name -f /tmp/Modelfile.$name" || return 1

  if [ "$wantChat" = no ]; then
    return 0
  fi

  ollamaHasChatFormat "$name"
}

# ollamaHasChatFormat makes sure ollama frames the chat. Without a template or
# an architecture renderer, ollama sends the bare prompt with no user and
# assistant turns, and the model answers differently from the other engines.
ollamaHasChatFormat() {
  local shown
  shown=$(ollamaRun ollama show --modelfile "$1" 2>/dev/null)

  if echo "$shown" | grep -q '^RENDERER '; then
    return 0
  fi
  if echo "$shown" | grep -q '^TEMPLATE ' && ! echo "$shown" | grep -qF 'TEMPLATE {{ .Prompt }}'; then
    return 0
  fi

  echo "ollama has no chat template and no renderer for $1." >&2
  echo "It would send the prompt with no user and assistant turns." >&2
  echo "Add a template for this model to modelFiles in benchmarks/compare.sh." >&2

  return 1
}

engineVersion() {
  case "$1" in
    yzma)
      echo "$yzma_version"
      ;;
    ollama)
      curl -s --max-time 5 "$OLLAMA_URL/api/version" 2>/dev/null |
        sed -n 's/.*"version"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -1
      ;;
    dmr)
      # The server runs the model, so its version goes in the table.
      # v1.2 prints a Client block and a Server block. v0.1 printed one line.
      docker model version 2>/dev/null |
        awk '/^Server:/{s=1} s&&/Version:/{print $2; exit}' | grep -oE 'v[0-9.]+' | head -1 ||
        docker model version 2>/dev/null | grep -oE 'v[0-9]+\.[0-9]+\.[0-9]+' | head -1
      ;;
  esac
}

# ready reports whether the engine can run the model, and what to do if not.
ready() {
  local engine=$1

  case "$engine" in
    yzma)
      if [ ! -f "$gguf" ]; then
        echo "no model at $gguf, run make download-compare-models" >&2
        return 1
      fi
      if [ "$suite" = multimodal ] && [ ! -f "$mmproj" ]; then
        echo "no projector at $mmproj, run make download-compare-models" >&2
        return 1
      fi
      ;;
    ollama)
      if ! curl -s --max-time 5 "$OLLAMA_URL/api/version" >/dev/null 2>&1; then
        echo "ollama is not answering at $OLLAMA_URL, start it" >&2
        return 1
      fi
      if [ ! -f "$gguf" ]; then
        echo "no model at $gguf, run make download-compare-models" >&2
        return 1
      fi
      local wantChat=yes
      if [ "$suite" = embeddings ]; then
        wantChat=no
      fi
      if ! ollamaImport "$local_name" "$wantChat"; then
        echo "ollama could not import $gguf, is $MODELS_DIR mounted at $OLLAMA_MODELS_DIR" >&2
        return 1
      fi
      ;;
    dmr)
      if ! curl -s --max-time 5 "$DMR_URL/engines/v1/models" >/dev/null 2>&1; then
        echo "Docker Model Runner is not answering at $DMR_URL, run: docker model status" >&2
        return 1
      fi
      # Ask the server, not the CLI. v1.2 lists huggingface.co but takes hf.co
      # in a request, so the listed name is no use as a check.
      if ! curl -s --max-time 10 "$DMR_URL/engines/v1/models" |
        grep -qiF "${ref#*/}"; then
        echo "Docker Model Runner does not have $model, run: docker model pull $ref" >&2
        return 1
      fi
      ;;
  esac

  return 0
}

# flagsFor returns the go test flags for one engine and suite.
flagsFor() {
  local engine=$1 suite=$2
  local flags=(-engine="$engine" -suite="$suite" -tokens="$tokens" -nctx="$nctx")
  if [ "$suite" = multimodal ] && [ -n "$image_size" ]; then
    flags+=(-image-size="$image_size")
  fi

  case "$engine" in
    yzma)
      flags+=(-model="$gguf")
      if [ "$suite" = multimodal ]; then
        flags+=(-mmproj="$mmproj")
      fi
      if [ -n "$device" ]; then
        flags+=(-device="$device")
      fi
      ;;
    ollama)
      flags+=(-server-model="$local_name" -ollama-url="$OLLAMA_URL/v1")
      ;;
    dmr)
      flags+=(-server-model="$ref" -dmr-url="$DMR_URL/engines/v1")
      ;;
  esac

  echo "${flags[@]}"
}

info=$work/device.txt
if has nvidia-smi; then nvidia-smi > "$info" 2>/dev/null || true; else : > "$info"; fi

answers=$work/answers.txt
: > "$answers"

# runCombination MODEL ENGINE SUITE runs one measurement and records it.
runCombination() {
  local model=$1 engine=$2 suite=$3

  modelFiles "$model"
  local_name=yzma-bench-$model

  if ! ready "$engine"; then
    echo "==> $model, $engine is not ready, skipping" >&2
    return 0
  fi

  local version
  version=$(engineVersion "$engine")
  if [ -z "$version" ]; then
    echo "==> no version for $engine, skipping" >&2
    return 0
  fi

  local flags
  read -r -a flags <<< "$(flagsFor "$engine" "$suite")"

  echo "==> $model, $engine, $suite"
  freeGPU

  # A long context model needs a limit, or Docker Model Runner asks the GPU for
  # a KV cache of many GB and llama.cpp fails. "docker model unload" forgets
  # the setting, so this must come after freeGPU.
  if [ "$engine" = dmr ]; then
    if [ "$suite" = embeddings ]; then
      docker model configure "$ref" --mode embedding >/dev/null 2>&1 || true
    else
      docker model configure "$ref" --context-size "$nctx" >/dev/null 2>&1 || true
      if ! dmrContextIsSet; then
        echo "the context size of $ref was not applied, skipping dmr" >&2
        return 0
      fi
    fi
  fi

  # Every engine must give the same answer, or the benchmark numbers compare
  # different work.
  (cd "$root/benchmarks/compare" && go test -run TestAnswer -v "${flags[@]}" 2>&1) |
    grep '^answer ' | sed "s/^/$model /" >> "$answers" || true

  local out=$work/$model-$engine-$suite.txt
  if ! {
    echo "\$ cd benchmarks/compare && go test -benchtime=$benchtime -count=$count -run=nada -bench BenchmarkCompare ${flags[*]}"
    (cd "$root/benchmarks/compare" && go test -benchtime="$benchtime" -count="$count" \
      -run=nada -bench BenchmarkCompare "${flags[@]}" 2>&1)
  } | tee "$out"; then
    echo "the $model run on $engine failed, the file keeps its current numbers" >&2
    return 0
  fi

  local args=(go run ./cmd/yzma-bench update --file "$file" --suite "compare-$suite"
    --backend "$engine" --machine "$machine" --label "$label" --model "$model"
    --engine-version "$version" --yzma "$yzma_version" --output "$out"
    --notes "$tokens tokens, greedy sampling, one request at a time." --device-info "$info")
  if [ "$engine" = yzma ] && [ -n "$llamacpp" ]; then
    args+=(--llamacpp "$llamacpp")
  fi
  if [ -n "$dry_run" ]; then
    args+=("$dry_run")
  fi

  "${args[@]}"
}

# The text and image suites use the models in $models. The embeddings suite
# uses its own embedding model, so it runs separately.
for suite in $suites; do
  if [ "$suite" = embeddings ]; then
    for engine in $engines; do
      runCombination "$embed_model" "$engine" embeddings
    done
    continue
  fi

  for model in $models; do
    for engine in $engines; do
      runCombination "$model" "$engine" "$suite"
    done
  done
done
if [ -s "$answers" ]; then
  echo
  echo "The engine answers. They must match."
  cat "$answers"

  # The prompt token count is the check that matters. Equal answers can hide a
  # difference, but a different count always means different work.
  echo
  awk '
    match($0, /answer [a-z]+\/[a-z]+: [0-9]+ prompt tokens/) {
      split($3, part, "/")
      engine = part[1]
      suite = part[2]
      sub(":", "", suite)
      key = $1 "/" suite
      count = $4
      if (!(key in seen)) { seen[key] = count; first[key] = engine; next }
      if (seen[key] != count) {
        bad[key] = bad[key] sprintf("  %s has %s, %s has %s\n", first[key], seen[key], engine, count)
      }
    }
    END {
      n = 0
      for (key in bad) {
        if (n == 0) print "The engines do not do the same work."
        printf "%s:\n%s", key, bad[key]
        n++
      }
      if (n == 0) {
        print "Every engine sends the same number of prompt tokens."
      } else {
        print ""
        print "A different count means a different prompt, or different image"
        print "preprocessing. The table shows the count in its own column, so a"
        print "reader can see which engines did the same work."
      }
    }
  ' "$answers"
fi
