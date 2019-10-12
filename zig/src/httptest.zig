const olin = @import("./olin/olin.zig");
const log = olin.log;
const resource = olin.resource;

const assert = @import("std").debug.assert;
const mem = @import("std").mem;
const fmt = @import("std").fmt;
const heap = @import("std").heap;
const Headers = @import("std").http.Headers;

const userAgent = "Olin+Zig@master";

export fn cwa_main() i32 {
    log.info("making request to http://127.0.0.1:3000/foo.txt");

    doRequest(heap.wasm_allocator) catch return 1;

    return 0;
}

fn doRequest(alloc: *mem.Allocator) !void {
    const fout = try resource.open("http://127.0.0.1:3000");
    var buf: []u8 = undefined;
    buf = try alloc.alloc(u8, 256);
    defer alloc.free(buf);
    var h = Headers.init(alloc);
    defer h.deinit();
    try h.append("User-Agent", userAgent, null);
    try h.append("Host", "127.0.0.1:3000", null);

    var res = try fmt.bufPrint(buf[0..], "GET /foo.txt HTTP/1.1\n{}\n\n", h);
    const n = try fout.write(res);
    log.info(res);

    try fout.flush();

    var resp: []u8 = undefined;
    resp = try alloc.alloc(u8, 2048);
    defer alloc.free(resp);

    const nresp = try fout.read(resp);
    log.info(resp[0..nresp]);
}
