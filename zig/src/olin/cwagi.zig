const fmt = @import("std").fmt;
const Allocator = @import("std").mem.Allocator;
const Headers = @import("std").http.Headers;

const http = @import("./http.zig");
const StatusCode = http.StatusCode;
const env_get = @import("./env.zig").get;
const Resource = @import("./resource.zig").Resource;

pub const Context = struct {
    method: []const u8,
    request_uri: []const u8,
    server_name: []const u8,
    body: Resource,
    content_length: u32,

    pub fn init(allocator: *Allocator) !Context {
        const content_length_str = env_get(allocator, "CONTENT_LENGTH") catch "0"[0..];
        defer allocator.free(content_length_str);
        return Context {
            .method = try env_get(allocator, "REQUEST_METHOD"),
            .request_uri = try env_get(allocator, "REQUEST_URI"),
            .server_name = try env_get(allocator, "SERVER_NAME"),
            .body = try Resource.stdin(),
            .content_length = try fmt.parseInt(u32, content_length_str, 10),
        };
    }

    pub fn destroy(self: Context, allocator: *Allocator) void {
        allocator.free(self.method);
        allocator.free(self.request_uri);
        allocator.free(self.server_name);
        self.body.close();
    }
};

pub const Response = struct {
    status: StatusCode,
    body: []u8,

    pub fn writeTo(self: Response, allocator: *Allocator, headers: *Headers, fout: Resource) !void {
        var header_tmp: [8192]u8 = undefined;

        var contentLength: [16]u8 = undefined;
        const clLen = fmt.formatIntBuf(contentLength[0..], self.body.len, 10, false, fmt.FormatOptions{});
        try headers.append("Content-Length", contentLength[0..clLen], null);

        const preamble = try fmt.bufPrint(header_tmp[0..], "Status: {} {}\n{}\n", .{@enumToInt(self.status), http.reasonPhrase(self.status), headers});

        _ = try fout.write(preamble);
        _ = try fout.write(self.body);
    }
};
