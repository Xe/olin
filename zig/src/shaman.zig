const olin = @import("./olin/olin.zig");
pub const os = olin;
pub const panic = os.panic;
const Resource = olin.resource.Resource;

pub fn main() anyerror!void {
    const fout = try Resource.stdout();
    const data = @embedFile("./shaman.aa");
    const n = try fout.write(data);
    fout.close();
}


