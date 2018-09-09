#![allow(unused_must_use)]

extern crate libcwa;

mod env;
mod http;
mod time;
mod resource;
mod runtime;
mod startup;
mod stdio;

#[no_mangle]
pub extern "C" fn cwa_main() -> i32 {
    let funcs = [
        env::test,
        http::test,
        time::test,
        resource::test,
        runtime::test,
        startup::test,
        stdio::test,
    ];

    for x in 0..funcs.len() {
        match funcs[x]() {
            Ok(()) => {},
            Err(e) => return e as i32,
        }
    }

    0
}

fn main() {}
