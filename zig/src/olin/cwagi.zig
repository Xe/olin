const fmt = @import("std").fmt;
const Allocator = @import("std").mem.Allocator;

const env_get = @import("./env.zig").get;
const Resource = @import("./resource.zig").Resource;

pub const Context = struct {
    method: []u8,
    request_uri: []u8,
    body: Resource,

    pub fn init(allocator: *Allocator) !Context {
        return Context {
            .method = try env_get(allocator, "REQUEST_METHOD"),
            .request_uri = try env_get(allocator, "REQUEST_URI"),
            .body = try Resource.stdin(),
        };
    }

    pub fn destroy(self: Context, allocator: *Allocator) void {
        allocator.free(self.method);
        allocator.free(self.request_uri);
    }
};

pub const Response = struct {
    status: u32,
    body: []u8,

    pub fn writeTo(self: Response, allocator: *Allocator, fout: Resource) !void {
        var header_tmp = try allocator.alloc(u8, 2048);

        const preamble = try fmt.bufPrint(header_tmp, "HTTP/1.1 {}\n", self.status);
        var n = try fout.write(preamble.ptr, preamble.len);

        const headers = @embedFile("./headers.txt");
        n = try fout.write(&headers, headers.len);

        const twoLines = "\n\n";
        n = try fout.write(&twoLines, twoLines.len);

        n = try fout.write(self.body.ptr, self.body.len);

        allocator.free(header_tmp);
    }
};
