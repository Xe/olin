const errs = @import("./error.zig");
const Allocator = @import("std").mem.Allocator;

extern fn startup_arg_len() usize;
extern fn startup_arg_at(id: usize, arg_data: [*]u8, arg_len: i32) i32;

pub fn args(allocator: *Allocator) ![][]u8 {
    const argc = startup_arg_len();
    var result = try allocator.alloc([]u8, argc);

    for (result) |v, i| {
        var arg_string: []u8 = try allocator.alloc(u8, 512);
        const size_of_arg = try errs.parse(startup_arg_at(i, arg_string.ptr, 512));
        result[i] = arg_string[0..@intCast(usize, size_of_arg)];
    }

    return result;
}

pub fn free_args(allocator: *Allocator, argv: [][]u8) void {
    for (argv) |v| {
        allocator.free(v);
    }

    allocator.free(argv);
}
