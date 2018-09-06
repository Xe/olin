use std::cell::RefCell;

use resource::Resource;
use raw;

thread_local! {
    static STDIN: RefCell<Resource> = RefCell::new(
        unsafe { Resource::from_raw(raw::io_get_stdin()) }
    );
    static STDOUT: RefCell<Resource> = RefCell::new(
        unsafe { Resource::from_raw(raw::io_get_stdout()) }
    );
    static STDERR: RefCell<Resource> = RefCell::new(
        unsafe { Resource::from_raw(raw::io_get_stderr()) }
    );
}

pub fn with_stdin<F: FnOnce(&mut Resource) -> T, T>(f: F) -> T {
    STDIN.with(|h| f(&mut *h.borrow_mut()))
}

pub fn with_stdout<F: FnOnce(&mut Resource) -> T, T>(f: F) -> T {
    STDOUT.with(|h| f(&mut *h.borrow_mut()))
}

pub fn with_stderr<F: FnOnce(&mut Resource) -> T, T>(f: F) -> T {
    STDERR.with(|h| f(&mut *h.borrow_mut()))
}
