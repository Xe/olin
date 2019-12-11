const olin = @import("./olin/olin.zig");
const std = @import("std");
var alloc = @import("std").heap.wasm_allocator;

export fn cwa_main() i32 {
    if (inner_main()) {} else |err| {
        olin.log.err(@errorName(err));
        return 1;
    }

    return 0;
}

fn inner_main() !void {
    const fin = try olin.resource.stdin();
    defer fin.close();
    const fout = try olin.resource.stdout();
    defer fout.close();
    const buflen = 2048;
    var buf: [buflen]u8 = undefined;
    var bufSlice = buf[0..];

    while (true) {
        const n = try fin.read(bufSlice);
        const nn = try fout.write(bufSlice);

        if (n != nn) {
            unreachable;
        }

        if (n < buflen) {
            break;
        }
    }
}
