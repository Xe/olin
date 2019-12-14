pub const os = @import("./olin/olin.zig");
pub const panic = os.panic;
const std = @import("std");

pub fn main() anyerror!void {
    const fin = std.io.getStdIn();
    const fout = std.io.getStdOut();
    const buflen = 928;
    var buf: [buflen]u8 = undefined;
    var bufSlice = buf[0..];

    while (true) {
        var n: usize = undefined;
        if (fin.read(bufSlice)) |retVal| {
            n = retVal;
        } else |errVal| {
            @panic(@errorName(errVal));
        }
        try fout.write(bufSlice[0..n]);

        if (n < buflen) {
            break;
        }
    }
}
