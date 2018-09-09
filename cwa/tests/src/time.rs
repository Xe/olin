extern crate chrono;
extern crate libcwa;

use time::chrono::TimeZone;
use self::libcwa::*;

/// This tests for https://github.com/CommonWA/cwa-spec/blob/master/ns/time.md
pub extern "C" fn test() -> Result<(), i32> {
    let now: i64 = time::ts();
    let dt = chrono::Utc.timestamp(now, 0);

    log::info(&format!("ts: {}, dt: {}", now, dt.to_rfc3339()));

    let now = time::now();
    log::info(&format!("time::now(): {}", now.to_rfc3339()));

    log::info("time test passed");
    Ok(())
}
