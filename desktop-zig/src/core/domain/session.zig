const std = @import("std");
const c = @cImport({
    @cInclude("time.h");
});

/// Get current unix timestamp
fn getCurrentTimestamp() i64 {
    return @intCast(c.time(null));
}

/// Session status in the logging lifecycle
pub const SessionStatus = enum {
    active,
    completed,
};

/// Sync status for offline-first architecture
pub const SyncStatus = enum {
    local,
    pending,
    synced,
};

/// Core domain entity - no I/O, no frameworks
/// Pure data and business logic only
pub const Session = struct {
    id: [36]u8, // UUID as string
    method: []const u8,
    amount: []const u8,
    unit: []const u8,
    notes: []const u8,
    status: SessionStatus,
    sync_status: SyncStatus,
    timestamp: i64, // Unix timestamp
    created_at: i64,
    updated_at: i64,

    const Self = @This();

    /// Create new session with defaults
    pub fn init(
        allocator: std.mem.Allocator,
        method: []const u8,
        amount: []const u8,
        unit: []const u8,
    ) !Self {
        const now = getCurrentTimestamp();

        // Generate UUID v4 using PRNG seeded with timestamp
        var id: [36]u8 = undefined;
        var uuid_bytes: [16]u8 = undefined;

        // Use timestamp and address for entropy (simple approach for demo)
        var prng = std.Random.DefaultPrng.init(@intCast(@as(u64, @bitCast(now)) ^ @intFromPtr(&id)));
        prng.fill(&uuid_bytes);

        // Set version (4) and variant bits
        uuid_bytes[6] = (uuid_bytes[6] & 0x0f) | 0x40;
        uuid_bytes[8] = (uuid_bytes[8] & 0x3f) | 0x80;

        // Format as string
        _ = std.fmt.bufPrint(&id, "{x:0>2}{x:0>2}{x:0>2}{x:0>2}-{x:0>2}{x:0>2}-{x:0>2}{x:0>2}-{x:0>2}{x:0>2}-{x:0>2}{x:0>2}{x:0>2}{x:0>2}{x:0>2}{x:0>2}", .{
            uuid_bytes[0],  uuid_bytes[1],  uuid_bytes[2],  uuid_bytes[3],
            uuid_bytes[4],  uuid_bytes[5],  uuid_bytes[6],  uuid_bytes[7],
            uuid_bytes[8],  uuid_bytes[9],  uuid_bytes[10], uuid_bytes[11],
            uuid_bytes[12], uuid_bytes[13], uuid_bytes[14], uuid_bytes[15],
        }) catch unreachable;

        return Self{
            .id = id,
            .method = try allocator.dupe(u8, method),
            .amount = try allocator.dupe(u8, amount),
            .unit = try allocator.dupe(u8, unit),
            .notes = "",
            .status = .active,
            .sync_status = .local,
            .timestamp = now,
            .created_at = now,
            .updated_at = now,
        };
    }

    /// Domain validation - no external dependencies
    pub fn validate(self: *const Self) !void {
        if (self.method.len == 0) {
            return error.MethodRequired;
        }
        if (self.amount.len == 0) {
            return error.AmountRequired;
        }
    }

    /// Mark session complete
    pub fn complete(self: *Self) void {
        self.status = .completed;
        self.updated_at = getCurrentTimestamp();
    }

    /// Free allocated memory
    pub fn deinit(self: *Self, allocator: std.mem.Allocator) void {
        if (self.method.len > 0) allocator.free(self.method);
        if (self.amount.len > 0) allocator.free(self.amount);
        if (self.unit.len > 0) allocator.free(self.unit);
        if (self.notes.len > 0) allocator.free(self.notes);
    }

    /// Get ID as slice
    pub fn getId(self: *const Self) []const u8 {
        return &self.id;
    }
};

/// Domain errors
pub const DomainError = error{
    MethodRequired,
    AmountRequired,
    InvalidSession,
};

test "session creation" {
    const allocator = std.testing.allocator;

    var session = try Session.init(allocator, "vape", "3", "puffs");
    defer session.deinit(allocator);

    try std.testing.expectEqualStrings("vape", session.method);
    try std.testing.expectEqualStrings("3", session.amount);
    try std.testing.expect(session.status == .active);
}

test "session validation" {
    const allocator = std.testing.allocator;

    var session = try Session.init(allocator, "vape", "3", "puffs");
    defer session.deinit(allocator);

    try session.validate();
}
