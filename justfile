default:
    just --list

# Build for linux-arm64
build:
    @echo "Building for linux-arm64..."
    GOOS=linux GOARCH=arm64 go build -o build/tempapi

