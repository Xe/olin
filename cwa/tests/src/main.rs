#![allow(unused_must_use)]

extern crate libcwa;

mod http;
mod ns;

#[no_mangle]
pub extern "C" fn cwa_main() -> i32 {
    let funcs = [
        ns::env::test,
        ns::resource::test,
        ns::runtime::test,
        ns::startup::test,
        ns::stdio::test,
        ns::time::test,
        http::test,
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
