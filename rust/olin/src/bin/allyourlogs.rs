#![no_main]
#![feature(start)]

extern crate olin;

use olin::{log, entrypoint};

entrypoint!();

fn main() -> Result<(), std::io::Error> {
    let string = "hi";
    log::error(&string);
    log::warning(&string);
    log::info(&string);

    Ok(())
}
