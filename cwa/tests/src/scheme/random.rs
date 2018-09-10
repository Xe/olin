extern crate olin;

use olin::{log, Resource};
use std::io::{Read, Write};

/// This tests for https://github.com/CommonWA/cwa-spec/blob/master/schemes/random.md
pub extern "C" fn test() -> Result<(), i32> {
    log::info("running scheme::random tests");

    let mut fin: Resource = Resource::open("random://").map_err(|e| {
        log::error(&format!("couldn't open: {:?}", e));
        1
    })?;

    let mut inp = [0u8; 16];
    fin.read(&mut inp).map_err(|e| {
        log::error(&format!("couldn't read: {:?}", e));
        1
    })?;

    log::info("flushing random file");
    fin.flush().map_err(|e| {
        log::error(&format!("error: {:?}", e));
        1
    })?;

    log::info("scheme::random tests passed");
    Ok(())
}
