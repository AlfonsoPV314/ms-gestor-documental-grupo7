$raiz=Split-Path $PSScriptRoot
$rutas=@("$raiz/.tools/node-v22.22.0-win-x64","$raiz/.tools/pnpm/node_modules/.bin","$raiz/.tools/gitleaks","$env:LOCALAPPDATA/Programs/DockerDesktop/resources/bin")
$env:PATH=($rutas -join ';')+';'+$env:PATH
$env:GOTOOLCHAIN='go1.26.11'

