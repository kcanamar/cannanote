const std = @import("std");
const Session = @import("../domain/domain.zig").Session;
const StoragePort = @import("../ports/ports.zig").StoragePort;
const StorageError = @import("../ports/ports.zig").StorageError;

/// Application service - orchestrates domain + ports
/// This is where business logic lives. It depends on ports (interfaces),
/// not concrete implementations, enabling easy testing and swapping.
pub const SessionService = struct {
    storage: StoragePort,
    allocator: std.mem.Allocator,

    const Self = @This();

    pub fn init(allocator: std.mem.Allocator, storage: StoragePort) Self {
        return Self{
            .storage = storage,
            .allocator = allocator,
        };
    }

    /// Log a new session
    pub fn logSession(
        self: *Self,
        method: []const u8,
        amount: []const u8,
        unit: []const u8,
        strain: []const u8,
        mood_before: u8,
        mind_before: u8,
        body_before: u8,
        notes: []const u8,
    ) !Session {
        var session = try Session.init(self.allocator, method, amount, unit, strain, mood_before, mind_before, body_before, notes);
        errdefer session.deinit(self.allocator);

        // Domain validation
        try session.validate();

        // Persist via port
        return try self.storage.addSession(&session);
    }

    /// Get all sessions, newest first
    pub fn getAllSessions(self: *Self) ![]Session {
        return try self.storage.getAllSessions(self.allocator);
    }

    /// Get single session by ID
    pub fn getSession(self: *Self, id: []const u8) !?Session {
        return try self.storage.getSession(id);
    }

    /// Complete a session
    pub fn completeSession(self: *Self, id: []const u8) !Session {
        var session = try self.storage.getSession(id) orelse return StorageError.NotFound;

        session.complete();
        return try self.storage.updateSession(&session);
    }

    /// Delete a session
    pub fn deleteSession(self: *Self, id: []const u8) !void {
        return try self.storage.deleteSession(id);
    }

    /// Get sessions pending sync
    pub fn getPendingSessions(self: *Self) ![]Session {
        return try self.storage.getPendingSessions(self.allocator);
    }
};
