extern crate olin;

use olin::{log, random};

pub extern "C" fn test() -> Result<(), i32> {
    log::info("running ns::random tests");

    log::info(&format!("i31: {}, i63: {}", random::i31(), random::i63()));

    log::info("ns::random tests passed");
    Ok(())
}
