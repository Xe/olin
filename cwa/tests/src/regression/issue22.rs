extern crate libcwa;

use libcwa::{env, log};

pub extern "C" fn test() -> Result<(), i32> {
    log::info("testing for issue 22: https://github.com/Xe/olin/issues/22");

    log::info("look for variable that does not exist");
    match env::get("DOES_NOT_EXIST") {
        None => log::info("this does not exist! :D"),
        Some(_) => {
            log::error("DOES_NOT_EXIST exists");

            return Err(1);
        },
    }

    log::info("issue 22 test passed");

    Ok(())
}
