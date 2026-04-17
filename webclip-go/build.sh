#!/bin/bash

# Configuration
APP_NAME="webclip"
BUILD_DIR="bin"
GO_BIN="/usr/local/go/bin/go"

# Clean previous builds
rm -rf $BUILD_DIR
mkdir -p $BUILD_DIR

# List of platforms to build for
# Format: GOOS/GOARCH
platforms=(
    "windows/amd64"
    "windows/arm64"
    "linux/amd64"
    "linux/arm"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

echo "Starting cross-compilation..."

for platform in "${platforms[@]}"
do
    # Split platform into OS and ARCH
    platform_split=(${platform//\// })
    export GOOS=${platform_split[0]}
    export GOARCH=${platform_split[1]}
    
    # Determine extension
    extension=""
    if [ $GOOS = "windows" ]; then
        extension=".exe"
    fi
    
    output_name="${APP_NAME}-${GOOS}-${GOARCH}${extension}"
    echo "Building: $output_name"
    
    $GO_BIN build -o "${BUILD_DIR}/${output_name}" .
    
    if [ $? -ne 0 ]; then
        echo "Error: Failed to build for $platform"
    fi
done

echo "Done! Binaries are in the '${BUILD_DIR}' folder."
ls -lh $BUILD_DIR
