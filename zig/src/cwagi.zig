const olin = @import("./olin/olin.zig");
const cwagi = olin.cwagi;

const std = @import("std");
const fmt = std.fmt;
var alloc = std.heap.wasm_allocator;

pub fn main() void {}

export fn cwa_main() i32 {
    const ctx = cwagi.Context.init(alloc) catch return 1;
    const resp: cwagi.Response = helloWorld(ctx) catch return 1;

    const fout = olin.resource.stdout() catch return 1;
    resp.writeTo(alloc, fout) catch return 1;

    return 0;
}

const message = "Hello, I am served from Zig compiled to wasm32-freestanding-none.\n\nI know the following about my environment:\n";

fn helloWorld(ctx: cwagi.Context) !cwagi.Response {
    var line: []u8 = try alloc.alloc(u8, 128);

    var buf = try std.Buffer.init(alloc, message);
    const runtime_meta = try olin.runtime.metadata(alloc);

    try buf.append(try fmt.bufPrint(line, "- I am running in {} which implements version {}.{} of the Common WebAssembly ABI.\n", runtime_meta.name, runtime_meta.spec_major, runtime_meta.spec_minor));
    try buf.append(try fmt.bufPrint(line, "- I think the time is {}\n", olin.time.unix()));
    try buf.append(try fmt.bufPrint(line, "- RUN_ID:    {}\n", try olin.env.get(alloc, "RUN_ID")));
    try buf.append(try fmt.bufPrint(line, "- WORKER_ID: {}\n", try olin.env.get(alloc, "WORKER_ID")));
    try buf.append(try fmt.bufPrint(line, "- Method:    {}\n", ctx.method));
    try buf.append(try fmt.bufPrint(line, "- URI:       {}\n", ctx.request_uri));

    alloc.free(line);

    try buf.append("\n\n");
    try buf.append(@embedFile("./cwagi_message.txt"));

    return cwagi.Response {
         .status = 200,
         .body = buf.toSlice(),
    };
}


