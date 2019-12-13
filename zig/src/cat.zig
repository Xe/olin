const olin = @import("./olin/olin.zig");
const std = @import("std");

export fn _start() i32 {
    if (inner_main()) {} else |err| {
        olin.log.err(@errorName(err));
        return 1;
    }

    return 0;
}

fn inner_main() !void {
    const fin = try olin.resource.stdin();
    const fout = try olin.resource.stdout();
    const buflen = 928;
    var buf: [buflen]u8 = undefined;
    var bufSlice = buf[0..];

    while (true) {
        const n = try fin.read(bufSlice);
        const nn = try fout.write(bufSlice[0..n]);

        if (n != nn) {
            unreachable;
        }

        if (n < buflen) {
            break;
        }
    }
}
