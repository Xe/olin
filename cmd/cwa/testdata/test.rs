extern crate core;
use std::str;

//#[link(name = "cwa", wasm_import_module = "cwa")]
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
}

#[repr(i32)]
pub enum Level {
    Error = 1,
    Warning = 3,
    Info = 6
}

/// Writes a line of log with the specified level to host logger.
pub fn log(level: Level, text: &str) {
    let text = text.as_bytes();

    unsafe {
        log_write(
            level as i32,
            text.as_ptr(),
            text.len()
        );
    }
}


#[no_mangle]
pub extern "C" fn cwa_main() -> i32 {
    log(Level::Info, "expecting spec major=0 and min=0");
    let minor: i32;
    let major: i32;

    unsafe {
        minor = runtime_spec_minor();
        major = runtime_spec_major();
    }

    if major != 0 && minor != 0 {
        log(Level::Error, &format!("minor: {}, major: {}", minor, major));
        return 1;
    }
    log(Level::Info, "passed");

    log(Level::Info, "getting runtime name, should be olin");
    let mut rt_name = [0u8; 16];
    let res: i32;
    unsafe {
        res = runtime_name(rt_name.as_mut_ptr(), 16);
    }

    if res < 0 {
        return res;
    }

    let runtime_name_str: &str;
    unsafe {
        runtime_name_str = core::str::from_utf8_unchecked(&rt_name[..res as usize]);
    }

    log(Level::Info, runtime_name_str);

    if runtime_name_str != "olin" {
        log(Level::Error, "Got runtime name, not olin");
        return 1;
    }
    log(Level::Info, "passed");

    log(Level::Info, "sleeping");
    unsafe {
        runtime_msleep(1);
    }
    log(Level::Info, "passed");

    log(Level::Info, "checking argc/argv");
    let argc: i32;
    unsafe {
        argc = startup_arg_len();
    }
    log(Level::Info, &format!("argc: {}", argc));

    for x in 0..argc {
        let mut arg_val = [0u8; 64];
        let res: i32;
        unsafe {
            res = startup_arg_at(x as i32, arg_val.as_mut_ptr(), 64);
        }

        if res < 0 {
            return res;
        }

        let arg_str: &str;
        unsafe {
            arg_str = core::str::from_utf8_unchecked(&arg_val[..res as usize]);
        }

        log(Level::Info, &format!("arg {}: {}", x, arg_str));
    }
    log(Level::Info, "passed");

    log(Level::Info, "env[\"MAGIC_CONCH\"] = \"yes\"");
    let envvar_name = "MAGIC_CONCH";
    let res: i32;
    let mut envvar_val = [0u8; 64];
    let val_str: &str;
    unsafe {
        res = env_get(envvar_name.as_bytes().as_ptr(), envvar_name.len(), envvar_val.as_mut_ptr(), 64);
    }

    if res < 0 {
        return res;
    }

    unsafe {
        val_str = core::str::from_utf8_unchecked(&envvar_val[..res as usize]);
    }

    if val_str != "yes" {
        log(Level::Error, &format!("wanted yes, got: {}", val_str));
        return 1;
    }
    log(Level::Info, "passed");

    return 0;
}

fn main() {}
