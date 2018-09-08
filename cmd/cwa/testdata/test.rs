#![allow(warnings)]

use std::io::{Read, Write};
use std::mem;
use std::str;

mod cwa;
use cwa::*;

#[no_mangle]
pub extern "C" fn cwa_main() -> i32 {
    match friendly_main() {
        Ok(()) => 0,
        Err(e) => e,
    }
}

#[inline(always)]
pub extern "C" fn friendly_main() -> Result<(), i32> {
    log::info("expecting spec major=0 and min=0");
    let minor: i32 = runtime::spec_major();
    let major: i32 = runtime::spec_major();

    log::info(&format!("minor: {}, major: {}", minor, major));

    if major != 0 && minor != 0 {
        log::error("version is wrong");
        return Err(1);
    }
    log::info("passed");

    log::info("getting runtime name, should be olin");
    let mut rt_name = [0u8; 32];
    let runtime_name = runtime::name_buf(rt_name.as_mut()).ok_or_else(|| {
        log::error("Runtime name larger than 32 byte limit");
        1
    })?;

    log::info(runtime_name);

    if runtime_name != "olin" {
        log::error("Got runtime name, not olin");
        return Err(1);
    }
    log::info("passed");

    log::info("sleeping");
    runtime::msleep(1);
    log::info("passed");

    log::info("checking argc/argv");
    let argc: i32 = startup::arg_len();
    log::info(&format!("argc: {}", argc));

    for x in 0..argc {
        let mut arg_val = [0u8; 64];
        let arg = startup::arg_at_buf(x, &mut arg_val).ok_or_else(|| {
            log::error(&format!("arg {} missing", x));
            1
        })?;
        log::info(&format!("arg {}: {}", x, arg));
    }
    log::info("passed");

    log::info("env[\"MAGIC_CONCH\"] = \"yes\"");
    let envvar_name = "MAGIC_CONCH";
    let res: i32;
    let mut envvar_val = [0u8; 64];
    let envvar_val = env::get_buf(envvar_name.as_bytes(), &mut envvar_val)
        .map(|s| str::from_utf8(s).expect("envvar wasn't UTF-8"))
        .map_err(|e| {
            log::error(&format!("couldn't get: {:?}", e));
            1
        })?;

    if envvar_val != "yes" {
        log::error(&format!("wanted yes, got: {}", envvar_val));
        return Err(1);
    }
    log::info("passed");

    log::info("trying to open a log:// file");
    {
        let mut fout: Resource =
            Resource::open("log://?prefix=test-log-please-ignore").map_err(|e| {
                log::error(&format!("couldn't open: {:?}", e));
                1
            })?;

        let res = fout.write(b"hi from inside log file");
        if let Err(err) = res {
            log::error("can't write message to log file");
            log::error(&format!("error: {:?}", err));
        }
    }
    log::info("successfully closed the file");

    log::info("opening a zero:// file");
    {
        let mut fout: Resource = Resource::open("zero://").map_err(|e| {
            log::error(&format!("error: {:?}", e));
            1
        })?;

        log::info("reading zeroes");
        let mut zeroes = [0u8, 16];
        let res = fout.read(&mut zeroes);
        if let Err(err) = res {
            log::error("can't read zeroes from zero file");
            log::error(&format!("error: {}", err));
        }

        log::info("verifying all zeroes are valid");
        for (x, val) in zeroes.iter().enumerate() {
            let val: u8 = *val;
            if val != 0 {
                log::error(&format!("expected zeroes[{}] to be 0, got: {}", x, val));
                return Err(1);
            }
        }
    }
    log::info("closed file");

    log::info("all functions passed basic usage");

    Ok(())
}

fn main() {}
