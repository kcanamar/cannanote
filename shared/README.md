# Shared Definitions

Future home of Protocol Buffer definitions for cross-platform model generation.

## Current Status

**Deferred** - Manual domain models in each platform for now.

Protocol Buffers will be added when:
- Desktop and mobile have stable domain models
- Sync API is production-ready
- Maintaining 3+ manual definitions becomes painful

## Planned Structure

```
shared/
├── proto/                  # Protocol Buffer definitions
│   ├── session.proto       # Session entity
│   ├── product.proto       # Product entity
│   ├── human.proto         # User entity
│   └── sync.proto          # Sync protocol messages
├── generated/              # Generated code (gitignored)
│   ├── rust/
│   ├── dart/
│   └── go/
└── scripts/
    └── generate-models.sh  # Run protoc for all platforms
```

## Why Protocol Buffers?

Single source of truth for data models across:
- **Go** (backend API)
- **Rust** (desktop app)
- **Dart** (mobile app)

Benefits:
- Change schema once, regenerate everywhere
- Compile-time errors catch mismatches
- Built-in serialization for sync protocol
- Language-agnostic documentation

## When to Add

Triggers for adding protobufs:
1. Adding 5+ entities with cross-platform sync
2. Sync API design is finalized
3. Manual model drift causes bugs
4. Team scales beyond 1 developer

Until then, manual models provide:
- Faster iteration
- Less tooling overhead
- More learning opportunities
