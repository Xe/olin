const olin = @import("./olin/olin.zig");
const env = olin.env;
const log = olin.log;
const random = olin.random;
const resource = olin.resource;
const time = olin.time;
const runtime = olin.runtime;
const startup = olin.startup;

pub const os = olin;
pub const panic = os.panic;

const std = @import("std");
const assert = std.debug.assert;
var alloc = std.heap.page_allocator;

pub fn main() anyerror!void {
    std.os.exit(do());
}

fn do() u8 {
    log.info("hi");
    log.warning("hi");
    log.err("hi");

    const ai32 = random.int32();
    const bi32 = random.int32();

    assert(ai32 != bi32);

    const ai64 = random.int64();
    const bi64 = random.int64();

    assert(ai64 != bi64);

    const now = time.unix();

    assert(now != 0);

    test_runtime_sleep();
    test_time();

    test_resource_log() catch return 1;
    test_resource_random() catch return 2;
    test_env_get() catch return 3;
    test_runtime_metadata() catch return 4;
    test_startup_args() catch return 5;

    return 0;
}

fn test_resource_log() !void {
    const msg = "hi there";
    const fout = try resource.open("log://?prefix=test");
    const n = try fout.write(msg);
    fout.close();
}

fn test_resource_random() !void {
    const fin = try resource.open("random://");
    var buf: []u8 = try alloc.alloc(u8, 32);
    defer alloc.free(buf);

    _ = try fin.read(buf);
}

fn test_env_get() !void {
    const key = "MAGIC_CONCH";
    log.info("getting env");
    const val = try env.get(alloc, key);

    log.info(val);
    alloc.free(val);
}

fn test_runtime_sleep() void {
    runtime.sleep(1);
}

fn test_runtime_metadata() !void {
    const metadata = try runtime.metadata(alloc);

    log.info(metadata.name);

    alloc.destroy(metadata);
}

fn test_time() void {
    const now = time.unix();
    @import("std").debug.assert(now != 0);
}

fn test_startup_args() !void {
    var args = try startup.args(alloc);

    for (args) |v| {
        log.info(v);
    }

    startup.free_args(alloc, args);
}
