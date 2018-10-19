#![no_main]

extern crate olin;

use olin::{env, log, panic};

#[no_mangle]
pub extern "C" fn cwa_main() -> i32 {
    panic::set_hook();

    match env::get("MESSAGE") {
        Err(env::Error::NotFound) => {
            log::error("MESSAGE doesn't exist");
            return 1;
        }
        Ok(val) => {
            log::info(&format!("{}", val));
        }
        _ => return 2,
    }

    0
}
