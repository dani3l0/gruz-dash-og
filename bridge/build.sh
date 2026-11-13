#!/bin/sh

# Disable CGO to avoid compatibility issues
export CGO_ENABLED=0

# Linux
GOOS=linux GOARCH=amd64 go build -o outgouge-linux-amd64
GOOS=linux GOARCH=arm64 go build -o outgouge-linux-arm64

# Windows
GOOS=windows GOARCH=amd64 go build -o outgouge-windows-amd64.exe

# macOS (Intel & Apple Silicon)
GOOS=darwin GOARCH=amd64 go build -o outgouge-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o outgouge-darwin-arm64
