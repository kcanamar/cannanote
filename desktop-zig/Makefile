.PHONY: dev build test clean fmt

# Run in development mode
dev:
	zig build run

# Build release binary
build:
	zig build -Doptimize=ReleaseSafe

# Run tests
test:
	zig build test

# Clean build artifacts
clean:
	rm -rf zig-out .zig-cache

# Format code
fmt:
	zig fmt src/

# Fetch dependencies
fetch:
	zig fetch --save 'git+https://github.com/david-vanderson/dvui#main'

# Build for all platforms (requires cross-compilation setup)
build-all: build-macos build-windows build-linux

build-macos:
	zig build -Doptimize=ReleaseSafe

build-windows:
	zig build -Doptimize=ReleaseSafe -Dtarget=x86_64-windows

build-linux:
	zig build -Doptimize=ReleaseSafe -Dtarget=x86_64-linux
