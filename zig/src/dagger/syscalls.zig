// External functions from the dagger API
pub extern fn open(data: [*]const u8, len: usize) i32;
pub extern fn read(fd: i32, data: [*]const u8, len: usize) i32;
pub extern fn write(fd: i32, data: [*]const u8, len: usize) i32;
pub extern fn close(fd: i32) i32;
pub extern fn flush(fd: i32) i32;
