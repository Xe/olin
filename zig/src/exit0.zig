const exit = @import("./olin/olin.zig").runtime.exit;

export fn _start() noreturn {
    exit(0);
}
