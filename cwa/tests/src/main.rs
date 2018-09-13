#![allow(unused_must_use)]

extern crate olin;

mod ns;
mod olintest;
mod regression;
mod scheme;

#[no_mangle]
pub extern "C" fn cwa_main() -> i32 {
    olin::panic::set_hook();

    let funcs = [
        ns::env::test,
        ns::random::test,
        ns::resource::test,
        ns::runtime::test,
        ns::startup::test,
        ns::stdio::test,
        ns::time::test,
        olintest::http::test,
        regression::issue22::test,
        regression::issue37::test,
        regression::issue39::test,
        scheme::http::test,
        scheme::log::test,
        scheme::null::test,
        scheme::random::test,
        scheme::zero::test,
    ];

    for func in &funcs {
        match func() {
            Ok(()) => {}
            Err(e) => return e as i32,
        }
    }

    0
}

fn main() {}
