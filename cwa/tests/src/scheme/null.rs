extern crate libcwa;

use libcwa::{log, Resource};
use std::io::{Read, Write};

pub extern "C" fn test() -> Result<(), i32> {
    log::info("running scheme::null tests");

    let mut fout: Resource =
        Resource::open("null://").map_err(|e| {
            log::error(&format!("couldn't open: {:?}", e));
            1
        })?;

    fout.write(b"entering rooms of fossil-light").map_err(|e| {
        log::error(&format!("couldn't write: {:?}", e));
        1
    });

    let mut inp = [0u8; 16];
    fout.read(&mut inp).map_err(|e| {
        log::error(&format!("couldn't read: {:?}", e));
        1
    });

    log::info("flushing null file");
    fout.flush().map_err(|e| {
        log::error(&format!("error: {:?}", e));
        1
    })?;

    log::info("scheme::null tests passed");
    Ok(())
}
