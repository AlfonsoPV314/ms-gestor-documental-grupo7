param([switch]$Calidad)
$ErrorActionPreference='Stop'
Set-Location (Split-Path $PSScriptRoot)
$dockerDir=Join-Path $env:LOCALAPPDATA 'Programs/DockerDesktop/resources/bin'
if(Test-Path $dockerDir){$env:PATH="$dockerDir;$env:PATH"}
if(!(Test-Path '.env')){
 $db=[Convert]::ToHexString([Security.Cryptography.RandomNumberGenerator]::GetBytes(24))
 $sonar=[Convert]::ToHexString([Security.Cryptography.RandomNumberGenerator]::GetBytes(24))
 (Get-Content '.env.example' -Raw).Replace('replace-with-local-database-password',$db).Replace('replace-with-local-sonar-database-password',$sonar) | Set-Content -Encoding utf8NoBOM '.env'
}
if(!(Test-Path '.secrets/firma.key')){
 $env:SIGNING_KEY_FILE=(Join-Path (Get-Location) '.secrets/firma.key')
 Push-Location apps/api
 try{ $env:GOTOOLCHAIN='go1.26.11'; go run ./cmd/claves; if($LASTEXITCODE -ne 0){throw 'No se pudo generar la clave'} }finally{Pop-Location}
}
if($Calidad){docker compose -f compose.quality.yaml up -d; if($LASTEXITCODE -ne 0){throw 'Falló SonarQube'};Write-Host 'SonarQube: http://localhost:9000 (esperar estado UP)'}
else{docker compose up -d --build --wait; if($LASTEXITCODE -ne 0){throw 'Falló Compose'};Write-Host 'Aplicación: http://localhost:5170 · API: http://localhost:5171/api/docs'}


