const std = @import("std");
const Session = @import("../domain/domain.zig").Session;

/// Storage errors
pub const StorageError = error{
    NotFound,
    DatabaseError,
    InitFailed,
    OutOfMemory,
};

/// Port defining what storage must do - not HOW
/// This is the hexagonal architecture "port" - an interface
/// that defines capabilities without implementation details.
pub const StoragePort = struct {
    ptr: *anyopaque,
    vtable: *const VTable,

    pub const VTable = struct {
        /// Initialize storage (create tables, etc)
        init: *const fn (ptr: *anyopaque) StorageError!void,

        /// Add new session
        addSession: *const fn (ptr: *anyopaque, session: *const Session) StorageError!Session,

        /// Get session by ID
        getSession: *const fn (ptr: *anyopaque, id: []const u8) StorageError!?Session,

        /// Get all sessions
        getAllSessions: *const fn (ptr: *anyopaque, allocator: std.mem.Allocator) StorageError![]Session,

        /// Update existing session
        updateSession: *const fn (ptr: *anyopaque, session: *const Session) StorageError!Session,

        /// Delete session
        deleteSession: *const fn (ptr: *anyopaque, id: []const u8) StorageError!void,

        /// Get sessions pending sync
        getPendingSessions: *const fn (ptr: *anyopaque, allocator: std.mem.Allocator) StorageError![]Session,

        /// Storage type identifier
        storageType: *const fn (ptr: *anyopaque) []const u8,
    };

    /// Initialize storage
    pub fn init(self: StoragePort) StorageError!void {
        return self.vtable.init(self.ptr);
    }

    /// Add new session
    pub fn addSession(self: StoragePort, session: *const Session) StorageError!Session {
        return self.vtable.addSession(self.ptr, session);
    }

    /// Get session by ID
    pub fn getSession(self: StoragePort, id: []const u8) StorageError!?Session {
        return self.vtable.getSession(self.ptr, id);
    }

    /// Get all sessions
    pub fn getAllSessions(self: StoragePort, allocator: std.mem.Allocator) StorageError![]Session {
        return self.vtable.getAllSessions(self.ptr, allocator);
    }

    /// Update existing session
    pub fn updateSession(self: StoragePort, session: *const Session) StorageError!Session {
        return self.vtable.updateSession(self.ptr, session);
    }

    /// Delete session
    pub fn deleteSession(self: StoragePort, id: []const u8) StorageError!void {
        return self.vtable.deleteSession(self.ptr, id);
    }

    /// Get sessions pending sync
    pub fn getPendingSessions(self: StoragePort, allocator: std.mem.Allocator) StorageError![]Session {
        return self.vtable.getPendingSessions(self.ptr, allocator);
    }

    /// Storage type identifier
    pub fn storageType(self: StoragePort) []const u8 {
        return self.vtable.storageType(self.ptr);
    }
};
