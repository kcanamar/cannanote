const std = @import("std");
const dvui = @import("dvui");

const core = @import("core/core.zig");
const adapters = @import("adapters/adapters.zig");

const Session = core.Session;
const SessionService = core.SessionService;
const SqliteStorageAdapter = adapters.SqliteStorageAdapter;

// Database file path
const DB_PATH = "cannanote.db";

// ============================================================================
// DVUI App Configuration
// ============================================================================

pub const dvui_app: dvui.App = .{
    .config = .{
        .options = .{
            .size = .{ .w = 900.0, .h = 600.0 },
            .min_size = .{ .w = 700.0, .h = 400.0 },
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
// Application State
// ============================================================================

var app_state: ?*AppState = null;

const FocusPanel = enum { form, list };
const FormField = enum { method, unit, amount, strain, mood, mind, body, notes };
const Mode = enum { normal, insert };

const AppState = struct {
    allocator: std.mem.Allocator,
    storage: SqliteStorageAdapter,
    service: SessionService,

    // Vim-like mode
    mode: Mode = .normal,
    focus_text_entry: bool = false, // Flag to focus text entry on next frame

    // Focus state for vim navigation
    focus_panel: FocusPanel = .form,
    form_field: FormField = .method,
    list_index: usize = 0,

    // Form state
    method_index: usize = 0,
    unit_index: usize = 0,
    mood_index: usize = 0,
    mind_index: usize = 0,
    body_index: usize = 0,
    amount_buffer: [32]u8 = [_]u8{0} ** 32,
    strain_buffer: [64]u8 = [_]u8{0} ** 64,
    notes_buffer: [256]u8 = [_]u8{0} ** 256,
    status_message: []const u8 = "",
    status_is_error: bool = false,

    // Session list
    sessions: []Session = &[_]Session{},

    // Options (matching web app)
    const methods = [_][]const u8{ "vape", "smoke", "edible", "tincture", "topical", "dab" };
    const units = [_][]const u8{ "puffs", "seconds", "hits", "grams", "joints", "mg", "ml", "drops", "pieces", "applications" };
    const scale_5 = [_][]const u8{ "-", "1", "2", "3", "4", "5" };

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
        // Clamp list index
        if (self.list_index >= self.sessions.len and self.sessions.len > 0) {
            self.list_index = self.sessions.len - 1;
        }
    }

    pub fn logSession(self: *AppState) void {
        const method = methods[self.method_index];
        const unit = units[self.unit_index];
        const amount = getBufferText(&self.amount_buffer);

        if (amount.len == 0) {
            self.setStatus("Amount required", true);
            return;
        }

        const strain = getBufferText(&self.strain_buffer);
        const notes = getBufferText(&self.notes_buffer);
        const mood: u8 = if (self.mood_index == 0) 0 else @intCast(self.mood_index);
        const mind: u8 = if (self.mind_index == 0) 0 else @intCast(self.mind_index);
        const body: u8 = if (self.body_index == 0) 0 else @intCast(self.body_index);

        _ = self.service.logSession(method, amount, unit, strain, mood, mind, body, notes) catch {
            self.setStatus("Error logging session", true);
            return;
        };

        self.setStatus("Session logged!", false);
        self.clearForm();
        self.refreshSessions();
    }

    pub fn deleteSelectedSession(self: *AppState) void {
        if (self.sessions.len == 0) return;
        if (self.list_index >= self.sessions.len) return;

        const session = self.sessions[self.list_index];
        self.service.deleteSession(session.getId()) catch {
            self.setStatus("Error deleting", true);
            return;
        };
        self.setStatus("Deleted", false);
        self.refreshSessions();
    }

    fn setStatus(self: *AppState, msg: []const u8, is_error: bool) void {
        self.status_message = msg;
        self.status_is_error = is_error;
    }

    fn clearForm(self: *AppState) void {
        @memset(&self.amount_buffer, 0);
        @memset(&self.strain_buffer, 0);
        @memset(&self.notes_buffer, 0);
        self.mood_index = 0;
        self.mind_index = 0;
        self.body_index = 0;
    }

    fn getBufferText(buffer: []const u8) []const u8 {
        var len: usize = 0;
        for (buffer) |byte| {
            if (byte == 0) break;
            len += 1;
        }
        return buffer[0..len];
    }

    // Vim navigation
    pub fn navDown(self: *AppState) void {
        switch (self.focus_panel) {
            .form => {
                const fields = std.enums.values(FormField);
                const idx = @intFromEnum(self.form_field);
                if (idx < fields.len - 1) {
                    self.form_field = fields[idx + 1];
                }
            },
            .list => {
                if (self.list_index < self.sessions.len -| 1) {
                    self.list_index += 1;
                }
            },
        }
    }

    pub fn navUp(self: *AppState) void {
        switch (self.focus_panel) {
            .form => {
                const fields = std.enums.values(FormField);
                const idx = @intFromEnum(self.form_field);
                if (idx > 0) {
                    self.form_field = fields[idx - 1];
                }
            },
            .list => {
                if (self.list_index > 0) {
                    self.list_index -= 1;
                }
            },
        }
    }

    pub fn navLeft(self: *AppState) void {
        self.focus_panel = .form;
    }

    pub fn navRight(self: *AppState) void {
        self.focus_panel = .list;
    }

    pub fn cycleFieldValue(self: *AppState, delta: i32) void {
        if (self.focus_panel != .form) return;

        switch (self.form_field) {
            .method => self.method_index = cycleIndex(self.method_index, methods.len, delta),
            .unit => self.unit_index = cycleIndex(self.unit_index, units.len, delta),
            .mood => self.mood_index = cycleIndex(self.mood_index, scale_5.len, delta),
            .mind => self.mind_index = cycleIndex(self.mind_index, scale_5.len, delta),
            .body => self.body_index = cycleIndex(self.body_index, scale_5.len, delta),
            else => {},
        }
    }

    fn cycleIndex(current: usize, max: usize, delta: i32) usize {
        if (delta > 0) {
            return if (current < max - 1) current + 1 else 0;
        } else {
            return if (current > 0) current - 1 else max - 1;
        }
    }

    // Check if current field is a text field
    pub fn isOnTextField(self: *AppState) bool {
        return self.focus_panel == .form and
            (self.form_field == .amount or self.form_field == .strain or self.form_field == .notes);
    }

    // Check if current field is a dropdown
    pub fn isOnDropdown(self: *AppState) bool {
        return self.focus_panel == .form and
            (self.form_field == .method or self.form_field == .unit or
            self.form_field == .mood or self.form_field == .mind or self.form_field == .body);
    }

    // Enter insert mode (only on text fields)
    pub fn enterInsertMode(self: *AppState) void {
        if (self.isOnTextField()) {
            self.mode = .insert;
            self.focus_text_entry = true; // Request focus on next render
        }
    }

    // Exit to normal mode
    pub fn exitInsertMode(self: *AppState) void {
        self.mode = .normal;
    }
};

// ============================================================================
// DVUI Lifecycle
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

    // Handle global keyboard input
    handleKeyboardInput(state);

    // Main layout: horizontal split
    var paned = dvui.paned(@src(), .{
        .direction = .horizontal,
        .collapsed_size = 400,
    }, .{ .expand = .both });

    // Left panel: Form
    if (paned.showFirst()) {
        renderFormPanel(state);
    }

    // Right panel: Session list
    if (paned.showSecond()) {
        renderListPanel(state);
    }

    paned.deinit();

    return .ok;
}

fn handleKeyboardInput(state: *AppState) void {
    for (dvui.events()) |*event| {
        if (event.evt != .key) continue;
        if (event.handled) continue;

        const key = event.evt.key;
        if (key.action != .down and key.action != .repeat) continue;

        // INSERT MODE - only handle Escape
        if (state.mode == .insert) {
            if (key.code == .escape) {
                state.exitInsertMode();
                event.handled = true;
            }
            // Let dvui handle all other keys for text entry
            continue;
        }

        // NORMAL MODE - full vim navigation
        switch (key.code) {
            .j => {
                state.navDown();
                event.handled = true;
            },
            .k => {
                state.navUp();
                event.handled = true;
            },
            .h => {
                state.navLeft();
                event.handled = true;
            },
            .l => {
                state.navRight();
                event.handled = true;
            },
            .tab => {
                if (key.mod.has(.lshift) or key.mod.has(.rshift)) {
                    state.navUp();
                } else {
                    state.navDown();
                }
                event.handled = true;
            },
            .i => {
                // Enter insert mode on text fields
                if (state.isOnTextField()) {
                    state.enterInsertMode();
                    event.handled = true;
                }
            },
            .space => {
                // Cycle dropdown values
                if (state.isOnDropdown()) {
                    state.cycleFieldValue(1);
                    event.handled = true;
                }
            },
            .enter => {
                // Enter ONLY submits the session
                if (state.focus_panel == .form) {
                    state.logSession();
                    event.handled = true;
                }
            },
            .d, .x => {
                if (state.focus_panel == .list) {
                    state.deleteSelectedSession();
                    event.handled = true;
                }
            },
            .g => {
                if (state.focus_panel == .list) {
                    state.list_index = 0;
                } else if (state.focus_panel == .form) {
                    state.form_field = .method;
                }
                event.handled = true;
            },
            .n => {
                state.focus_panel = .form;
                state.form_field = .method;
                event.handled = true;
            },
            .r => {
                state.refreshSessions();
                event.handled = true;
            },
            .escape => {
                // Clear status or other cancel action
                state.status_message = "";
                event.handled = true;
            },
            else => {},
        }
    }
}

fn renderFormPanel(state: *AppState) void {
    var scroll = dvui.scrollArea(@src(), .{}, .{ .expand = .both });
    defer scroll.deinit();

    const form_focused = state.focus_panel == .form;

    // Mode indicator
    const mode_str: []const u8 = if (state.mode == .insert) "-- INSERT --" else "-- NORMAL --";
    dvui.labelNoFmt(@src(), mode_str, .{}, .{ .id_extra = 950 });
    _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 5 }, .id_extra = 951 });

    // Title with focus indicator
    const title: []const u8 = if (form_focused) "[ New Session ]" else "  New Session  ";
    dvui.labelNoFmt(@src(), title, .{}, .{});
    _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 5 }, .id_extra = 952 });

    // Help text based on mode
    const help: []const u8 = if (state.mode == .insert)
        "Esc:normal"
    else
        "j/k:nav i:edit Space:cycle Enter:submit h/l:panel";
    dvui.labelNoFmt(@src(), help, .{}, .{ .id_extra = 900 });
    _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 10 }, .id_extra = 901 });

    // Method
    renderDropdownField(state, "Method", &AppState.methods, &state.method_index, .method, form_focused, 0);

    // Unit
    renderDropdownField(state, "Unit", &AppState.units, &state.unit_index, .unit, form_focused, 1);

    // Amount
    renderTextField(state, "Amount", &state.amount_buffer, .amount, form_focused, 0);

    // Strain
    renderTextField(state, "Strain", &state.strain_buffer, .strain, form_focused, 1);

    _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 8 }, .id_extra = 100 });
    dvui.labelNoFmt(@src(), "Pre-Session State (1-5)", .{}, .{ .id_extra = 101 });
    _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 4 }, .id_extra = 102 });

    // Mood Before
    renderDropdownField(state, "Mood", &AppState.scale_5, &state.mood_index, .mood, form_focused, 2);

    // Mind Before
    renderDropdownField(state, "Mind", &AppState.scale_5, &state.mind_index, .mind, form_focused, 3);

    // Body Before
    renderDropdownField(state, "Body", &AppState.scale_5, &state.body_index, .body, form_focused, 4);

    // Notes
    renderTextField(state, "Notes", &state.notes_buffer, .notes, form_focused, 2);

    _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 10 }, .id_extra = 200 });

    // Log button
    if (dvui.button(@src(), "Log Session [Enter]", .{}, .{})) {
        state.logSession();
    }

    // Status
    if (state.status_message.len > 0) {
        _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 5 }, .id_extra = 201 });
        dvui.labelNoFmt(@src(), state.status_message, .{}, .{ .id_extra = 202 });
    }
}

