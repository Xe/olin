extern fn log_write(level: i32, data: [*]const u8, len: usize) void;

const log_lvl_error = 1;
const log_lvl_warning = 3;
const log_lvl_info = 6;

pub fn err(data: []const u8) void {
    log_write(log_lvl_error, data.ptr, data.len);
}

test "err" {
    err("unknown");
}

pub fn warning(data: []const u8) void {
    log_write(log_lvl_warning, data.ptr, data.len);
}

test "warning" {
    warning("unknown");
}

pub fn info(data: []const u8) void {
    log_write(log_lvl_info, data.ptr, data.len);
}

test "info" {
    info("info");
}
