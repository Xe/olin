const dagger = @import("./dagger/dagger.zig");
const Stream = dagger.Stream;

export fn dagger_main() i32 {
    return main() catch 1;
}

fn main() !i32 {
    const l = try Stream.log();
    _ = try l.write_slice("hello");

    const stdout = try Stream.stdout();
    _ = try stdout.write_slice("hello\n");
    const stderr = try Stream.stderr();
    _ = try stderr.write_slice("hello\n");

    const fout = try Stream.open("file:///testdata/test");

    _ = try l.close();
    _ = try stdout.flush();
    _ = try stdout.close();
    _ = try stderr.flush();
    _ = try stderr.close();
    return 0;
}
