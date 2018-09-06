#[link(wasm_import_module = "cwa")]
extern "C" {
    pub fn log_write(level: i32, text_ptr: *const u8, text_len: usize);
    pub fn env_get(
        key_ptr: *const u8, key_len: usize,
        value_buf_ptr: *mut u8, value_buf_len: usize
    ) -> i32;
    pub fn runtime_spec_major() -> i32;
    pub fn runtime_spec_minor() -> i32;
    pub fn runtime_name(out_ptr: *mut u8, out_len: usize) -> i32;
    pub fn runtime_msleep(ms: i32);
    pub fn startup_arg_len() -> i32;
    pub fn startup_arg_at(id: i32, out_ptr: *mut u8, out_len: usize) -> i32;

    pub fn resource_open(url_ptr: *const u8, url_len: usize) -> i32;
    pub fn resource_read(id: i32, data_ptr: *mut u8, data_len: usize) -> i32;
    pub fn resource_write(id: i32, data_ptr: *const u8, data_len: usize) -> i32;
    pub fn resource_close(id: i32);

    pub fn io_get_stdin() -> i32;
    pub fn io_get_stdout() -> i32;
    pub fn io_get_stderr() -> i32;
}
