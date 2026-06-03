param(
  [string]$HostName = "114.55.176.121",
  [string]$RemoteUser = "root",
  [string]$RemoteDir = "/opt/go-admin-edu",
  [string]$PasswordEnv = "GO_ADMIN_EDU_SSH_PASS",
  [switch]$SkipBuild,
  [switch]$SkipPortalDataCheck
)

$params = @{
  HostName = $HostName
  RemoteUser = $RemoteUser
  RemoteDir = $RemoteDir
  PasswordEnv = $PasswordEnv
}

if ($SkipBuild) {
  $params.SkipBuild = $true
}
if ($SkipPortalDataCheck) {
  $params.SkipPortalDataCheck = $true
}

& (Join-Path $PSScriptRoot "deploy-117.ps1") @params
exit $LASTEXITCODE
