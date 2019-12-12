use std::panic;

pub fn set_hook() {
    panic::set_hook(Box::new(|pi| {
        crate::log::error(&format!("panic: {:?}", pi));
    }));
}
