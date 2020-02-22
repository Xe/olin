const olin = @import("./olin/olin.zig");
pub const os = olin;
pub const panic = os.panic;
const log = olin.log;
const runtime = olin.runtime;
const std = @import("std");
var alloc = std.heap.page_allocator;

pub fn main() anyerror!void {
    const metadata = try runtime.metadata(alloc);

    log.info(metadata.name);

    alloc.destroy(metadata);
}
