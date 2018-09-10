extern crate olin;

use olin::{log, startup};

/// This tests for https://github.com/CommonWA/cwa-spec/blob/master/ns/startup.md
pub extern "C" fn test() -> Result<(), i32> {
    log::info("running ns::startup tests");

    log::info("checking argc/argv");
    let argc: i32 = startup::arg_len();
    log::info(&format!("argc: {}", argc));

    for x in 0..argc {
        let mut arg_val = [0u8; 64];
        let arg = startup::arg_at_buf(x, &mut arg_val).ok_or_else(|| {
            log::error(&format!("arg {} missing", x));
            1
        })?;
        log::info(&format!("arg {}: {}", x, arg));
    }
    log::info("passed");

    log::info("ns::startup tests passed");
    Ok(())
}
