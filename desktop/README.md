# CannaNote Desktop

Pure Rust desktop application using egui for UI.

## Architecture

Hexagonal architecture with three layers:

```
src/
├── core/                 # Business logic (no I/O)
│   ├── domain/           # Entities (Session, Product)
│   ├── ports/            # Interfaces (StoragePort)
│   └── application/      # Services (SessionService)
├── adapters/             # I/O implementations
│   └── storage/          # MemoryAdapter, SqliteAdapter
└── ui/                   # egui screens and widgets
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

# Lint with clippy
make lint
```

## Dependencies

- **eframe/egui** - Immediate mode GUI
- **tokio** - Async runtime
- **sqlx** - Database (SQLite)
- **uuid** - Unique identifiers
- **chrono** - Date/time handling

## Cross-Platform Build

Native binary compiles for macOS, Windows, and Linux from single codebase.

```bash
# macOS (native)
cargo build --release

# Windows (requires mingw)
cargo build --release --target x86_64-pc-windows-gnu

# Linux (requires musl or cross)
cargo build --release --target x86_64-unknown-linux-gnu
```

For CI builds, use GitHub Actions with native runners for each platform.
