const olin = @import("./olin/olin.zig");
const Resource = olin.resource.Resource;

export fn cwa_main() i32 {
  main() catch unreachable;
  return 0;
}

fn main() !void {
  const out = try Resource.stdout();
  const msg = "Hello, world!\n";
  const ign = try out.write(&msg, msg.len);
}
