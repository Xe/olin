pub const os = @import("./olin/olin.zig");
pub const panic = os.panic;
const std = @import("std");

pub fn main() anyerror!void {
    std.debug.warn("All your base are belong to us.\n", .{});
}
