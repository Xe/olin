pub const os = @import("./olin/olin.zig");
pub const panic = os.panic;
pub const std = @import("std");

pub fn main() anyerror!void {
    std.os.exit(1);
}
