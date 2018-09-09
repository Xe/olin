extern crate chrono;
extern crate libcwa;

use self::chrono::TimeZone;
use libcwa::log;

/// This tests for https://github.com/CommonWA/cwa-spec/blob/master/ns/time.md
pub extern "C" fn test() -> Result<(), i32> {
    log::info("running ns::time tests");

    let now: i64 = libcwa::time::ts();
    let dt = chrono::Utc.timestamp(now, 0);

    log::info(&format!("ts: {}, dt: {}", now, dt.to_rfc3339()));

    let now = libcwa::time::now();
    log::info(&format!("time::now(): {}", now.to_rfc3339()));

    log::info("ns::time tests passed");
    Ok(())
}
