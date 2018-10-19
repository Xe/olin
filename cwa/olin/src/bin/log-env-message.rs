#![no_main]
#![no_std]

extern crate olin;

use olin::{env, log};

#[no_mangle]
pub extern "C" fn cwa_main() -> i32 {
    match env::get("MESSAGE") {
        Err(env::Error::NotFound) => {
            log::error("MESSAGE doesn't exist");
            return 1;
        }
        Ok(val) => {
            log::info(&val);
        }
        _ => return 2,
    }

    0
}
