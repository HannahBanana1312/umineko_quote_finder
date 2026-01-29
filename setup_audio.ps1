$AudioDir = "internal\quote\data\audio"
$ZipUrl = "https://waifuvault.moe/f/da75978f-6ba4-474c-b063-f3f77a249470/voice.zip"
$TmpZip = "$env:TEMP\voice.zip"
$TmpDir = "$env:TEMP\voice"

if (Test-Path $AudioDir) {
    Write-Output "Audio directory already exists at $AudioDir, skipping download."
    exit 0
}

Write-Output "Downloading voice files..."
Invoke-WebRequest -Uri $ZipUrl -OutFile $TmpZip

Write-Output "Extracting..."
New-Item -ItemType Directory -Force -Path "internal\quote\data" | Out-Null
Expand-Archive -Path $TmpZip -DestinationPath $TmpDir
Move-Item -Path "$TmpDir\voice" -Destination $AudioDir

Remove-Item -Recurse -Force $TmpZip, $TmpDir
Write-Output "Done. Audio files extracted to $AudioDir"
