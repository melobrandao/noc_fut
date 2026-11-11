$ErrorActionPreference = 'Stop'
Push-Location $PSScriptRoot
try {
    if (Test-Path -LiteralPath '.\noc_fut.exe') { & '.\noc_fut.exe' }
    else { go run ./cmd/noc_fut }
} finally { Pop-Location }
