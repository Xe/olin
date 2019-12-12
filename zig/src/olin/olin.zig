pub const env = @import("./env.zig");
pub const err = @import("./error.zig");
pub const log = @import("./log.zig");
pub const random = @import("./random.zig");
pub const resource = @import("./resource.zig");
pub const time = @import("./time.zig");
pub const runtime = @import("./runtime.zig");
pub const startup = @import("./startup.zig");

// cwagi
pub const cwagi = @import("./cwagi.zig");

// not directly used, but imported like this to force the compiler to actually consider it.
pub const panic = @import("./panic.zig").panic;

fn hack(inp: i32) usize {
    if (err.parse(inp)) |resp| {
        return @intCast(usize, resp);
    } else |errVal| {
        errno = @intCast(u12, inp*-1);
        return 0;
    }
}

pub export fn _start() noreturn {
    runtime.exit(@intCast(i32, @import("std").special.start.callMain()));
}

pub const io = struct {
    pub fn getStdInHandle() bits.fd_t {
        if (bits.STDIN_FILENO == -1) {
            bits.STDIN_FILENO = resource.io_get_stdin();
        }

        return bits.STDIN_FILENO;
    }

    pub fn getStdErrHandle() bits.fd_t {
        if (bits.STDERR_FILENO == -1) {
            bits.STDERR_FILENO = resource.io_get_stderr();
        }

        return bits.STDERR_FILENO;
    }

    pub fn getStdOutHandle() bits.fd_t {
        if (bits.STDOUT_FILENO == -1) {
            bits.STDOUT_FILENO = resource.io_get_stdout();
        }

        return bits.STDOUT_FILENO;
    }
};

pub const system = struct {
    const std = @import("std");
    const mem = std.mem;

    pub fn raise(sig: bits.signal_t) noreturn {
        runtime.exit(sig);
    }

    pub fn write(fd: i32, buf: [*]const u8, count: usize) usize {
        return hack(resource.resource_write(fd, buf, count));
    }

    pub fn read(fd: i32, buf: [*]u8, count: usize) usize {
        return hack(resource.resource_read(fd, buf, count));
    }

    pub fn close(fd: i32) usize {
        resource.resource_close(fd);
    }

    pub fn openat(fd: bits.fd_t, path: [*:0]const u8, flags: u32, mode: usize) usize {
        return open(path, flags, mode);
    }

    pub fn open(path: [*:0]const u8, flags: u32, mode: usize) usize {
        const inner_path = mem.toSlice(path);
        return hack(resource.resource_open(inner_path.ptr, inner_path.len));
    }

    pub fn lseek(fd: i32, offset: i64, whence: u2) i64 {
        return std.os.EINVAL;
    }

    pub fn getErrno(arg: i64) u12 {
        return errno;
    }

    pub fn exit(status: u8) noreturn {
        runtime.exit(@intCast(i32, status));
    }

    pub fn abort() noreturn {
        exit(1);
    }
};

var errno: u12 = 0;

pub const bits = @import("bits.zig");
