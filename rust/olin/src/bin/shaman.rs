#![no_main]

extern crate olin;

use std::io::Write;

#[no_mangle]
pub extern "C" fn cwa_main() -> i32 {
    olin::panic::set_hook();

    let bytes = include_bytes!("shaman.aa");
    let mut out = olin::stdio::out();

    out.write(bytes)
        .map_err(|e| {
            olin::log::error(&format!("can't write to stdout: {:?}", e));
            1
        }).unwrap();

    0
}
