// Adapters module - hexagonal architecture periphery
// Contains concrete implementations of ports

pub const storage = @import("storage/storage.zig");

pub const MemoryStorageAdapter = storage.MemoryStorageAdapter;
pub const SqliteStorageAdapter = storage.SqliteStorageAdapter;
