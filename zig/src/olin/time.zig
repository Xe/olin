extern fn time_now() i64;

pub fn unix() i64 {
    return time_now();
}

test "now isn't zero" {
    const now = time_now();
    @import("std").debug.assert(now != 0);
}
