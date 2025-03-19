#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Create build directory
BUILD_DIR="build"
mkdir -p "$BUILD_DIR"

# Print with color
print_status() {
    echo -e "${GREEN}[BUILD]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Clean previous builds
print_status "Cleaning previous builds..."
rm -rf "$BUILD_DIR"/*

# Build for current platform
print_status "Building mast for current platform..."
go build -o "$BUILD_DIR/mast"

# Check if build was successful
if [ $? -eq 0 ]; then
    print_status "Build successful! Binary created as '$BUILD_DIR/mast'"
    # Make the binary executable
    chmod +x "$BUILD_DIR/mast"
else
    print_error "Build failed!"
    exit 1
fi

# Build for other platforms
print_status "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o "$BUILD_DIR/mast-linux-amd64"

print_status "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o "$BUILD_DIR/mast-windows-amd64.exe"

# Create release archives
print_status "Creating release archives..."
cd "$BUILD_DIR"
tar -czf mast-linux-amd64.tar.gz mast-linux-amd64
zip mast-windows-amd64.zip mast-windows-amd64.exe
cd ..

# Check if archives were created
if [ -f "$BUILD_DIR/mast-linux-amd64.tar.gz" ] && [ -f "$BUILD_DIR/mast-windows-amd64.zip" ]; then
    print_status "Release archives created successfully!"
    print_status "Linux: $BUILD_DIR/mast-linux-amd64.tar.gz"
    print_status "Windows: $BUILD_DIR/mast-windows-amd64.zip"
else
    print_error "Failed to create release archives!"
    exit 1
fi

# Run tests
print_status "Running tests..."
go test ./...

# Print version
print_status "Current version:"
"$BUILD_DIR/mast" version

print_status "Build process completed successfully!"