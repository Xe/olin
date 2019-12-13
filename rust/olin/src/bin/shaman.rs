#![no_main]
#![feature(start)]

extern crate olin;

use std::io::Write;

olin::entrypoint!();

fn main() -> Result<(), std::io::Error> {
    let bytes = include_bytes!("shaman.aa");
    let mut out = olin::stdio::out();

    out.write(bytes)
        .map_err(|e| {
            olin::log::error(&format!("can't write to stdout: {:?}", e));
            1
        }).unwrap();

    Ok(())
}
