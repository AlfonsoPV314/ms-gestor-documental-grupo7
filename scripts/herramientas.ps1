$ErrorActionPreference='Stop'
Set-Location (Split-Path $PSScriptRoot)
New-Item -ItemType Directory -Force '.tools' | Out-Null
$node='node-v22.22.0-win-x64.zip'
if(!(Test-Path '.tools/node-v22.22.0-win-x64/node.exe')){
 Invoke-WebRequest "https://nodejs.org/dist/v22.22.0/$node" -OutFile ".tools/$node"
 $checks=Invoke-RestMethod 'https://nodejs.org/dist/v22.22.0/SHASUMS256.txt'
 $esperado=($checks -split "`n" | Where-Object {$_ -match "\s+$([regex]::Escape($node))$"}).Split(' ',[StringSplitOptions]::RemoveEmptyEntries)[0]
 if((Get-FileHash ".tools/$node" -Algorithm SHA256).Hash.ToLower() -ne $esperado){throw 'Checksum Node inválido'}
 Expand-Archive -LiteralPath ".tools/$node" -DestinationPath '.tools' -Force
}
$env:PATH="$(Join-Path (Get-Location) '.tools/node-v22.22.0-win-x64');$env:PATH"
if(!(Test-Path '.tools/pnpm/node_modules/pnpm/bin/pnpm.cjs')){npm install --prefix .tools/pnpm --ignore-scripts --no-audit --no-fund pnpm@10.13.1;if($LASTEXITCODE -ne 0){throw 'Falló instalación pnpm'}}
if(!(Test-Path '.tools/gitleaks/gitleaks.exe')){
 $zip='gitleaks_8.30.1_windows_x64.zip'
 Invoke-WebRequest "https://github.com/gitleaks/gitleaks/releases/download/v8.30.1/$zip" -OutFile ".tools/$zip"
 $checks=Invoke-RestMethod 'https://github.com/gitleaks/gitleaks/releases/download/v8.30.1/gitleaks_8.30.1_checksums.txt'
 $esperado=($checks -split "`n" | Where-Object {$_ -match "\s+$([regex]::Escape($zip))$"}).Split(' ',[StringSplitOptions]::RemoveEmptyEntries)[0]
 if((Get-FileHash ".tools/$zip" -Algorithm SHA256).Hash.ToLower() -ne $esperado){throw 'Checksum Gitleaks inválido'}
 Expand-Archive -LiteralPath ".tools/$zip" -DestinationPath '.tools/gitleaks' -Force
}
Write-Host 'Herramientas locales listas. Ejecutar: . ./scripts/activar.ps1'