fn renderDropdownField(
    state: *AppState,
    label_text: []const u8,
    options: []const []const u8,
    index: *usize,
    field: FormField,
    panel_focused: bool,
    id_extra: usize,
) void {
    const is_focused = panel_focused and state.form_field == field;

    var hbox = dvui.box(@src(), .{ .dir = .horizontal }, .{ .id_extra = id_extra });
    defer hbox.deinit();

    // Focus indicator - more visible
    const indicator: []const u8 = if (is_focused) "[>]" else "   ";
    dvui.labelNoFmt(@src(), indicator, .{}, .{ .id_extra = id_extra });

    // Label with padding
    var label_buf: [32]u8 = undefined;
    const padded_label = std.fmt.bufPrint(&label_buf, "{s: <8}", .{label_text}) catch label_text;
    dvui.labelNoFmt(@src(), padded_label, .{}, .{ .id_extra = id_extra + 100 });

    // Dropdown
    _ = dvui.dropdown(@src(), options, .{ .choice = index }, .{}, .{ .id_extra = id_extra + 200 });
}

fn renderTextField(
    state: *AppState,
    label_text: []const u8,
    buffer: []u8,
    field: FormField,
    panel_focused: bool,
    id_extra: usize,
) void {
    const is_focused = panel_focused and state.form_field == field;
    const is_editing = is_focused and state.mode == .insert;

    var hbox = dvui.box(@src(), .{ .dir = .horizontal }, .{ .id_extra = id_extra + 300 });
    defer hbox.deinit();

    // Focus indicator - shows mode
    const indicator: []const u8 = if (is_editing) "[I]" else if (is_focused) "[>]" else "   ";
    dvui.labelNoFmt(@src(), indicator, .{}, .{ .id_extra = id_extra + 400 });

    // Label with padding
    var label_buf: [32]u8 = undefined;
    const padded_label = std.fmt.bufPrint(&label_buf, "{s: <8}", .{label_text}) catch label_text;
    dvui.labelNoFmt(@src(), padded_label, .{}, .{ .id_extra = id_extra + 500 });

    if (is_editing) {
        // INSERT MODE: Show editable text entry
        var te = dvui.textEntry(@src(), .{ .text = .{ .buffer = buffer } }, .{ .id_extra = id_extra + 600 });

        // Focus the text entry if just entered insert mode
        if (state.focus_text_entry) {
            dvui.focusWidget(te.data().id, null, null);
            state.focus_text_entry = false;
        }

        te.deinit();
    } else {
        // NORMAL MODE: Show as clickable label
        const text_value = AppState.getBufferText(buffer);
        var display_buf: [64]u8 = undefined;
        const display = if (text_value.len > 0)
            std.fmt.bufPrint(&display_buf, "[{s}]", .{text_value}) catch "[...]"
        else
            "[empty - press i to edit]";

        // Clickable button that enters insert mode
        if (dvui.button(@src(), display, .{}, .{ .id_extra = id_extra + 600 })) {
            state.form_field = field;
            state.focus_panel = .form;
            state.enterInsertMode();
        }
    }
}

