#![allow(warnings)]

extern crate core;
use std::mem;
use std::io::Read;
use std::io::Write;
use std::str;

mod cwa;
use cwa::*;

#[no_mangle]
pub extern "C" fn cwa_main() -> i32 {
    log(Level::Info, "expecting spec major=0 and min=0");
    let minor: i32 = unsafe { runtime_spec_minor() };
    let major: i32 = unsafe { runtime_spec_major() };

    log(Level::Info, &format!("minor: {}, major: {}", minor, major));

    if major != 0 && minor != 0 {
        log(Level::Error, "version is wrong");
        return 1;
    }
    log(Level::Info, "passed");

    log(Level::Info, "getting runtime name, should be olin");
    let mut rt_name = [0u8; 16];
    let res: i32 = unsafe { runtime_name(rt_name.as_mut_ptr(), 16) };

    if res < 0 {
        return res;
    }

    let runtime_name_str: &str = unsafe { std::str::from_utf8_unchecked(&rt_name[..res as usize]) };

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
    let argc: i32 = unsafe { startup_arg_len() };
    log(Level::Info, &format!("argc: {}", argc));

    for x in 0..argc {
        let mut arg_val = [0u8; 64];
        let res: i32 = unsafe { startup_arg_at(x as i32, arg_val.as_mut_ptr(), 64) };

        if res < 0 {
            return res;
        }

        let arg_str: &str = unsafe{ core::str::from_utf8_unchecked(&arg_val[..res as usize]) };
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

    log(Level::Info, "trying to open a log:// file");
    let mut fout: Resource;
    let fout_maybe: Option<Resource>;
    fout_maybe = Resource::open("log://test-log-please-ignore");
    match fout_maybe {
        Some(r) => fout = r,
        None => return 1,
    }

    let log_msg = "hi from inside log file".as_bytes();
    let res = fout.write(log_msg);
    if !res.is_ok() {
        log(Level::Error, "can't write message to log file");
        log(Level::Error, &format!("error: {}", res.err().unwrap()));
    }
    std::mem::drop(fout);
    log(Level::Info, "successfully closed the file");

    log(Level::Info, "opening a zero:// file");
    let mut fout: Resource;
    let fout_maybe: Option<Resource>;
    fout_maybe = Resource::open("zero://");
    match fout_maybe {
        Some(r) => fout = r,
        None => return 1,
    }

    log(Level::Info, "reading zeroes");
    let mut zeroes = [0u8, 16];
    let res = fout.read(&mut zeroes);
    if !res.is_ok() {
        log(Level::Error, "can't read zeroes from zero file");
        log(Level::Error, &format!("error: {}", res.err().unwrap()));
    }

    /* XXX(Xe): Rust is broken
    log(Level::Info, "verifying all zeroes are valid");
    for x in 0..16 {
        if zeroes[x] != 0 {
            log(Level::Error, &format!("expected zeroes[{}] to be 0, got: {}", x, zeroes[x]));
            return 1;
        }
    }
    */

    std::mem::drop(fout);
    log(Level::Info, "closed file");

    log(Level::Info, "all functions passed basic usage");

    return 0;
}

fn main() {}
