extern crate olin;

use olin::{log, Resource};
use std::io::{Read, Write};

/// This tests for https://github.com/CommonWA/cwa-spec/blob/master/schemes/zero.md
pub extern "C" fn test() -> Result<(), i32> {
    log::info("running scheme::zero tests");

    let mut fout: Resource = Resource::open("zero://").map_err(|e| {
        log::error(&format!("couldn't open: {:?}", e));
        1
    })?;

    fout.write(b"so ancient they are swarmed by truth.")
        .map_err(|e| {
            log::error(&format!("couldn't write: {:?}", e));
            1
        });

    let mut inp = [0u8; 16];
    fout.read(&mut inp).map_err(|e| {
        log::error(&format!("couldn't read: {:?}", e));
        1
    });

    log::info("verifying all zeroes are valid");
    for (x, val) in inp.iter().enumerate() {
        let val: u8 = *val;
        if val != 0 {
            log::error(&format!("expected zeroes[{}] to be 0, got: {}", x, val));
            return Err(1);
        }
    }

    log::info("flushing zero file");
    fout.flush().map_err(|e| {
        log::error(&format!("error: {:?}", e));
        1
    })?;

    log::info("scheme::zero tests passed");
    Ok(())
}
