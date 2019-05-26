export fn cwa_main() i32 {
    force_crash() catch unreachable;

    return 0;
}

pub fn main() void {}

extern fn random_i32() i32;
const ErrorStuff = error{
    Stuff,
};

fn force_crash() !void {
    if (@mod(random_i32(), 2) == 1) {
        return error.Stuff;
    }
}
