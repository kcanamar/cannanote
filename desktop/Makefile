.PHONY: dev build test clean fmt lint

# Development
dev:
	cargo run

# Build release binary
build:
	cargo build --release

# Run tests
test:
	cargo test

# Clean build artifacts
clean:
	cargo clean

# Format code
fmt:
	cargo fmt

# Lint code
lint:
	cargo clippy

# Build for all platforms (requires cross-compilation setup)
build-all: build-macos build-windows build-linux

build-macos:
	cargo build --release

build-windows:
	cargo build --release --target x86_64-pc-windows-gnu

build-linux:
	cargo build --release --target x86_64-unknown-linux-gnu
