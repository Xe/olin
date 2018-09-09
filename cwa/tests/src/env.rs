extern crate libcwa;

use libcwa::log;
use libcwa::env;
use std::str;

pub extern "C" fn test() -> Result<(), i32> {
    log::info("running env tests");

    log::info("env[\"MAGIC_CONCH\"] = \"yes\"");
    let envvar_name = "MAGIC_CONCH";
    let mut envvar_val = [0u8; 64];
    let envvar_val = env::get_buf(envvar_name.as_bytes(), &mut envvar_val)
        .map(|s| str::from_utf8(s).expect("envvar wasn't UTF-8"))
        .map_err(|e| {
            log::error(&format!("couldn't get: {:?}", e));
            1
        })?;

    if envvar_val != "yes" {
        log::error(&format!("wanted yes, got: {}", envvar_val));
        return Err(1);
    }
    log::info("passed");

    log::info("env tests passed");
    Ok(())
}
