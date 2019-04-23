const errs = @import("./error.zig");
const Allocator = @import("std").mem.Allocator;

extern fn runtime_spec_minor() i32;
extern fn runtime_spec_major() i32;
extern fn runtime_name(name_data: [*]const u8, name_len: usize) i32;

pub const Metadata = struct {
    spec_minor: i32,
    spec_major: i32,
    name: []u8,
};

// The user must free the Metadata.
pub fn metadata(alloc: *Allocator) !*Metadata {
    const runtime_name_limit = 32;
    var result: *Metadata = undefined;
    var name: []u8 = undefined;

    result = try alloc.create(Metadata);
    name = try alloc.alloc(u8, runtime_name_limit);
    result.spec_minor = runtime_spec_minor();
    result.spec_major = runtime_spec_major();

    const name_len = try errs.parse(runtime_name(name.ptr, runtime_name_limit));
    result.name = name[0..@intCast(usize, name_len)];

    return result;
}

extern fn runtime_msleep(ms_len: i32) void;

pub fn sleep(ms_len: i32) void {
    runtime_msleep(ms_len);
}


