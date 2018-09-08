extern crate chrono;
extern crate libcwa;

use time::chrono::TimeZone;
use std::io;
use libcwa::*;

pub fn test() -> Result<(), i32> {
    let now: i64 = time::ts();
    let dt = chrono::Utc.timestamp(now, 0);

    log::info(&format!("ts: {}, dt: {}", now, dt.to_rfc3339()));

    let now = time::now();
    log::info(&format!("time::now(): {}", now.to_rfc3339()));

    log::info("time test passed");
    Ok(())
}
