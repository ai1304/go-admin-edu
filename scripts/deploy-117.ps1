param(
  [string]$HostName = "114.55.176.121",
  [string]$RemoteUser = "root",
  [string]$RemoteDir = "/opt/go-admin-edu",
  [string]$PasswordEnv = "GO_ADMIN_EDU_SSH_PASS",
  [switch]$SkipBuild,
  [switch]$SkipPortalDataCheck
)

$ErrorActionPreference = "Stop"

$Root = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$StateDir = Join-Path $Root ".local"
$RemoteScript = Join-Path $StateDir "deploy-114-remote.sh"
$ArtifactFile = Join-Path $Root "deploy-artifacts.tgz"

function Write-Info($Message) {
  Write-Host "[deploy-114] $Message" -ForegroundColor Cyan
}

function Test-Command($Name) {
  if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
    throw "Command '$Name' was not found. Please install it or add it to PATH."
  }
}

function Invoke-Step($Message, [scriptblock]$Block) {
  Write-Info $Message
  & $Block
}

function Build-Api {
  $apiOut = Join-Path $Root "deploy\artifacts\api\go-admin"
  $goCache = Join-Path $Root ".local\go-build-cache"
  New-Item -ItemType Directory -Force -Path (Split-Path $apiOut) | Out-Null
  New-Item -ItemType Directory -Force -Path $goCache | Out-Null

  $oldGoos = $env:GOOS
  $oldGoarch = $env:GOARCH
  $oldCgo = $env:CGO_ENABLED
  $oldGoCache = $env:GOCACHE
  try {
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    $env:CGO_ENABLED = "0"
    $env:GOCACHE = $goCache
    Push-Location (Join-Path $Root "go-admin")
    go build -trimpath -ldflags="-s -w" -o $apiOut .
    if ($LASTEXITCODE -ne 0) {
      throw "Backend build failed with exit code $LASTEXITCODE."
    }
  } finally {
    Pop-Location
    $env:GOOS = $oldGoos
    $env:GOARCH = $oldGoarch
    $env:CGO_ENABLED = $oldCgo
    $env:GOCACHE = $oldGoCache
  }
}

function Build-Web {
  npm.cmd --prefix (Join-Path $Root "web\apps\admin") run build
  if ($LASTEXITCODE -ne 0) {
    throw "Admin web build failed with exit code $LASTEXITCODE."
  }
  npm.cmd --prefix (Join-Path $Root "web\apps\portal") run build
  if ($LASTEXITCODE -ne 0) {
    throw "Portal web build failed with exit code $LASTEXITCODE."
  }
}

function Sync-Artifacts {
  $admin = Join-Path $Root "deploy\artifacts\admin"
  $portal = Join-Path $Root "deploy\artifacts\portal"
  $api = Join-Path $Root "deploy\artifacts\api"
  $apiBinary = Join-Path $api "go-admin"

  foreach ($path in @($admin, $portal, $api)) {
    $resolvedParent = Split-Path $path
    New-Item -ItemType Directory -Force -Path $path | Out-Null
    if (-not (Resolve-Path $resolvedParent).Path.StartsWith($Root)) {
      throw "Refusing to write outside workspace: $path"
    }
  }

  Remove-Item -LiteralPath (Join-Path $admin "*") -Recurse -Force -ErrorAction SilentlyContinue
  Remove-Item -LiteralPath (Join-Path $portal "*") -Recurse -Force -ErrorAction SilentlyContinue
  if (-not (Test-Path $apiBinary)) {
    throw "API artifact was not found: $apiBinary"
  }
  Copy-Item -Path (Join-Path $Root "web\apps\admin\dist\*") -Destination $admin -Recurse -Force
  Copy-Item -Path (Join-Path $Root "web\apps\portal\dist\*") -Destination $portal -Recurse -Force
  Copy-Item -Path (Join-Path $Root "deploy\api-runtime.Dockerfile") -Destination (Join-Path $api "Dockerfile") -Force
}

function Package-Artifacts {
  if (Test-Path $ArtifactFile) {
    Remove-Item -LiteralPath $ArtifactFile -Force
  }
  Push-Location $Root
  try {
    tar -czf deploy-artifacts.tgz -C deploy artifacts
    if ($LASTEXITCODE -ne 0) {
      throw "Packaging deploy artifacts failed with exit code $LASTEXITCODE."
    }
  } finally {
    Pop-Location
  }
  $item = Get-Item $ArtifactFile
  Write-Info ("Packaged {0:N1} MB: {1}" -f ($item.Length / 1MB), $item.FullName)
}

