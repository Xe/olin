use std::panic;

pub fn set_hook() {
    panic::set_hook(Box::new(|pi| {
        ::log::error(&format!("panic: {:?}", pi));
    }));
}
