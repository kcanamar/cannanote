// Core module - hexagonal architecture center
// Contains domain entities, ports, and application services

pub const domain = @import("domain/domain.zig");
pub const ports = @import("ports/ports.zig");
pub const application = @import("application/application.zig");

// Re-exports for convenience
pub const Session = domain.Session;
pub const SessionStatus = domain.SessionStatus;
pub const SyncStatus = domain.SyncStatus;

pub const StoragePort = ports.StoragePort;
pub const StorageError = ports.StorageError;

pub const SessionService = application.SessionService;
