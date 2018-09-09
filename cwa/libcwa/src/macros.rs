#[macro_export]
macro_rules! main {
    ($body:block) => {
        #[no_mangle]
        pub extern "C" fn cwa_main() -> i32 {
            { $body }
        }
    }
}

#[macro_export]
macro_rules! slice_raw_ptr_or_null {
    ($t:expr) => {
        if $t.len() == 0 {
            ::std::ptr::null()
        } else {
            &$t[0] as *const _
        }
    }
}

#[macro_export]
macro_rules! slice_raw_ptr_or_null_mut {
    ($t:expr) => {
        if $t.len() == 0 {
            ::std::ptr::null_mut()
        } else {
            &mut $t[0] as *mut _
        }
    }
}
