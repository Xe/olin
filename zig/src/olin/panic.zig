const builtin = @import("builtin");
const log = @import("log");

pub fn panic(msg: []const u8, error_return_trace: ?*builtin.StackTrace) noreturn {
    log.err(msg);

    unreachable;
}
