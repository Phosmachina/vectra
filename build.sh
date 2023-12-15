#!/bin/bash

# cleanup previous builds
rm vectra*

# set version
VERSION="1.1.0"

# build for Linux
GOOS=linux GOARCH=amd64 go build -o "vectra-linux-amd64-$VERSION" -v ./app.go

# build for windows
GOOS=windows GOARCH=amd64 go build -o "vectra-win-amd64-$VERSION.exe" -v ./app.go

echo "Build complete: $(ls vectra*)"
