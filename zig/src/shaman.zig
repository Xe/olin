const olin = @import("./olin/olin.zig");
const Resource = olin.resource.Resource;

export fn _start() noreturn {
    olin.runtime.exit(main() catch olin.runtime.exit(1));
}

fn main() !i32 {
    const fout = try Resource.stdout();
    const data = @embedFile("./shaman.aa");
    const n = try fout.write(data);
    fout.close();
    return 0;
}


