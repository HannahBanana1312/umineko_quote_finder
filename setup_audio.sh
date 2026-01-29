#!/bin/sh
set -e

AUDIO_DIR="internal/quote/data/audio"
ZIP_URL="https://waifuvault.moe/f/da75978f-6ba4-474c-b063-f3f77a249470/voice.zip"

if [ -d "$AUDIO_DIR" ]; then
    echo "Audio directory already exists at $AUDIO_DIR, skipping download."
    exit 0
fi

echo "Downloading voice files..."
curl -fSL -o /tmp/voice.zip "$ZIP_URL"

echo "Extracting..."
mkdir -p internal/quote/data
unzip -qo /tmp/voice.zip -d /tmp/voice
mv /tmp/voice/voice "$AUDIO_DIR"

rm -rf /tmp/voice.zip /tmp/voice
echo "Done. Audio files extracted to $AUDIO_DIR"
