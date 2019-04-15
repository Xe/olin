const errs = @import("./error.zig");

extern fn env_get(
    key_data: [*]const u8,
    key_len: usize,
    value_data: [*]const u8,
    value_len: usize,
) i32;

pub fn get(key: []const u8) ![]u8 {
    const len = 2048;
    var buf: [len]u8 = undefined;
    const result: i32 = env_get(key.ptr, key.len, &buf, len);
    const value_len = errs.parse(result) catch |err| return err;
    return buf[0..@intCast(usize, value_len)];
}
