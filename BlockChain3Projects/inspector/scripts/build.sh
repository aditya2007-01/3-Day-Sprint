#!/bin/bash
set -e

echo "🏗️  Building BHIV Inspector for multiple platforms..."

# Create dist directory
mkdir -p dist

# Build for Linux
echo "📦 Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/inspector-linux ./cmd

# Build for Windows
echo "📦 Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/inspector.exe ./cmd

# Build for macOS
echo "📦 Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/inspector-mac ./cmd

# Build for macOS Intel
echo "📦 Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/inspector-mac-intel ./cmd

echo ""
echo "✅ Build complete! Artifacts:"
ls -lh dist/
echo ""
echo "📊 File sizes:"
du -sh dist/*
