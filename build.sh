#!/bin/bash
#
# Build script for claude-code-logs
#

set -e

# Get version info
VERSION=${VERSION:-"dev"}
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

# Default values
OUTPUT="claude-code-logs"
GOOS=${GOOS:-$(go env GOOS)}
GOARCH=${GOARCH:-$(go env GOARCH)}

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -o, --output NAME    Output binary name (default: claude-code-logs)"
    echo "  -v, --version VER    Version string (default: dev)"
    echo "  --os OS              Target OS (default: current)"
    echo "  --arch ARCH          Target architecture (default: current)"
    echo "  --all                Build for all supported platforms"
    echo "  -h, --help           Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                   # Build for current platform"
    echo "  $0 -v 1.0.0          # Build with version 1.0.0"
    echo "  $0 --all             # Build for all platforms"
}

build() {
    local os=$1
    local arch=$2
    local output=$3

    echo "Building for ${os}/${arch}..."
    CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build -ldflags "${LDFLAGS}" -o "${output}" .
    echo "  -> ${output}"
}

# Parse arguments
BUILD_ALL=false
while [[ $# -gt 0 ]]; do
    case $1 in
        -o|--output)
            OUTPUT="$2"
            shift 2
            ;;
        -v|--version)
            VERSION="$2"
            LDFLAGS="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"
            shift 2
            ;;
        --os)
            GOOS="$2"
            shift 2
            ;;
        --arch)
            GOARCH="$2"
            shift 2
            ;;
        --all)
            BUILD_ALL=true
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

# Run go mod tidy first
echo "Running go mod tidy..."
go mod tidy

if [ "$BUILD_ALL" = true ]; then
    # Create dist directory
    mkdir -p dist

    # Build for all supported platforms (matching goreleaser config)
    build darwin amd64 "dist/claude-code-logs_darwin_amd64"
    build darwin arm64 "dist/claude-code-logs_darwin_arm64"
    build linux amd64 "dist/claude-code-logs_linux_amd64"

    echo ""
    echo "All builds complete! Binaries are in ./dist/"
else
    # Build for single platform
    build ${GOOS} ${GOARCH} "${OUTPUT}"
    echo ""
    echo "Build complete!"
fi
