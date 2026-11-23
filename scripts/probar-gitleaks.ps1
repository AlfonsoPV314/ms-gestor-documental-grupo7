$ErrorActionPreference='Stop'
Set-Location (Split-Path $PSScriptRoot)
. ./scripts/activar.ps1
$ruta=Join-Path (Get-Location) ('.artifacts/gitleaks-'+[guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Force $ruta | Out-Null
Push-Location $ruta
try{
 git init -q
 git config user.name 'Prueba local'
 git config user.email 'prueba@example.invalid'
 $aws='AKIA'+(-join (1..16 | ForEach-Object {[char](Get-Random -Minimum 65 -Maximum 91)}))
 "aws_access_key_id=$aws" | Set-Content -Encoding utf8NoBOM 'secreto-sintetico.txt'
 git add secreto-sintetico.txt
 $hook="#!/bin/sh`nexec gitleaks protect --staged --redact --no-banner`n"
 [IO.File]::WriteAllText((Join-Path $ruta '.git/hooks/pre-commit'),$hook,[Text.UTF8Encoding]::new($false))
 $salida=git commit -m 'test: secreto sintético que debe rechazarse' 2>&1
 $codigo=$LASTEXITCODE
 $salida | Set-Content -Encoding utf8NoBOM (Join-Path $ruta 'resultado.log')
 if($codigo -eq 0){throw 'El hook dejó pasar el secreto sintético'}
 if(($salida -join "`n") -notmatch 'leaks found'){throw 'El commit falló por una causa distinta a Gitleaks'}
 Write-Host 'Correcto: pre-commit bloqueó un secreto sintético. Log redactado en .artifacts.'
}finally{Pop-Location}
