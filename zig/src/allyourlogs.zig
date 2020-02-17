const olin = @import("./olin/olin.zig");
pub const os = olin;
pub const panic = os.panic;
const std = @import("std");
const log = olin.log;

pub fn main() anyerror!void {
    log.info("hi");
    log.warning("hi");
    log.err("hi");
    std.os.exit(0);
}
