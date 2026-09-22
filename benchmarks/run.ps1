# run.ps1 runs the benchmarks of this machine and puts each result in
# benchmarks/windows.md. See benchmarks/README.md.
#
#   .\benchmarks\run.ps1 -Machine ryzen-9-7950x -Label "AMD Ryzen 9 7950X"
#
# PowerShell does not run a script until the policy of the machine permits it.
# Start the script from a PowerShell prompt in the directory of the repository.
#
#   powershell -ExecutionPolicy Bypass -File .\benchmarks\run.ps1

[CmdletBinding()]
param(
  [string]$Machine = "",
  [string]$Label = "",
  [string]$Backend = "all",
  [string]$Suite = "text,multimodal",
  [string]$LlamaCpp = "",
  [int]$NCtx = 0,
  [int]$Threads = 0,
  [switch]$Threadpool,
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
  $env:YZMA_BENCHMARK_MMMODEL = Join-Path $modelsDir "SmolVLM-256M-Instruct-Q8_0.gguf"
}
if (-not $env:YZMA_BENCHMARK_MMPROJ) {
  $env:YZMA_BENCHMARK_MMPROJ = Join-Path $modelsDir "mmproj-SmolVLM-256M-Instruct-Q8_0.gguf"
}

if (-not (Test-Path $env:YZMA_BENCHMARK_MODEL)) {
  throw "no model at $($env:YZMA_BENCHMARK_MODEL), run make download-benchmark-models"
}

# A nightly build has no upstream_tag, because its own tag names the assets.
$install = Join-Path $env:YZMA_LIB "yzma-install.json"
if (-not $LlamaCpp -and (Test-Path $install)) {
  $record = Get-Content $install -Raw | ConvertFrom-Json
  $LlamaCpp = if ($record.upstream_tag) { $record.upstream_tag } else { $record.tag }
}
if (-not $LlamaCpp) {
  throw "no tag of the llama.cpp build, run make download-llama.cpp or give -LlamaCpp"
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
# An empty list gives an empty file, because a pipeline with no output is null.
function Write-Lines($path, $lines) {
  $text = [string[]]@(@($lines) | Where-Object { $null -ne $_ } | ForEach-Object { "$_" })
  [System.IO.File]::WriteAllLines($path, $text)
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
    "--backend", $backend, "--machine", $Machine, "--label", $Label,
    "--llamacpp", $LlamaCpp, "--yzma", $yzmaVersion, "--output", $output)
  # PowerShell removes an empty argument, thus --device would take the name of
  # the next flag as its value.
  if ($device) { $argv += @("--device", $device) }
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

    # Each argument is one string, because PowerShell can divide a bare
    # argument that has a dash and a variable.
    $argv = @("test", "-benchtime=$BenchTime", "-count=$Count", "-run=nada",
      "-bench", "$bench", "-nctx=$ctx")
    if ($Threads -gt 0) { $argv += "-threads=$Threads" }
    if ($Threadpool -and $backend -eq "cpu") { $argv += "-threadpool" }
    if ($device) { $argv += "-device=$device" }
    $command = "> cd $pkg; go " + ($argv -join " ")
    Write-Host $command
    # Show each line when it comes. A suite takes many minutes.
    Push-Location (Join-Path $root $pkg)
    & go @argv 2>&1 | Tee-Object -Variable result | Out-Host
    Pop-Location

    Write-Lines $out (@($command) + $result)

    Update-Section $suiteName $backend $record $out $info
  }
}

# The device list of llama.cpp says which backends this machine has. The first
# go command builds the packages, thus it is slow.
Write-Host "==> the devices of this machine"
$devices = Join-Path $work "devices.txt"
Write-Lines $devices (& go run . system -lib $env:YZMA_LIB)
Get-Content $devices | Write-Host

$lines = Get-Content $devices
if (-not $Label) {
  for ($i = 0; $i -lt $lines.Count; $i++) {
    if ($lines[$i] -match "Backend:\s*CPU") {
      $description = $lines | Select-Object -Skip $i | Where-Object { $_ -match "Description:" } | Select-Object -First 1
      if ($description) { $Label = ($description -replace ".*Description:\s*", "") }
      # Some builds give "CPU" as the description, which is no name for a
      # machine.
      if ($Label -eq "CPU") { $Label = "" }
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
    # A name of a variable has no case in PowerShell. The name $backend here
    # would write on the $Backend parameter and stop each device.
    $devBackend = $Matches[1].Trim().ToLower()
    $info = Join-Path $work "$devBackend-$device.txt"
    if ($devBackend -eq "cpu") {
      if ((Test-Wanted "cpu") -and -not $ranCPU) {
        Write-Lines $info @()
        Invoke-Suites "cpu" $device "" $info
        $ranCPU = $true
      }
    } elseif (Test-Wanted $devBackend) {
      Get-DeviceInfo $devBackend $info
      Invoke-Suites $devBackend $device $device $info
    }
  }
}
