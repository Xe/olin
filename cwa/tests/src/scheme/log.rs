extern crate libcwa;

use std::io::Write;
use libcwa::{log, Resource};

pub extern "C" fn test() -> Result<(), i32> {
    log::info("running scheme::log tests");

    log::info("log://");
    {
        let mut fout: Resource =
            Resource::open("log://").map_err(|e| {
                log::error(&format!("couldn't open: {:?}", e));
                1
            })?;

        fout.write(b"I am listening for a sound beyond sound").map_err(|e| {
            log::error(&format!("couldn't write log: {:?}", e));
            1
        });

        fout.flush().map_err(|e| {
            log::error(&format!("couldn't flush: {:?}", e));
            1
        });
    }
    log::info("passed");

    log::info("log://?prefix=inner");
    {
        let mut fout: Resource =
            Resource::open("log://?prefix=inner").map_err(|e| {
                log::error(&format!("couldn't open: {:?}", e));
                1
            })?;

        fout.write(b"that stalks the nightland of my dreams").map_err(|e| {
            log::error(&format!("couldn't write log: {:?}", e));
            1
        });

        fout.flush().map_err(|e| {
            log::error(&format!("couldn't flush: {:?}", e));
            1
        });
    }
    log::info("passed");

    log::info("scheme::log tests passed");
    Ok(())
}
