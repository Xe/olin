#![allow(unused_must_use)]
#![no_main]
#![feature(start)]

extern crate olin;

mod ns;
mod olintest;
mod regression;
mod scheme;

olin::entrypoint!();

fn main() -> Result<(), std::io::Error> {
    let mut fail_count = 0;

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
            Err(e) => {
                olin::log::error(&format!("test error: {:?}", e));
                fail_count += 1;
            },
        }
    }

    olin::runtime::exit(fail_count);
}
