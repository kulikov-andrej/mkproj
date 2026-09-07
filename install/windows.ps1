$ErrorActionPreference = "Stop"

$Repo = "kulikov-andrej/mkproj"
$AssetName = "mkproj-windows-amd64.exe"

$InstallDir = Join-Path $env:LOCALAPPDATA "Programs\mkproj"
$InstallPath = Join-Path $InstallDir "mkproj.exe"

$ConfigDir = Join-Path $env:APPDATA "mkproj"
$TemplatesDir = Join-Path $ConfigDir "templates"

Write-Host "Installing mkproj..."

$Release = Invoke-RestMethod `
    -Uri "https://api.github.com/repos/$Repo/releases/latest" `
    -Headers @{
        Accept = "application/vnd.github+json"
    }

$Asset = $Release.assets |
    Where-Object { $_.name -eq $AssetName } |
    Select-Object -First 1

if (-not $Asset) {
    throw "Release asset not found: $AssetName"
}

if (-not $Asset.digest) {
    throw "Release asset has no digest"
}

if (-not $Asset.digest.StartsWith("sha256:")) {
    throw "Unsupported asset digest: $($Asset.digest)"
}

$ExpectedHash = $Asset.digest.Substring("sha256:".Length)

$TempFile = Join-Path `
    ([System.IO.Path]::GetTempPath()) `
    "mkproj-$([System.Guid]::NewGuid()).exe"

try {
    Write-Host "Downloading $($Release.tag_name)..."

    Invoke-WebRequest `
        -Uri $Asset.browser_download_url `
        -OutFile $TempFile

    $ActualHash = (
        Get-FileHash `
            -Path $TempFile `
            -Algorithm SHA256
    ).Hash.ToLowerInvariant()

    if ($ActualHash -ne $ExpectedHash.ToLowerInvariant()) {
        throw "SHA-256 verification failed"
    }

    New-Item `
        -ItemType Directory `
        -Force `
        -Path $InstallDir |
        Out-Null

    Move-Item `
        -Force `
        -Path $TempFile `
        -Destination $InstallPath

    $TempFile = $null

    $UserPath = [Environment]::GetEnvironmentVariable(
        "Path",
        "User"
    )

    $PathEntries = @(
        $UserPath -split ";" |
        Where-Object { $_ }
    )

    $AlreadyInPath = $PathEntries |
        Where-Object {
            $_.TrimEnd("\") -ieq $InstallDir.TrimEnd("\")
        }

    if (-not $AlreadyInPath) {
        $NewPath = if ($UserPath) {
            "$UserPath;$InstallDir"
        }
        else {
            $InstallDir
        }

        [Environment]::SetEnvironmentVariable(
            "Path",
            $NewPath,
            "User"
        )

        Write-Host "Added mkproj to user PATH."
    }

    Write-Host ""
    Write-Host "Installed:"
    Write-Host "  $InstallPath"

    if (Test-Path $TemplatesDir) {
        if (-not (Get-Item $TemplatesDir).PSIsContainer) {
            throw "Template path is not a directory: $TemplatesDir"
        }
    }
    else {
        New-Item `
            -ItemType Directory `
            -Path $TemplatesDir `
            -Force |
            Out-Null

        Write-Host "Created template directory:"
        Write-Host "  $TemplatesDir"
    }

    & $InstallPath --version

    if (-not $AlreadyInPath) {
        Write-Host ""
        Write-Host "Open a new terminal to use 'mkproj' from PATH."
    }
}
finally {
    if ($TempFile -and (Test-Path $TempFile)) {
        Remove-Item -Force $TempFile
    }
}