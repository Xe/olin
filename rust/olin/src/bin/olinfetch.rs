#![no_main]
#![feature(start)]

extern crate olin;

use olin::{log, runtime, stdio, time, entrypoint};
use std::io::Write;

entrypoint!();

fn main() -> Result<(), std::io::Error> {
    let mut rt_name = [0u8; 32];
    let runtime_name = runtime::name_buf(rt_name.as_mut())
        .ok_or_else(|| {
            log::error("Runtime name larger than 32 byte limit");
            1
        })
        .unwrap();

    let mut out = stdio::out();

    write!(out, "CPU:\t\t{}\n", "wasm32").expect("write to work");
    write!(
        out,
        "Runtime:\t{} {}.{}\n",
        runtime_name,
        runtime::spec_major(),
        runtime::spec_minor()
    )
        .expect("write to work");
    write!(out, "Now:\t\t{}\n", time::now().to_rfc3339()).expect("write to work");
    Ok(())
}
