const std = @import("std");
const core = @import("../../core/core.zig");
const Session = core.Session;
const SyncStatus = core.SyncStatus;
const StoragePort = core.StoragePort;
const StorageError = core.StorageError;

const c = @cImport({
    @cInclude("sqlite3.h");
});

/// SQLite-based storage adapter for persistent local storage
/// Implements the StoragePort interface for hexagonal architecture
pub const SqliteStorageAdapter = struct {
    db: ?*c.sqlite3,
    allocator: std.mem.Allocator,
    db_path: []const u8,

    const Self = @This();

    /// Create adapter with database file path
    pub fn init(allocator: std.mem.Allocator, db_path: []const u8) Self {
        return Self{
            .db = null,
            .allocator = allocator,
            .db_path = db_path,
        };
    }

    pub fn deinit(self: *Self) void {
        if (self.db) |db| {
            _ = c.sqlite3_close(db);
            self.db = null;
        }
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
        const self: *Self = @ptrCast(@alignCast(ptr));

        // Open database
        const result = c.sqlite3_open(self.db_path.ptr, &self.db);
        if (result != c.SQLITE_OK) {
            return StorageError.DatabaseError;
        }

        // Create sessions table if not exists
        const create_sql =
            \\CREATE TABLE IF NOT EXISTS sessions (
            \\    id TEXT PRIMARY KEY,
            \\    method TEXT NOT NULL,
            \\    amount TEXT NOT NULL,
            \\    unit TEXT NOT NULL,
            \\    notes TEXT DEFAULT '',
            \\    status TEXT DEFAULT 'active',
            \\    sync_status TEXT DEFAULT 'local',
            \\    timestamp INTEGER NOT NULL,
            \\    created_at INTEGER NOT NULL,
            \\    updated_at INTEGER NOT NULL
            \\);
        ;

        var err_msg: [*c]u8 = null;
        const exec_result = c.sqlite3_exec(self.db, create_sql, null, null, &err_msg);
        if (exec_result != c.SQLITE_OK) {
            if (err_msg) |msg| {
                c.sqlite3_free(msg);
            }
            return StorageError.DatabaseError;
        }
    }

    fn addSessionImpl(ptr: *anyopaque, session: *const Session) StorageError!Session {
        const self: *Self = @ptrCast(@alignCast(ptr));
        const db = self.db orelse return StorageError.DatabaseError;

        const insert_sql =
            \\INSERT INTO sessions (id, method, amount, unit, notes, status, sync_status, timestamp, created_at, updated_at)
            \\VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
        ;

        var stmt: ?*c.sqlite3_stmt = null;
        if (c.sqlite3_prepare_v2(db, insert_sql, -1, &stmt, null) != c.SQLITE_OK) {
            return StorageError.DatabaseError;
        }
        defer _ = c.sqlite3_finalize(stmt);

        // Bind parameters
        _ = c.sqlite3_bind_text(stmt, 1, session.getId().ptr, @intCast(session.getId().len), null);
        _ = c.sqlite3_bind_text(stmt, 2, session.method.ptr, @intCast(session.method.len), null);
        _ = c.sqlite3_bind_text(stmt, 3, session.amount.ptr, @intCast(session.amount.len), null);
        _ = c.sqlite3_bind_text(stmt, 4, session.unit.ptr, @intCast(session.unit.len), null);
        _ = c.sqlite3_bind_text(stmt, 5, session.notes.ptr, @intCast(session.notes.len), null);
        _ = c.sqlite3_bind_text(stmt, 6, statusToString(session.status).ptr, -1, null);
        _ = c.sqlite3_bind_text(stmt, 7, syncStatusToString(session.sync_status).ptr, -1, null);
        _ = c.sqlite3_bind_int64(stmt, 8, session.timestamp);
        _ = c.sqlite3_bind_int64(stmt, 9, session.created_at);
        _ = c.sqlite3_bind_int64(stmt, 10, session.updated_at);

        if (c.sqlite3_step(stmt) != c.SQLITE_DONE) {
            return StorageError.DatabaseError;
        }

        return session.*;
    }

    fn getSessionImpl(ptr: *anyopaque, id: []const u8) StorageError!?Session {
        const self: *Self = @ptrCast(@alignCast(ptr));
        const db = self.db orelse return StorageError.DatabaseError;

        const select_sql = "SELECT id, method, amount, unit, notes, status, sync_status, timestamp, created_at, updated_at FROM sessions WHERE id = ?;";

        var stmt: ?*c.sqlite3_stmt = null;
        if (c.sqlite3_prepare_v2(db, select_sql, -1, &stmt, null) != c.SQLITE_OK) {
            return StorageError.DatabaseError;
        }
        defer _ = c.sqlite3_finalize(stmt);

        _ = c.sqlite3_bind_text(stmt, 1, id.ptr, @intCast(id.len), null);

        if (c.sqlite3_step(stmt) == c.SQLITE_ROW) {
            return rowToSession(stmt, self.allocator) catch return StorageError.OutOfMemory;
        }

        return null;
    }

    fn getAllSessionsImpl(ptr: *anyopaque, allocator: std.mem.Allocator) StorageError![]Session {
        const self: *Self = @ptrCast(@alignCast(ptr));
        const db = self.db orelse return StorageError.DatabaseError;

        const select_sql = "SELECT id, method, amount, unit, notes, status, sync_status, timestamp, created_at, updated_at FROM sessions ORDER BY timestamp DESC;";

        var stmt: ?*c.sqlite3_stmt = null;
        if (c.sqlite3_prepare_v2(db, select_sql, -1, &stmt, null) != c.SQLITE_OK) {
            return StorageError.DatabaseError;
        }
        defer _ = c.sqlite3_finalize(stmt);

        var result: std.ArrayList(Session) = .empty;
        errdefer result.deinit(allocator);

        while (c.sqlite3_step(stmt) == c.SQLITE_ROW) {
            const session = rowToSession(stmt, allocator) catch return StorageError.OutOfMemory;
            result.append(allocator, session) catch return StorageError.OutOfMemory;
        }

        return result.toOwnedSlice(allocator) catch return StorageError.OutOfMemory;
    }

    fn updateSessionImpl(ptr: *anyopaque, session: *const Session) StorageError!Session {
        const self: *Self = @ptrCast(@alignCast(ptr));
        const db = self.db orelse return StorageError.DatabaseError;

        const update_sql =
            \\UPDATE sessions SET method = ?, amount = ?, unit = ?, notes = ?, status = ?, sync_status = ?, updated_at = ?
            \\WHERE id = ?;
        ;

        var stmt: ?*c.sqlite3_stmt = null;
        if (c.sqlite3_prepare_v2(db, update_sql, -1, &stmt, null) != c.SQLITE_OK) {
            return StorageError.DatabaseError;
        }
        defer _ = c.sqlite3_finalize(stmt);

        _ = c.sqlite3_bind_text(stmt, 1, session.method.ptr, @intCast(session.method.len), null);
        _ = c.sqlite3_bind_text(stmt, 2, session.amount.ptr, @intCast(session.amount.len), null);
        _ = c.sqlite3_bind_text(stmt, 3, session.unit.ptr, @intCast(session.unit.len), null);
        _ = c.sqlite3_bind_text(stmt, 4, session.notes.ptr, @intCast(session.notes.len), null);
        _ = c.sqlite3_bind_text(stmt, 5, statusToString(session.status).ptr, -1, null);
        _ = c.sqlite3_bind_text(stmt, 6, syncStatusToString(session.sync_status).ptr, -1, null);
        _ = c.sqlite3_bind_int64(stmt, 7, session.updated_at);
        _ = c.sqlite3_bind_text(stmt, 8, session.getId().ptr, @intCast(session.getId().len), null);

        if (c.sqlite3_step(stmt) != c.SQLITE_DONE) {
            return StorageError.DatabaseError;
        }

        if (c.sqlite3_changes(db) == 0) {
            return StorageError.NotFound;
        }

        return session.*;
    }

    fn deleteSessionImpl(ptr: *anyopaque, id: []const u8) StorageError!void {
        const self: *Self = @ptrCast(@alignCast(ptr));
        const db = self.db orelse return StorageError.DatabaseError;

        const delete_sql = "DELETE FROM sessions WHERE id = ?;";

        var stmt: ?*c.sqlite3_stmt = null;
        if (c.sqlite3_prepare_v2(db, delete_sql, -1, &stmt, null) != c.SQLITE_OK) {
            return StorageError.DatabaseError;
        }
        defer _ = c.sqlite3_finalize(stmt);

        _ = c.sqlite3_bind_text(stmt, 1, id.ptr, @intCast(id.len), null);

        if (c.sqlite3_step(stmt) != c.SQLITE_DONE) {
            return StorageError.DatabaseError;
        }
    }

    fn getPendingSessionsImpl(ptr: *anyopaque, allocator: std.mem.Allocator) StorageError![]Session {
        const self: *Self = @ptrCast(@alignCast(ptr));
        const db = self.db orelse return StorageError.DatabaseError;

        const select_sql = "SELECT id, method, amount, unit, notes, status, sync_status, timestamp, created_at, updated_at FROM sessions WHERE sync_status = 'pending' ORDER BY timestamp DESC;";

        var stmt: ?*c.sqlite3_stmt = null;
        if (c.sqlite3_prepare_v2(db, select_sql, -1, &stmt, null) != c.SQLITE_OK) {
            return StorageError.DatabaseError;
        }
        defer _ = c.sqlite3_finalize(stmt);

        var result: std.ArrayList(Session) = .empty;
        errdefer result.deinit(allocator);

        while (c.sqlite3_step(stmt) == c.SQLITE_ROW) {
            const session = rowToSession(stmt, allocator) catch return StorageError.OutOfMemory;
            result.append(allocator, session) catch return StorageError.OutOfMemory;
        }

        return result.toOwnedSlice(allocator) catch return StorageError.OutOfMemory;
    }

    fn storageTypeImpl(_: *anyopaque) []const u8 {
        return "sqlite";
    }

    // Helper functions
    fn rowToSession(stmt: ?*c.sqlite3_stmt, allocator: std.mem.Allocator) !Session {
        const id_ptr = c.sqlite3_column_text(stmt, 0);
        const method_ptr = c.sqlite3_column_text(stmt, 1);
        const amount_ptr = c.sqlite3_column_text(stmt, 2);
        const unit_ptr = c.sqlite3_column_text(stmt, 3);
        const notes_ptr = c.sqlite3_column_text(stmt, 4);
        const status_ptr = c.sqlite3_column_text(stmt, 5);
        const sync_status_ptr = c.sqlite3_column_text(stmt, 6);

        var id: [36]u8 = undefined;
        if (id_ptr) |ptr| {
            const len = @min(36, std.mem.len(ptr));
            @memcpy(id[0..len], ptr[0..len]);
            if (len < 36) @memset(id[len..], 0);
        }

        return Session{
            .id = id,
            .method = try dupeString(allocator, method_ptr),
            .amount = try dupeString(allocator, amount_ptr),
            .unit = try dupeString(allocator, unit_ptr),
            .notes = try dupeString(allocator, notes_ptr),
            .status = stringToStatus(status_ptr),
            .sync_status = stringToSyncStatus(sync_status_ptr),
            .timestamp = c.sqlite3_column_int64(stmt, 7),
            .created_at = c.sqlite3_column_int64(stmt, 8),
            .updated_at = c.sqlite3_column_int64(stmt, 9),
        };
    }

    fn dupeString(allocator: std.mem.Allocator, ptr: ?[*:0]const u8) ![]const u8 {
        if (ptr) |p| {
            const len = std.mem.len(p);
            if (len == 0) return "";
            return try allocator.dupe(u8, p[0..len]);
        }
        return "";
    }

    fn statusToString(status: core.SessionStatus) [:0]const u8 {
        return switch (status) {
            .active => "active",
            .completed => "completed",
        };
    }

    fn syncStatusToString(status: SyncStatus) [:0]const u8 {
        return switch (status) {
            .local => "local",
            .pending => "pending",
            .synced => "synced",
        };
    }

    fn stringToStatus(ptr: ?[*:0]const u8) core.SessionStatus {
        if (ptr) |p| {
            const str = std.mem.span(p);
            if (std.mem.eql(u8, str, "completed")) return .completed;
        }
        return .active;
    }

    fn stringToSyncStatus(ptr: ?[*:0]const u8) SyncStatus {
        if (ptr) |p| {
            const str = std.mem.span(p);
            if (std.mem.eql(u8, str, "pending")) return .pending;
            if (std.mem.eql(u8, str, "synced")) return .synced;
        }
        return .local;
    }
};

test "sqlite adapter init" {
    const allocator = std.testing.allocator;

    var adapter = SqliteStorageAdapter.init(allocator, ":memory:");
    defer adapter.deinit();

    var storage = adapter.port();
    try storage.init();

    try std.testing.expectEqualStrings("sqlite", storage.storageType());
}
