param(
  [switch]$Infra,
  [switch]$Restart,
  [switch]$Stop,
  [switch]$Visible
)

$ErrorActionPreference = "Stop"

$Root = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$StateDir = Join-Path $Root ".local"
$LogDir = Join-Path $StateDir "logs"
$PidFile = Join-Path $StateDir "local-services.json"

New-Item -ItemType Directory -Force -Path $LogDir | Out-Null

function Write-Info($Message) {
  Write-Host "[local] $Message" -ForegroundColor Cyan
}

function Test-Command($Name) {
  if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
    throw "Command '$Name' was not found. Please install it or add it to PATH."
  }
}

function Read-State {
  if (-not (Test-Path $PidFile)) {
    return @()
  }

  $content = Get-Content -Raw -Path $PidFile
  if ([string]::IsNullOrWhiteSpace($content)) {
    return @()
  }

  $items = $content | ConvertFrom-Json
  return @($items | ForEach-Object { $_ })
}

function Save-State($Items) {
  @($Items) | ConvertTo-Json -Depth 4 | Set-Content -Encoding UTF8 -Path $PidFile
}

function Stop-LocalServices {
  $items = Read-State
  if ($items.Count -eq 0) {
    Write-Info "No recorded local services to stop."
    return
  }

  foreach ($item in $items) {
    $proc = Get-Process -Id $item.pid -ErrorAction SilentlyContinue
    if ($proc) {
      Write-Info "Stopping $($item.name) (PID $($item.pid))"
      taskkill.exe /PID $item.pid /T /F | Out-Null
    } else {
      Write-Info "$($item.name) is not running (PID $($item.pid))."
    }
  }

  Remove-Item -LiteralPath $PidFile -Force -ErrorAction SilentlyContinue
}

function Join-CommandLine($Executable, $Arguments) {
  $parts = @($Executable) + @($Arguments)
  return ($parts | ForEach-Object {
    if ($_ -match '[\s"]') {
      '"' + ($_ -replace '"', '\"') + '"'
    } else {
      $_
    }
  }) -join " "
}

function Quote-PSString($Value) {
  return "'" + ($Value -replace "'", "''") + "'"
}

function Join-PSArray($Values) {
  return "@(" + (($Values | ForEach-Object { Quote-PSString $_ }) -join ", ") + ")"
}

function Start-LocalService($Name, $WorkingDirectory, $Executable, $Arguments, $Port) {
  $stdout = Join-Path $LogDir "$Name.out.log"
  $stderr = Join-Path $LogDir "$Name.err.log"
  $runner = Join-Path $StateDir "$Name.runner.ps1"

  Remove-Item -LiteralPath $stdout, $stderr -Force -ErrorAction SilentlyContinue

  $commandLine = Join-CommandLine -Executable $Executable -Arguments $Arguments
  $psWorkingDirectory = Quote-PSString $WorkingDirectory
  $psExecutable = Quote-PSString $Executable
  $psArguments = Join-PSArray $Arguments
  $psStdout = Quote-PSString $stdout
  $psStderr = Quote-PSString $stderr

  if ($Visible) {
    @"
Set-Location -LiteralPath $psWorkingDirectory
`$serviceArgs = $psArguments
& $psExecutable @serviceArgs
"@ | Set-Content -Encoding UTF8 -Path $runner

    $process = Start-Process `
      -FilePath "powershell.exe" `
      -ArgumentList "-NoExit -NoProfile -ExecutionPolicy Bypass -File `"$runner`"" `
      -WorkingDirectory $WorkingDirectory `
      -PassThru
  } else {
    @"
Set-Location -LiteralPath $psWorkingDirectory
`$serviceArgs = $psArguments
& $psExecutable @serviceArgs 1> $psStdout 2> $psStderr
"@ | Set-Content -Encoding UTF8 -Path $runner

    $process = Start-Process `
      -FilePath "powershell.exe" `
      -ArgumentList "-NoProfile -ExecutionPolicy Bypass -File `"$runner`"" `
      -WorkingDirectory $WorkingDirectory `
      -PassThru `
      -WindowStyle Hidden
  }

  [pscustomobject]@{
    name = $Name
    pid = $process.Id
    port = $Port
    cwd = $WorkingDirectory
    command = $commandLine
    stdout = $stdout
    stderr = $stderr
    runner = $runner
  }
}

if ($Stop) {
  Stop-LocalServices
  return
}

if ($Restart) {
  Stop-LocalServices
}

Test-Command "go"
Test-Command "npm.cmd"

if ($Infra) {
  Test-Command "docker"
  Write-Info "Starting infra containers: mysql, redis, minio"
  docker compose -f (Join-Path $Root "docker-compose.yml") up -d mysql redis minio
}

$services = @(
  @{
    Name = "backend"
    WorkingDirectory = Join-Path $Root "go-admin"
    Executable = "go"
    Arguments = @("run", ".", "server", "-c", "config/settings.yml", "-a", "true")
    Port = 8000
  },
  @{
    Name = "admin-web"
    WorkingDirectory = Join-Path $Root "web\apps\admin"
    Executable = "npm.cmd"
    Arguments = @("run", "dev")
    Port = 1798
  },
  @{
    Name = "portal-web"
    WorkingDirectory = Join-Path $Root "web\apps\portal"
    Executable = "npm.cmd"
    Arguments = @("run", "dev")
    Port = 1799
  }
)

$state = @()
foreach ($service in $services) {
  Write-Info "Starting $($service.Name) on port $($service.Port)"
  $state += Start-LocalService `
    -Name $service.Name `
    -WorkingDirectory $service.WorkingDirectory `
    -Executable $service.Executable `
    -Arguments $service.Arguments `
    -Port $service.Port
}

Save-State $state

Write-Host ""
Write-Info "All local services have been started."
Write-Host "Backend:     http://localhost:8000"
Write-Host "Admin web:   http://localhost:1798"
Write-Host "Portal web:  http://localhost:1799"
Write-Host "Logs:        $LogDir"
Write-Host "Stop:        powershell -ExecutionPolicy Bypass -File scripts\start-local.ps1 -Stop"
Write-Host ""
Write-Host "Tips:"
Write-Host "  Start infra too:  powershell -ExecutionPolicy Bypass -File scripts\start-local.ps1 -Infra"
Write-Host "  Restart all:      powershell -ExecutionPolicy Bypass -File scripts\start-local.ps1 -Restart"
Write-Host "  Show windows:     powershell -ExecutionPolicy Bypass -File scripts\start-local.ps1 -Visible"
