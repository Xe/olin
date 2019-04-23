const errs = @import("./error.zig");
const OlinError = errs.OlinError;

extern fn resource_open(data: [*]const u8, len: usize) i32;
extern fn resource_read(fd: i32, data: [*]const u8, len: usize) i32;
extern fn resource_write(fd: i32, data: [*]const u8, len: usize) i32;
extern fn resource_close(fd: i32) void;
extern fn resource_flush(fd: i32) i32;

extern fn io_get_stdin() i32;
extern fn io_get_stdout() i32;
extern fn io_get_stderr() i32;

fn fd_check(fd: i32) errs.OlinError!Resource {
    if(errs.parse(fd)) |fd_for_handle| {
        return Resource {
            .fd = fd_for_handle,
        };
    } else |err| {
        return err;
    }
}

pub const open = Resource.open;
pub const stdin = Resource.stdin;
pub const stdout = Resource.stdout;
pub const stderr = Resource.stderr;

pub const Resource = struct {
    fd: i32,

    pub fn open(url: []const u8) !Resource {
        const fd = resource_open(url.ptr, url.len);
        return fd_check(fd);
    }

    pub fn stdin() !Resource{
        const fd = io_get_stdin();
        return fd_check(fd);
    }

    pub fn stdout() !Resource{
        const fd = io_get_stdout();
        return fd_check(fd);
    }

    pub fn stderr() !Resource{
        const fd = io_get_stderr();
        return fd_check(fd);
    }

    pub fn write(self: Resource, data: [*]const u8, len: usize) OlinError!i32 {
        const n = resource_write(self.fd, data, len);

        if (errs.parse(n)) |nresp| {
            return nresp;
        } else |err| {
            return err;
        }
    }

    pub fn read(self: Resource, data: [*]u8, len: usize) OlinError!i32 {
        const n = resource_read(self.fd, data, len);

        if (errs.parse(n)) |nresp| {
            return nresp;
        } else |err| {
            return err;
        }
    }

    pub fn close(self: Resource) void {
        resource_close(self.fd);
        return;
    }

    pub fn flush(self: Resource) OlinError!void {
        const n = resource_flush(self.fd);

        if (errs.parse(n)) {} else |err| {
            return err;
        }
    }
};

