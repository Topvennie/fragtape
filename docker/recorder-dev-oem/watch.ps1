$src = "C:\Users\Fragtape\Desktop\Shared\recorder\dist\recorder.exe"
$dstDir = "C:\Users\Fragtape\Desktop\Shared\recorder"
$dst = Join-Path $dstDir "recorder.exe"
$tmp = Join-Path $dstDir "recorder.exe.new"
$procName = "recorder"

New-Item -ItemType Directory -Force -Path $dstDir | Out-Null

function Stop-Recorder {
  $p = Get-Process $procName -ErrorAction SilentlyContinue
  if ($p) {
    $p | Stop-Process -Force -ErrorAction SilentlyContinue
    try { $p.WaitForExit() } catch {}
    Start-Sleep -Milliseconds 150
  }
}

function Wait-StableFile([string]$path, [int]$stableMs = 500) {
  $prevLen = -1
  $prevTime = [datetime]::MinValue
  $stableFor = 0

  while ($true) {
    try {
      $fi = Get-Item $path -ErrorAction Stop
      $len = $fi.Length
      $t = $fi.LastWriteTimeUtc

      if ($len -eq $prevLen -and $t -eq $prevTime) {
        $stableFor += 100
      } else {
        $stableFor = 0
        $prevLen = $len
        $prevTime = $t
      }

      if ($stableFor -ge $stableMs) { return $fi }
    } catch {
      $stableFor = 0
    }
    Start-Sleep -Milliseconds 100
  }
}

function Copy-And-Run {
  if (-not (Test-Path $src)) { return }

  $fi = Wait-StableFile $src 500
  Write-Host ("[{0}] Detected build: utc={1}, size={2}" -f (Get-Date), $fi.LastWriteTimeUtc, $fi.Length)

  Stop-Recorder

  Copy-Item -Force -ErrorAction Stop $src $tmp
  Move-Item -Force -ErrorAction Stop $tmp $dst

  Write-Host ("[{0}] Starting {1} (cwd {2})" -f (Get-Date), $dst, $dstDir)
  Start-Process -FilePath $dst -WorkingDirectory $dstDir | Out-Null
}

Write-Host ("Watching {0}..." -f $src)
Write-Host ("Copying to {0} and running from {1}" -f $dst, $dstDir)

Copy-And-Run

$lastTime = [datetime]::MinValue
$lastSize = -1

while ($true) {
  try {
    if (Test-Path $src) {
      $fi = Get-Item $src
      $t = $fi.LastWriteTimeUtc
      $s = $fi.Length

      if ($t -ne $lastTime -or $s -ne $lastSize) {
        $lastTime = $t
        $lastSize = $s
        Copy-And-Run
      }
    }
  } catch {
    Write-Host ("[{0}] Watch error: {1}" -f (Get-Date), $_.Exception.Message)
  }

  Start-Sleep -Milliseconds 250
}

