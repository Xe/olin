const syscalls = @import("./syscalls.zig");

pub const DaggerError = error {
    Unknown,
};

fn parse_error(inp: i32) DaggerError!i32 {
    switch (inp) {
        -1 => {
            return error.Unknown;
        },
        else => {
            return inp;
        },
    }
}

pub const Stream = struct{
    sd: i32,

    pub fn open(url: []const u8)!Stream {
        const sd = try parse_error(syscalls.open(url.ptr, url.len));
        return Stream{
            .sd = sd,
        };
    }

    pub fn stdin() !Stream{
        return open("stdin://");
    }

    pub fn stdout() !Stream{
        return open("stdout://");
    }

    pub fn stderr() !Stream{
        return open("stderr://");
    }

    pub fn log() !Stream{
        return open("log://");
    }

    pub fn write_slice(self: Stream, data: []const u8) !i32 {
        return self.write(data.ptr, data.len);
    }

    pub fn write(self: Stream, data: [*]const u8, len: usize) !i32{
        const n = try parse_error(syscalls.write(self.sd, data, len));
        return n;
    }

    pub fn read(self: Stream, data: [*]u8, len: usize) !i32{
        const n = try parse_error(syscalls.read(self.sd, data, len));
        return n;
    }

    pub fn close(self: Stream) !void {
        const unused = try parse_error(syscalls.close(self.sd));
    }

    pub fn flush(self: Stream) !void {
        const unused = try parse_error(syscalls.flush(self.sd));
    }
};
