const dagger = @import("./dagger/dagger.zig");
const Stream = dagger.Stream;

export fn dagger_main() i32 {
  main() catch unreachable;
  return 0;
}

fn main() !void {
  const out = try Stream.open("stdout://");
  const msg = "Hello, world!\n";
  const ign = try out.write(&msg, msg.len);
}
