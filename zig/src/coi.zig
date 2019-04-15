const olin = @import("./olin/olin.zig");
const env = olin.env;
const log = olin.log;
const random = olin.random;
const resource = olin.resource;
const time = olin.time;

const std = @import("std");
const assert = std.debug.assert;

export fn cwa_main() i32 {
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

    test_resource_log() catch return 1;
    test_env_get() catch return 2;

    return 0;
}

fn test_resource_log() !void {
    const msg = "hi there";
    const open = resource.Resource.open;
    const fout = try open("log://?prefix=test");
    const n = try fout.write(&msg, msg.len);
}

fn test_env_get() !void {
    const key = "MAGIC_CONCH";
    log.info("getting env");
    const val = try olin.env.get(key);
    const cmp = c"yes";

    log.info(val);
}
