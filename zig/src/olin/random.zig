extern fn random_i32() i32;
extern fn random_i64() i64;

pub fn int32() i32 { return random_i32(); }
pub fn int64() i64 { return random_i64(); }

test "i32 randomness" {
    const a = int32();
    const b = int32();

    assert(a != b);
}

test "i64 randomness" {
    const a = int64();
    const b = int64();

    assert(a != b);
}
