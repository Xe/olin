const olin = @import("./olin/olin.zig");
pub const os = olin;
pub const panic = os.panic;
const log = olin.log;
const startup = olin.startup;
const std = @import("std");
const fmt = std.fmt;
const alloc = std.heap.page_allocator;

pub fn main() anyerror!void {
    var args = try startup.args(alloc);
    defer startup.free_args(alloc, args);

    for (args) |arg| {
        log.info(arg);
    }
}
