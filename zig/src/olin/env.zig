const errs = @import("./error.zig");
const Allocator = @import("std").mem.Allocator;

extern fn env_get(
    key_data: [*]const u8,
    key_len: usize,
    value_data: [*]u8,
    value_len: usize,
) i32;

pub fn get(allocator: *Allocator, key: []const u8) ![]u8 {
    const needed_size_i = env_get(key.ptr, key.len, @intToPtr([*]u8, 1), 0);
    const needed_size_n = try errs.parse(needed_size_i);
    const needed_size = @intCast(usize, needed_size_n);

    var buf: []u8 = try allocator.alloc(u8, @intCast(usize, needed_size));

    const result: i32 = env_get(key.ptr, key.len, buf.ptr, needed_size);
    const value_len = try errs.parse(result);

    return buf[0..@intCast(usize, value_len)];
}
