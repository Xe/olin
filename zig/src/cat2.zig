pub const os = @import("./olin/olin.zig");
const std = @import("std");
const alloc = std.heap.page_allocator;

pub fn main() anyerror!void {
    std.debug.warn("All your base are belong to us.\n", .{});
}