function Write-RemoteScript {
  New-Item -ItemType Directory -Force -Path $StateDir | Out-Null

  $remoteScriptContent = @'
#!/bin/sh
set -eu

COMPOSE="docker compose -p go-admin-edu -f deploy/docker-compose.runtime.yml"
SSH_OPTS="-o StrictHostKeyChecking=no -o UserKnownHostsFile=/tmp/known_hosts"

if [ -z "${SSH_PASS:-}" ]; then
  echo "SSH_PASS is required" >&2
  exit 1
fi

apk add --no-cache openssh-client sshpass >/dev/null

cd /workspace

echo "[remote] ensure directories"
sshpass -p "$SSH_PASS" ssh $SSH_OPTS "$REMOTE_USER@$REMOTE_HOST" "mkdir -p '$REMOTE_DIR/deploy' '$REMOTE_DIR/web/apps/admin' '$REMOTE_DIR/web/apps/portal' '$REMOTE_DIR/go-admin/config'"

echo "[remote] upload artifacts and runtime files"
sshpass -p "$SSH_PASS" scp $SSH_OPTS deploy-artifacts.tgz "$REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR/deploy-artifacts.tgz"
sshpass -p "$SSH_PASS" scp $SSH_OPTS deploy/docker-compose.runtime.yml "$REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR/deploy/docker-compose.runtime.yml"
sshpass -p "$SSH_PASS" scp $SSH_OPTS go-admin/config/settings.docker.yml "$REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR/go-admin/config/settings.docker.yml"
sshpass -p "$SSH_PASS" scp $SSH_OPTS web/apps/admin/nginx.conf "$REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR/web/apps/admin/nginx.conf"
sshpass -p "$SSH_PASS" scp $SSH_OPTS web/apps/portal/nginx.conf "$REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR/web/apps/portal/nginx.conf"

echo "[remote] unpack, migrate, restart"
sshpass -p "$SSH_PASS" ssh $SSH_OPTS "$REMOTE_USER@$REMOTE_HOST" "
set -eu
cd '$REMOTE_DIR'
mkdir -p deploy web/apps/admin web/apps/portal go-admin/config
sed -i 's#publicEndpoint: http://localhost:9000#publicEndpoint: http://'"${REMOTE_HOST}"':9000#' go-admin/config/settings.docker.yml || true
sed -i 's#publicEndpoint: http://[0-9.]*:9000#publicEndpoint: http://'"${REMOTE_HOST}"':9000#' go-admin/config/settings.docker.yml || true
rm -rf deploy/artifacts
tar -xzf deploy-artifacts.tgz -C deploy
chmod +x deploy/artifacts/api/go-admin
$COMPOSE up -d mysql redis minio
$COMPOSE run --rm api sh -c 'mkdir -p temp/logs && /main migrate -c /config/settings.yml'
$COMPOSE up -d api admin-web portal-web
$COMPOSE restart api admin-web portal-web
$COMPOSE ps
"

echo "[remote] health checks"
sshpass -p "$SSH_PASS" ssh $SSH_OPTS "$REMOTE_USER@$REMOTE_HOST" '
set -eu
curl -fsS http://127.0.0.1:9000/minio/health/live >/dev/null
api_code=$(curl -sS -o /tmp/go-admin-edu-api-check.out -w "%{http_code}" http://127.0.0.1:8000/api/v1/app-config || true)
admin_code=$(curl -sS -o /dev/null -w "%{http_code}" -I http://127.0.0.1:18080 || true)
portal_code=$(curl -sS -o /dev/null -w "%{http_code}" -I http://127.0.0.1:18081 || true)
echo "minio:200"
echo "api-app-config:${api_code}"
head -c 240 /tmp/go-admin-edu-api-check.out || true
echo
echo "admin:${admin_code}"
echo "portal:${portal_code}"
test "$api_code" = "200"
test "$admin_code" = "200"
test "$portal_code" = "200"
'

if [ "${SKIP_PORTAL_DATA_CHECK:-0}" != "1" ]; then
  echo "[remote] portal data checks"
  sshpass -p "$SSH_PASS" ssh $SSH_OPTS "$REMOTE_USER@$REMOTE_HOST" '
set -eu
for path in \
  /api/v1/portal/resources \
  /api/v1/portal/courses \
  /api/v1/portal/activities \
  /api/v1/portal/cases \
  /api/v1/portal/experts \
  /api/v1/portal/news
do
  code=$(curl -sS -o /tmp/go-admin-edu-portal-check.out -w "%{http_code}" "http://127.0.0.1:8000${path}?pageIndex=1&pageSize=1" || true)
  bytes=$(wc -c < /tmp/go-admin-edu-portal-check.out)
  echo "${path}:${code}:${bytes}"
  test "$code" = "200"
done
'
fi

echo "[remote] deploy ok"
'@

  [System.IO.File]::WriteAllText($RemoteScript, $remoteScriptContent, [System.Text.UTF8Encoding]::new($false))
}

$password = [Environment]::GetEnvironmentVariable($PasswordEnv)
if ([string]::IsNullOrWhiteSpace($password)) {
  throw "Environment variable '$PasswordEnv' is empty. Set it before deploying."
}

Test-Command "docker"
Test-Command "tar"

if (-not $SkipBuild) {
  Test-Command "go"
  Test-Command "npm.cmd"
  Invoke-Step "Build backend binary" { Build-Api }
  Invoke-Step "Build admin and portal web apps" { Build-Web }
} else {
  Write-Info "Skipping build; existing dist and API artifact will be used."
}

Invoke-Step "Sync frontend dist into deploy/artifacts" { Sync-Artifacts }
Invoke-Step "Package deploy artifacts" { Package-Artifacts }
Invoke-Step "Prepare remote deploy helper" { Write-RemoteScript }

$skipData = if ($SkipPortalDataCheck) { "1" } else { "0" }

Invoke-Step "Upload, migrate, restart, and check 114 server" {
  docker run --rm `
    -e SSH_PASS="$password" `
    -e REMOTE_HOST="$HostName" `
    -e REMOTE_USER="$RemoteUser" `
    -e REMOTE_DIR="$RemoteDir" `
    -e SKIP_PORTAL_DATA_CHECK="$skipData" `
    -v "${Root}:/workspace" `
    -v "${RemoteScript}:/remote_deploy.sh" `
    alpine:3.20 sh /remote_deploy.sh
  if ($LASTEXITCODE -ne 0) {
    throw "Remote deploy failed with exit code $LASTEXITCODE."
  }
}

Write-Host ""
Write-Info "Deploy completed."
Write-Host "Admin:  http://$HostName`:18080"
Write-Host "Portal: http://$HostName`:18081"
Write-Host "API:    http://$HostName`:8000/api/v1/app-config"