fn renderListPanel(state: *AppState) void {
    var scroll = dvui.scrollArea(@src(), .{}, .{ .expand = .both, .id_extra = 1 });
    defer scroll.deinit();

    const list_focused = state.focus_panel == .list;

    // Title with focus indicator
    const title: []const u8 = if (list_focused) "[ Sessions ]" else "  Sessions  ";
    dvui.labelNoFmt(@src(), title, .{}, .{ .id_extra = 1 });
    _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 5 }, .id_extra = 1 });

    // Help
    dvui.labelNoFmt(@src(), "j/k:nav d/x:delete g:top n:new", .{}, .{ .id_extra = 800 });
    _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 8 }, .id_extra = 801 });

    // Refresh button
    if (dvui.button(@src(), "Refresh [r]", .{}, .{ .id_extra = 1 })) {
        state.refreshSessions();
    }

    _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 8 }, .id_extra = 2 });

    // Session count
    var count_buf: [32]u8 = undefined;
    const count_text = std.fmt.bufPrint(&count_buf, "{d} sessions", .{state.sessions.len}) catch "?";
    dvui.labelNoFmt(@src(), count_text, .{}, .{ .id_extra = 802 });

    _ = dvui.spacer(@src(), .{ .min_size_content = .{ .h = 8 }, .id_extra = 3 });

    if (state.sessions.len == 0) {
        dvui.labelNoFmt(@src(), "No sessions yet", .{}, .{ .id_extra = 700 });
    } else {
        for (state.sessions, 0..) |session, i| {
            const is_selected = list_focused and i == state.list_index;

            // Build display string
            var buf: [256]u8 = undefined;
            var time_buf: [20]u8 = undefined;
            const time_str = formatTimestamp(session.timestamp, &time_buf);

            const display = blk: {
                if (session.strain.len > 0) {
                    break :blk std.fmt.bufPrint(&buf, "{s}{s} | {s} {s} {s}", .{
                        if (is_selected) @as([]const u8, ">") else @as([]const u8, " "),
                        time_str,
                        session.amount,
                        session.unit,
                        session.strain,
                    }) catch "Error";
                } else {
                    break :blk std.fmt.bufPrint(&buf, "{s}{s} | {s} {s} ({s})", .{
                        if (is_selected) @as([]const u8, ">") else @as([]const u8, " "),
                        time_str,
                        session.amount,
                        session.unit,
                        session.method,
                    }) catch "Error";
                }
            };

            dvui.labelNoFmt(@src(), display, .{}, .{ .id_extra = i });

            // Delete button
            var del_buf: [8]u8 = undefined;
            const del_label = std.fmt.bufPrint(&del_buf, "[x]", .{}) catch "x";
            if (dvui.button(@src(), del_label, .{}, .{ .id_extra = i + 1000 })) {
                state.list_index = i;
                state.deleteSelectedSession();
            }
        }
    }
}

const c_time = @cImport({
    @cInclude("time.h");
});

fn formatTimestamp(timestamp: i64, buf: []u8) []const u8 {
    // Use C localtime for proper timezone handling
    var time_val: c_time.time_t = @intCast(timestamp);
    const local = c_time.localtime(&time_val);

    if (local) |tm| {
        return std.fmt.bufPrint(buf, "{d:0>2}/{d:0>2} {d:0>2}:{d:0>2}", .{
            @as(u8, @intCast(tm.*.tm_mon + 1)),
            @as(u8, @intCast(tm.*.tm_mday)),
            @as(u8, @intCast(tm.*.tm_hour)),
            @as(u8, @intCast(tm.*.tm_min)),
        }) catch "??/??";
    }

    return "??/??";
}

// ============================================================================
// Tests
// ============================================================================

test {
    _ = @import("core/core.zig");
    _ = @import("adapters/adapters.zig");
}
