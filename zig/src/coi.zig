const olin = @import("./olin/olin.zig");
const log = olin.log;
const random = olin.random;
const resource = olin.resource;
const time = olin.time;

const std = @import("std");
const assert = std.debug.assert;

export fn cwa_main() i32 {
    var fails: i32 = 0;

    if (test_log() != 0) {
        log.err("log failed");
        fails = fails + 1;
    }

    if (test_random() != 0) {
        log.err("random failed");
        fails = fails + 1;
    }

    if (test_time() != 0) {
        log.err("time failed");
        fails = fails + 1;
    }

    if (test_resource_log() != 0) {
        log.err("resource log failed");
        fails = fails + 1;
    }

    return fails;
}

fn test_log() i32 {
    log.info("hi");
    log.warning("hi");
    log.err("hi");

    return 0;
}

fn test_random() i32 {
    const ai32 = random.int32();
    const bi32 = random.int32();

    assert(ai32 != bi32);

    const ai64 = random.int64();
    const bi64 = random.int64();

    assert(ai64 != bi64);

    return 0;
}

fn test_time() i32 {
    const now = time.unix();

    assert(now != 0);

    return 0;
}

fn test_resource_log() i32 {
    if(resource.Resource.open("log://?prefix=test")) |fout| {
        const msg = "hi there";
        if(fout.write(&msg, msg.len)) |name| {
            return 0;
        } else |err| {
            log.err(@errorName(err));
            return 1;
        }
    } else |err| {
        log.err(@errorName(err));
        return 1;
    }
}
