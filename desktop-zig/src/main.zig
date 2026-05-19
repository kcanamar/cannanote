const std = @import("std");
const dvui = @import("dvui");

const core = @import("core/core.zig");
const adapters = @import("adapters/adapters.zig");

const Session = core.Session;
const SessionService = core.SessionService;
const SqliteStorageAdapter = adapters.SqliteStorageAdapter;

// Database file path (in user's home directory)
const DB_PATH = "cannanote.db";

// ============================================================================
// DVUI App Configuration
// ============================================================================

pub const dvui_app: dvui.App = .{
    .config = .{
        .options = .{
            .size = .{ .w = 500.0, .h = 450.0 },
            .min_size = .{ .w = 400.0, .h = 350.0 },
            .title = "CannaNote",
        },
    },
    .frameFn = appFrame,
    .initFn = appInit,
    .deinitFn = appDeinit,
};

pub const main = dvui.App.main;
pub const panic = dvui.App.panic;
pub const std_options: std.Options = .{
    .logFn = dvui.App.logFn,
};

// ============================================================================
// Application State (Global for dvui frame access)
// ============================================================================

var app_state: ?*AppState = null;

const AppState = struct {
    allocator: std.mem.Allocator,
    storage: SqliteStorageAdapter,
    service: SessionService,

    // Form state
    method_index: usize = 0,
    amount_buffer: [32]u8 = [_]u8{0} ** 32,
    status_message: []const u8 = "",

    // Session list
    sessions: []Session = &[_]Session{},

    const methods = [_][]const u8{ "vape", "smoke", "edible", "tincture" };

    pub fn init(allocator: std.mem.Allocator) !*AppState {
        const self = try allocator.create(AppState);
        self.* = AppState{
            .allocator = allocator,
            .storage = SqliteStorageAdapter.init(allocator, DB_PATH),
            .service = undefined,
            .sessions = &[_]Session{},
        };

        const port = self.storage.port();
        try port.init();
        self.service = SessionService.init(allocator, port);

        // Load existing sessions from database
        self.refreshSessions();

        return self;
    }

    pub fn deinit(self: *AppState) void {
        if (self.sessions.len > 0) {
            self.allocator.free(self.sessions);
        }
        self.storage.deinit();
        self.allocator.destroy(self);
    }

    pub fn refreshSessions(self: *AppState) void {
        if (self.sessions.len > 0) {
            self.allocator.free(self.sessions);
        }
        self.sessions = self.service.getAllSessions() catch &[_]Session{};
    }

    pub fn logSession(self: *AppState) void {
        const method = methods[self.method_index];

        // Find actual text length (scan for null terminator)
        var len: usize = 0;
        for (self.amount_buffer) |byte| {
            if (byte == 0) break;
            len += 1;
        }

        const amount = self.amount_buffer[0..len];

        if (amount.len == 0) {
            self.status_message = "Please enter an amount";
            return;
        }

        _ = self.service.logSession(method, amount, "puffs") catch {
            self.status_message = "Error logging session";
            return;
        };

        self.status_message = "Session logged!";
        @memset(&self.amount_buffer, 0);
        self.refreshSessions();
    }
};

// ============================================================================
// DVUI Lifecycle Functions
// ============================================================================

fn appInit(_: *dvui.Window) anyerror!void {
    const allocator = std.heap.page_allocator;
    app_state = try AppState.init(allocator);
}

fn appDeinit() void {
    if (app_state) |state| {
        state.deinit();
        app_state = null;
    }
}

fn appFrame() anyerror!dvui.App.Result {
    const state = app_state orelse return .close;

    // Title
    dvui.label(@src(), "CannaNote", .{}, .{});
    _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 20 } });

    // Log Session section
    {
        dvui.label(@src(), "Log Session", .{}, .{});
        _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 10 } });

        // Method label and dropdown
        dvui.labelNoFmt(@src(), "Method:", .{}, .{});
        _ = dvui.dropdown(@src(), &AppState.methods, .{ .choice = &state.method_index }, .{}, .{});

        _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 10 } });

        // Amount label and input
        dvui.labelNoFmt(@src(), "Amount:", .{}, .{});
        var te = dvui.textEntry(@src(), .{ .text = .{ .buffer = &state.amount_buffer } }, .{});
        te.deinit();

        _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 10 } });

        // Log button
        if (dvui.button(@src(), "Log Session", .{}, .{})) {
            state.logSession();
        }

        // Status message
        if (state.status_message.len > 0) {
            _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 5 } });
            dvui.labelNoFmt(@src(), state.status_message, .{}, .{});
        }
    }

    _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 20 } });

    // Sessions section
    {
        dvui.label(@src(), "Sessions", .{}, .{});

        if (dvui.button(@src(), "Refresh", .{}, .{})) {
            state.refreshSessions();
        }

        _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 10 } });

        // Session list
        if (state.sessions.len == 0) {
            dvui.labelNoFmt(@src(), "No sessions yet. Log your first session above!", .{}, .{});
        } else {
            for (state.sessions, 0..) |session, i| {
                var buf: [128]u8 = undefined;
                const label_text = std.fmt.bufPrint(&buf, "{s} {s} ({s})", .{
                    session.amount,
                    session.unit,
                    session.method,
                }) catch "Error";
                dvui.labelNoFmt(@src(), label_text, .{}, .{ .id_extra = i });
            }
        }
    }

    return .ok;
}

// ============================================================================
// Tests
// ============================================================================

test {
    @import("std").testing.refAllDecls(@This());
    _ = @import("core/core.zig");
    _ = @import("adapters/adapters.zig");
}
