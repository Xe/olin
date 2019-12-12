extern crate olin;

use olin::{log, sys};

pub extern "C" fn test() -> Result<(), i32> {
    log::info("testing for issue 37: https://github.com/Xe/olin/issues/37");

    let mut envvar_val = [0u8; 64];
    let resp = unsafe { sys::env_get("BUTTS".as_bytes().as_ptr(), 5, envvar_val.as_mut_ptr(), 64) };
    log::info(&format!("env_get(\"BUTTS\", 5) => {}", resp));
    match resp {
        -4 => log::info("got expected value"),
        _ => {
            log::error("expected -4");
            return Err(1);
        }
    }

    log::info("issue 37 test passed");

    Ok(())
}
