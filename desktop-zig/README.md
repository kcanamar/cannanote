# CannaNote Desktop (Zig)

Pure Zig desktop application using dvui for UI.

## Philosophy

- **Zero C/C++ dependencies** - Pure Zig supply chain
- **First principles** - Minimal abstractions, explicit code
- **Hexagonal architecture** - Clean separation of concerns

## Architecture

```
src/
├── main.zig              # dvui app entry
├── core/                 # Business logic (no I/O)
│   ├── domain/           # Entities (Session)
│   ├── ports/            # Interfaces (StoragePort)
│   └── application/      # Services (SessionService)
└── adapters/             # I/O implementations
    └── storage/          # MemoryAdapter, (future: SQLite)
```

## Requirements

- Zig 0.16.0+
- SDL3 (for dvui backend)

### Install SDL3 on macOS

```bash
brew install sdl3
```

## Development

```bash
# Run in development mode
make dev

# Build release binary
make build

# Run tests
make test

# Format code
make fmt

# Clean build artifacts
make clean
```

## Cross-Platform Build

Zig makes cross-compilation simple:

```bash
# macOS (native)
zig build -Doptimize=ReleaseSafe

# Windows
zig build -Doptimize=ReleaseSafe -Dtarget=x86_64-windows

# Linux
zig build -Doptimize=ReleaseSafe -Dtarget=x86_64-linux
```

## Dependencies

- **dvui** - Immediate mode GUI, pure Zig
- No other runtime dependencies

## Why Zig?

1. **Supply chain security** - No npm, pip, cargo dependencies to audit
2. **Simplicity** - No hidden control flow, no hidden allocations
3. **C interop** - Seamless when needed (SQLite future)
4. **Cross-compilation** - Build for any platform from any platform
5. **Philosophy** - Language design values align with project values
