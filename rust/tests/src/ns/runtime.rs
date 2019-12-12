extern crate olin;

use olin::{log, runtime};

/// This tests for https://github.com/CommonWA/cwa-spec/blob/master/ns/runtime.md
pub extern "C" fn test() -> Result<(), i32> {
    log::info("running ns::runtime tests");

    log::info("expecting spec major=0 and min=0");
    let minor: i32 = runtime::spec_major();
    let major: i32 = runtime::spec_major();

    log::info(&format!("minor: {}, major: {}", minor, major));

    if major != 0 && minor != 0 {
        log::error("version is wrong");
        return Err(1);
    }
    log::info("passed");

    log::info("getting runtime name, should be olin");
    let mut rt_name = [0u8; 32];
    let runtime_name = runtime::name_buf(rt_name.as_mut()).ok_or_else(|| {
        log::error("Runtime name larger than 32 byte limit");
        1
    })?;

    log::info(runtime_name);

    if runtime_name != "olin" {
        log::error("Got runtime name, not olin");
        return Err(1);
    }
    log::info("passed");

    log::info("sleeping");
    runtime::msleep(1);
    log::info("passed");

    log::info("ns::runtime tests passed");
    Ok(())
}
