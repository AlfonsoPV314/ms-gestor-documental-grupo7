param([switch]$Mutacion)
$ErrorActionPreference='Stop'
Set-Location (Split-Path $PSScriptRoot)
. ./scripts/activar.ps1
docker compose up -d --wait postgres
if($LASTEXITCODE -ne 0){throw 'PostgreSQL no disponible'}
$db=(Get-Content .env | Where-Object {$_ -like 'POSTGRES_PASSWORD=*'}).Split('=',2)[1]
$env:TEST_DATABASE_URL="postgres://documental:${db}@host.docker.internal:5173/documental?sslmode=disable"
docker build -f apps/api/Dockerfile.test -t documental-tests:local .
if($LASTEXITCODE -ne 0){throw 'Falló runner E2E'}
$comando='scripts/e2e.mjs'
if($Mutacion){$comando='scripts/mutacion.mjs'}
docker run --rm --user 0 --add-host host.docker.internal:host-gateway -e TEST_DATABASE_URL -v "${PWD}:/src" -v documental-go-mod:/go/pkg/mod -v documental-go-cache:/root/.cache/go-build documental-tests:local node $comando
if($LASTEXITCODE -ne 0){throw 'Falló validación E2E'}
