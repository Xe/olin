extern fn log_write(level: i32, data: [*]const u8, len: usize) void;

const log_error = 1;
const log_warning = 3;
const log_info = 6;

export fn cwa_main() i32 {
    const msg = "coi la munje\n";
    log_write(log_info, &msg, msg.len);

    return 0;
}
