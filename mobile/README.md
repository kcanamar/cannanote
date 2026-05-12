# CannaNote Mobile

Flutter mobile application for iOS and Android.

## Architecture

Hexagonal architecture with three layers:

```
lib/
├── core/                 # Business logic (no I/O)
│   ├── domain/           # Entities (Session, Product)
│   ├── ports/            # Interfaces (StoragePort)
│   └── application/      # Services (SessionService)
├── adapters/             # I/O implementations
│   └── storage/          # MemoryAdapter, DriftAdapter
├── di/                   # GetIt dependency injection
└── ui/                   # Flutter widgets
    ├── screens/
    └── widgets/
```

## Development

```bash
# Get dependencies
make deps

# Run in development mode
make dev

# Build release APK
make build-android

# Build release iOS
make build-ios

# Run tests
make test

# Generate code (after modifying Drift tables)
make gen

# Analyze code
make lint
```

## Dependencies

- **get_it** - Service locator for dependency injection
- **drift** - SQLite ORM with type-safe queries
- **uuid** - Unique identifiers
- **http** - REST API client

## Dependency Injection

Using GetIt as a simple service locator:

```dart
// Register at startup (di/injection.dart)
getIt.registerLazySingleton<StoragePort>(
  () => MemoryStorageAdapter(),
);

// Use anywhere
final storage = getIt<StoragePort>();
```

This approach:
- No Flutter coupling (works in tests, services)
- Clean hexagonal architecture
- Easy to swap implementations
