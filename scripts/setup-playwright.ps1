$ErrorActionPreference = 'Stop'
Push-Location (Split-Path $PSScriptRoot -Parent)
try {
    # Instala o driver e o Chromium headless para validar a API e a interface.
    go run github.com/mxschmitt/playwright-go/cmd/playwright@v0.6401.0 install --only-shell chromium
    if ($LASTEXITCODE -ne 0) { throw 'Falha ao instalar Playwright.' }
} finally { Pop-Location }
