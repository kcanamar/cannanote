// Ports module
// Interfaces defining external capabilities

pub const storage = @import("storage_port.zig");

pub const StoragePort = storage.StoragePort;
pub const StorageError = storage.StorageError;
