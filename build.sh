#!/bin/bash
# build.sh
echo "Building CSV Reader Application..."

# Set environment variables
export GOOS=linux  # or darwin for Mac
export GOARCH=amd64
export CGO_ENABLED=0

# Build the application
go build -o csv-reader

echo "Build complete! csv-reader has been created."