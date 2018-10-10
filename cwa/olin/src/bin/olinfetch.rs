#![no_main]

extern crate olin;

use olin::{log, panic, runtime, stdio, time};
use std::io::Write;

#[no_mangle]
pub extern "C" fn cwa_main() -> i32 {
    panic::set_hook();

    let mut out = stdio::out();

    let mut rt_name = [0u8; 32];
    let runtime_name = runtime::name_buf(rt_name.as_mut())
        .ok_or_else(|| {
            log::error("Runtime name larger than 32 byte limit");
            1
        }).unwrap();

    write!(out, "CPU:\t\t{}\n", "wasm32");
    write!(
        out,
        "Runtime:\t{} {}.{}\n",
        runtime_name,
        runtime::spec_major(),
        runtime::spec_minor()
    );
    write!(out, "Now:\t\t{}\n", time::now().to_rfc3339());

    0
}
