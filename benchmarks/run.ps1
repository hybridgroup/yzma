# run.ps1 runs the benchmarks of this machine and puts each result in
# benchmarks/windows.md. See benchmarks/README.md.
#
#   .\benchmarks\run.ps1 -Machine ryzen-9-7950x -Label "AMD Ryzen 9 7950X"

[CmdletBinding()]
param(
  [string]$Machine = "",
  [string]$Label = "",
  [string]$Backend = "all",
  [string]$Suite = "text,multimodal",
  [string]$LlamaCpp = "",
  [int]$NCtx = 0,
  [int]$Count = 5,
  [string]$BenchTime = "10s",
  [switch]$DryRun
)

$ErrorActionPreference = "Stop"
$root = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
Set-Location $root

$file = "benchmarks/windows.md"
$modelsDir = if ($env:MODELS_DIR) { $env:MODELS_DIR } else { Join-Path $HOME "models" }
if (-not $env:YZMA_LIB) { $env:YZMA_LIB = Join-Path $root "lib" }
if (-not $env:YZMA_BENCHMARK_MODEL) {
  $env:YZMA_BENCHMARK_MODEL = Join-Path $modelsDir "SmolLM-135M.Q2_K.gguf"
}
if (-not $env:YZMA_BENCHMARK_MMMODEL) {
  $env:YZMA_BENCHMARK_MMMODEL = Join-Path $modelsDir "Qwen3-VL-2B-Instruct.Q4_K_M.gguf"
}
if (-not $env:YZMA_BENCHMARK_MMPROJ) {
  $env:YZMA_BENCHMARK_MMPROJ = Join-Path $modelsDir "Qwen3-VL-2B-Instruct.mmproj-Q8_0.gguf"
}

if (-not (Test-Path $env:YZMA_BENCHMARK_MODEL)) {
  throw "no model at $($env:YZMA_BENCHMARK_MODEL), run make download-benchmark-models"
}

$install = Join-Path $env:YZMA_LIB "yzma-install.json"
if (-not $LlamaCpp -and (Test-Path $install)) {
  $LlamaCpp = (Get-Content $install -Raw | ConvertFrom-Json).upstream_tag
}

if (-not $Machine) {
  $Machine = if ($env:YZMA_BENCH_MACHINE) { $env:YZMA_BENCH_MACHINE } else { $env:COMPUTERNAME }
}
$Machine = ($Machine.ToLower() -replace "[^a-z0-9._-]", "-")

$yzmaVersion = (Select-String -Path "version.go" -Pattern 'currentVersion = "(.*)"').Matches[0].Groups[1].Value
$work = New-Item -ItemType Directory -Path (Join-Path $env:TEMP ("yzma-bench-" + [guid]::NewGuid()))
$suites = $Suite.Split(",") | ForEach-Object { $_.Trim() }
$wanted = $Backend.Split(",") | ForEach-Object { $_.Trim().ToLower() }

function Test-Wanted($name) {
  return ($Backend -eq "all") -or ($wanted -contains $name)
}

# Write-Lines writes a file with no byte order mark, which the tool expects.
function Write-Lines($path, $lines) {
  [System.IO.File]::WriteAllLines($path, [string[]]($lines | ForEach-Object { "$_" }))
}

# Get-DeviceInfo collects what the machine says about the device of a backend.
function Get-DeviceInfo($backend, $path) {
  try {
    switch ($backend) {
      "cuda" { Write-Lines $path (nvidia-smi) }
      "vulkan" { Write-Lines $path (vulkaninfo --summary) }
      default { Write-Lines $path @() }
    }
  } catch {
    Write-Lines $path @()
  }
}

function Update-Section($suiteName, $backend, $device, $output, $info) {
  $argv = @("run", "./cmd/yzma-bench", "update", "--file", $file, "--suite", $suiteName,
    "--backend", $backend, "--device", $device, "--machine", $Machine, "--label", $Label,
    "--llamacpp", $LlamaCpp, "--yzma", $yzmaVersion, "--output", $output)
  if ((Test-Path $info) -and (Get-Item $info).Length -gt 0) { $argv += @("--device-info", $info) }
  if ($DryRun) { $argv += "--dry-run" }
  & go @argv
  if ($LASTEXITCODE -ne 0) { throw "yzma-bench failed" }
}

# Invoke-Suites runs both suites for one device. The benchmarks take -nctx and
# -device, which go test only gives to the package that it runs in.
function Invoke-Suites($backend, $device, $record, $info) {
  $ctx = $NCtx
  if ($ctx -eq 0) { $ctx = if ($backend -eq "cpu") { 8192 } else { 32000 } }

  foreach ($suiteName in $suites) {
    switch ($suiteName) {
      "text" { $pkg = "pkg/llama"; $bench = "BenchmarkInference" }
      "multimodal" { $pkg = "pkg/mtmd"; $bench = "BenchmarkMultimodalInference" }
      default { throw "unknown suite: $suiteName" }
    }
    if ($suiteName -eq "multimodal" -and -not (Test-Path $env:YZMA_BENCHMARK_MMMODEL)) {
      Write-Warning "no multimodal model at $($env:YZMA_BENCHMARK_MMMODEL), the suite is not run"
      continue
    }

    $out = Join-Path $work "$suiteName-$backend-$device.txt"
    Write-Host "==> $suiteName, $backend ($device)"

    $flag = "-device=$device"
    $command = "> cd $pkg; go test -benchtime=$BenchTime -count=$Count -run=nada -bench $bench -nctx=$ctx $flag"
    Push-Location (Join-Path $root $pkg)
    $result = & go test -benchtime=$BenchTime -count=$Count -run=nada -bench $bench -nctx=$ctx $flag 2>&1
    Pop-Location

    Write-Lines $out (@($command) + $result)
    $result | Write-Host

    Update-Section $suiteName $backend $record $out $info
  }
}

# The device list of llama.cpp says which backends this machine has.
$devices = Join-Path $work "devices.txt"
Write-Lines $devices (& go run . system -lib $env:YZMA_LIB)

$lines = Get-Content $devices
if (-not $Label) {
  for ($i = 0; $i -lt $lines.Count; $i++) {
    if ($lines[$i] -match "Backend:\s*CPU") {
      $description = $lines | Select-Object -Skip $i | Where-Object { $_ -match "Description:" } | Select-Object -First 1
      if ($description) { $Label = ($description -replace ".*Description:\s*", "") }
      break
    }
  }
  if (-not $Label) { $Label = $Machine }
}

$device = ""
$ranCPU = $false
foreach ($line in $lines) {
  if ($line -match "^Device \d+:\s*(.*)$") { $device = $Matches[1].Trim() }
  elseif ($line -match "Backend:\s*(.*)$") {
    $backend = $Matches[1].Trim().ToLower()
    $info = Join-Path $work "$backend-$device.txt"
    if ($backend -eq "cpu") {
      if ((Test-Wanted "cpu") -and -not $ranCPU) {
        Write-Lines $info @()
        Invoke-Suites "cpu" $device "" $info
        $ranCPU = $true
      }
    } elseif (Test-Wanted $backend) {
      Get-DeviceInfo $backend $info
      Invoke-Suites $backend $device $device $info
    }
  }
}
