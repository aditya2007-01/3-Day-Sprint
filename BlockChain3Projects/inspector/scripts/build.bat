@echo off
echo Building BHIV Inspector for multiple platforms...

REM Create dist directory
if not exist dist mkdir dist

echo.
echo [1/4] Building for Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -o dist/inspector-linux ./cmd

echo [2/4] Building for Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o dist/inspector.exe ./cmd

echo [3/4] Building for macOS (arm64)...
set GOOS=darwin
set GOARCH=arm64
go build -ldflags="-s -w" -o dist/inspector-mac ./cmd

echo [4/4] Building for macOS (amd64)...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-s -w" -o dist/inspector-mac-intel ./cmd

echo.
echo Build complete! Artifacts:
dir dist
echo.
pause
