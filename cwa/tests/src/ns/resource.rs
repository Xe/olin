extern crate olin;

use olin::{log, Resource};
use std::io::{Read, Write};

/// This tests for https://github.com/CommonWA/cwa-spec/blob/master/ns/runtime.md
pub extern "C" fn test() -> Result<(), i32> {
    log::info("running ns::resource tests");

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
        let mut zeroes = [0u8; 16];
        let res = fout.read(&mut zeroes);
        if let Err(err) = res {
            log::error("can't read zeroes from zero file");
            log::error(&format!("error: {}", err));
            return Err(1);
        }

        log::info("verifying all zeroes are valid");
        for (x, val) in zeroes.iter().enumerate() {
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
    }
    log::info("closed file");

    log::info("ns::resource tests passed");
    Ok(())
}
