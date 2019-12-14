const olin = @import("./olin/olin.zig");
pub const os = olin;
pub const panic = os.panic;
const cwagi = olin.cwagi;

const std = @import("std");
const fmt = std.fmt;
const Headers = std.http.Headers;
var alloc = std.heap.page_allocator;

pub fn main() anyerror!void {
    const fout = try olin.resource.stdout();
    const ctx = try cwagi.Context.init(alloc);
    defer ctx.destroy(alloc);

    var h = Headers.init(alloc);
    defer h.deinit();

    const resp: cwagi.Response = try helloWorld(ctx, &h);
    try resp.writeTo(alloc, &h, fout);
}

const message = "Hello, I am served from Zig compiled to wasm32-freestanding-none.\n\nI know the following about my environment:\n";

fn helloWorld(ctx: cwagi.Context, headers: *Headers) !cwagi.Response {
    var line: [2048]u8 = undefined;
    try headers.append("Content-Type", "text/plain", null);
    try headers.append("Olin-Lang", "Zig", null);

    var buf = try std.Buffer.init(alloc, message);
    const runtime_meta = try olin.runtime.metadata(alloc);
    const runID = try olin.env.get(alloc, "RUN_ID");
    defer alloc.free(runID);
    const workerID = try olin.env.get(alloc, "WORKER_ID");
    defer alloc.free(workerID);

    try buf.append(try fmt.bufPrint(line[0..], "- I am running in {} which implements version {}.{} of the Common WebAssembly ABI.\n- I think the time is {}\n- RUN_ID:    {}\n- WORKER_ID: {}\n- Method:    {}\n- URI:       {}\n\n",
                                    .{
                                        runtime_meta.name,
                                        runtime_meta.spec_major,
                                        runtime_meta.spec_minor,
                                        olin.time.unix(),
                                        runID,
                                        workerID,
                                        ctx.method,
                                        ctx.request_uri,
                                     }
                                    ));

    try buf.append(@embedFile("./cwagi_message.txt"));

    return cwagi.Response {
         .status = olin.http.StatusCode.OK,
         .body = buf.toSlice(),
    };
}
