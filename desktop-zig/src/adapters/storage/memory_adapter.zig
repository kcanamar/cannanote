const std = @import("std");
const core = @import("../../core/core.zig");
const Session = core.Session;
const SyncStatus = core.SyncStatus;
const StoragePort = core.StoragePort;
const StorageError = core.StorageError;

/// Simple in-memory storage for testing and development
/// This adapter implements the StoragePort interface using
/// in-memory hash maps. Data persists only during app lifetime.
pub const MemoryStorageAdapter = struct {
    sessions: std.StringHashMap(Session),
    allocator: std.mem.Allocator,

    const Self = @This();

    pub fn init(allocator: std.mem.Allocator) Self {
        return Self{
            .sessions = std.StringHashMap(Session).init(allocator),
            .allocator = allocator,
        };
    }

    pub fn deinit(self: *Self) void {
        // Free all stored sessions
        var it = self.sessions.iterator();
        while (it.next()) |entry| {
            var session = entry.value_ptr.*;
            session.deinit(self.allocator);
        }
        self.sessions.deinit();
    }

    /// Get as StoragePort interface
    pub fn port(self: *Self) StoragePort {
        return StoragePort{
            .ptr = self,
            .vtable = &vtable,
        };
    }

    // VTable implementation
    const vtable = StoragePort.VTable{
        .init = initImpl,
        .addSession = addSessionImpl,
        .getSession = getSessionImpl,
        .getAllSessions = getAllSessionsImpl,
        .updateSession = updateSessionImpl,
        .deleteSession = deleteSessionImpl,
        .getPendingSessions = getPendingSessionsImpl,
        .storageType = storageTypeImpl,
    };

    fn initImpl(ptr: *anyopaque) StorageError!void {
        _ = ptr;
        // Nothing to initialize for memory storage
    }

    fn addSessionImpl(ptr: *anyopaque, session: *const Session) StorageError!Session {
        const self: *Self = @ptrCast(@alignCast(ptr));

        // Clone the session
        const id_copy = self.allocator.dupe(u8, session.getId()) catch return StorageError.OutOfMemory;

        self.sessions.put(id_copy, session.*) catch return StorageError.DatabaseError;

        return session.*;
    }

    fn getSessionImpl(ptr: *anyopaque, id: []const u8) StorageError!?Session {
        const self: *Self = @ptrCast(@alignCast(ptr));

        return self.sessions.get(id);
    }

    fn getAllSessionsImpl(ptr: *anyopaque, allocator: std.mem.Allocator) StorageError![]Session {
        const self: *Self = @ptrCast(@alignCast(ptr));

        var result: std.ArrayList(Session) = .empty;
        errdefer result.deinit(allocator);

        var it = self.sessions.iterator();
        while (it.next()) |entry| {
            result.append(allocator, entry.value_ptr.*) catch return StorageError.OutOfMemory;
        }

        // Sort by timestamp descending (newest first)
        const items = result.toOwnedSlice(allocator) catch return StorageError.OutOfMemory;
        std.mem.sort(Session, items, {}, struct {
            fn lessThan(_: void, a: Session, b: Session) bool {
                return a.timestamp > b.timestamp;
            }
        }.lessThan);

        return items;
    }

    fn updateSessionImpl(ptr: *anyopaque, session: *const Session) StorageError!Session {
        const self: *Self = @ptrCast(@alignCast(ptr));

        const key = session.getId();
        if (self.sessions.contains(key)) {
            self.sessions.put(key, session.*) catch return StorageError.DatabaseError;
            return session.*;
        }

        return StorageError.NotFound;
    }

    fn deleteSessionImpl(ptr: *anyopaque, id: []const u8) StorageError!void {
        const self: *Self = @ptrCast(@alignCast(ptr));

        if (self.sessions.fetchRemove(id)) |kv| {
            var session = kv.value;
            session.deinit(self.allocator);
            self.allocator.free(kv.key);
        }
    }

    fn getPendingSessionsImpl(ptr: *anyopaque, allocator: std.mem.Allocator) StorageError![]Session {
        const self: *Self = @ptrCast(@alignCast(ptr));

        var result: std.ArrayList(Session) = .empty;
        errdefer result.deinit(allocator);

        var it = self.sessions.iterator();
        while (it.next()) |entry| {
            if (entry.value_ptr.sync_status == SyncStatus.pending) {
                result.append(allocator, entry.value_ptr.*) catch return StorageError.OutOfMemory;
            }
        }

        return result.toOwnedSlice(allocator) catch return StorageError.OutOfMemory;
    }

    fn storageTypeImpl(ptr: *anyopaque) []const u8 {
        _ = ptr;
        return "memory";
    }
};

test "memory adapter add and get session" {
    const allocator = std.testing.allocator;

    var adapter = MemoryStorageAdapter.init(allocator);
    defer adapter.deinit();

    var storage = adapter.port();
    try storage.init();

    var session = try Session.init(allocator, "vape", "3", "puffs");
    defer session.deinit(allocator);

    _ = try storage.addSession(&session);

    const retrieved = try storage.getSession(session.getId());
    try std.testing.expect(retrieved != null);
    try std.testing.expectEqualStrings("vape", retrieved.?.method);
}
