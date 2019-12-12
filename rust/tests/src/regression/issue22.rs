extern crate olin;

use olin::{env, log};

pub extern "C" fn test() -> Result<(), i32> {
    log::info("testing for issue 22: https://github.com/Xe/olin/issues/22");

    log::info("look for variable that does not exist");
    match env::get("DOES_NOT_EXIST") {
        Err(env::Error::NotFound) => log::info("this does not exist! :D"),
        Ok(_) => {
            log::error("DOES_NOT_EXIST exists");

            return Err(1);
        }
        _ => {
            log::error("got other error");

            return Err(2);
        }
    }

    log::info("issue 22 test passed");

    Ok(())
}
