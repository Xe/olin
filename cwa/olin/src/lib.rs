#![allow(warnings)]

use std::io::{self, Read, Write};

#[macro_use]
pub mod macros;

pub mod panic;

pub mod sys {
    extern {
        pub fn log_write(level: i32, text_ptr: *const u8, text_len: usize);
        pub fn env_get(
            key_ptr: *const u8,
            key_len: usize,
            value_buf_ptr: *mut u8,
            value_buf_len: usize,
        ) -> i32;
        pub fn runtime_spec_major() -> i32;
        pub fn runtime_spec_minor() -> i32;
        pub fn runtime_name(out_ptr: *mut u8, out_len: usize) -> i32;
        pub fn runtime_msleep(ms: i32);

        pub fn startup_arg_len() -> i32;
        pub fn startup_arg_at(id: i32, out_ptr: *mut u8, out_len: usize) -> i32;

        pub fn resource_open(url_ptr: *const u8, url_len: usize) -> i32;
        pub fn resource_read(id: i32, data_ptr: *mut u8, data_len: usize) -> i32;
        pub fn resource_write(id: i32, data_ptr: *const u8, data_len: usize) -> i32;
        pub fn resource_close(id: i32);
        pub fn resource_flush(id: i32) -> i32;

        pub fn time_now() -> i64;

        pub fn io_get_stdin() -> i32;
        pub fn io_get_stdout() -> i32;
        pub fn io_get_stderr() -> i32;

        pub fn random_i32() -> i32;
        pub fn random_i64() -> i64;
    }
}

mod err {
    use std::io;

    pub const UNKNOWN: i32 = -1;
    pub const INVALID_ARGUMENT: i32 = -2;
    pub const PERMISSION_DENIED: i32 = -3;
    pub const NOT_FOUND: i32 = -4;

    #[repr(i32)]
    #[derive(Debug)]
    pub enum Error {
        Unknown = UNKNOWN,
        InvalidArgument = INVALID_ARGUMENT,
        PermissionDenied = PERMISSION_DENIED,
        NotFound = NOT_FOUND,
    }
    impl Error {
        pub fn check(n: i32) -> Result<i32, Error> {
            match n {
                n if n >= 0 => Ok(n),
                INVALID_ARGUMENT => Err(Error::InvalidArgument),
                PERMISSION_DENIED => Err(Error::PermissionDenied),
                NOT_FOUND => Err(Error::NotFound),
                _ => Err(Error::Unknown),
            }
        }
    }
    pub fn check_io(error: i32) -> Result<i32, io::ErrorKind> {
        match error {
            n if n >= 0 => Ok(n),
            INVALID_ARGUMENT => Err(io::ErrorKind::InvalidInput),
            PERMISSION_DENIED => Err(io::ErrorKind::PermissionDenied),
            NOT_FOUND => Err(io::ErrorKind::NotFound),
            _ => Err(io::ErrorKind::Other),
        }
    }
}

pub mod log {
    use super::sys;

    #[repr(i32)]
    #[derive(Debug)]
    pub enum Level {
        Error = 1,
        Warning = 3,
        Info = 6,
    }

    /// Writes a line of log with the specified level to host logger.
    pub fn write(level: Level, text: &str) {
        let text = text.as_bytes();
        unsafe { sys::log_write(level as i32, text.as_ptr(), text.len()) }
    }

    pub fn error(text: &str) {
        write(Level::Error, text)
    }
    pub fn warning(text: &str) {
        write(Level::Warning, text)
    }
    pub fn info(text: &str) {
        write(Level::Info, text)
    }
}

pub mod env {
    use super::{err, sys};

    #[derive(Debug)]
    pub enum Error {
        NotFound,
        TooSmall(u32),
    }
    pub fn get_buf<'a>(key: &[u8], value: &'a mut [u8]) -> Result<&'a mut [u8], Error> {
        let ret = unsafe { sys::env_get(key.as_ptr(), key.len(), value.as_mut_ptr(), value.len()) };
        match ret {
            err::NOT_FOUND => Err(Error::NotFound),
            n if (n as usize) <= value.len() => Ok(&mut value[0..n as usize]),
            n => Err(Error::TooSmall(n as u32)),
        }
    }

    /// Returns the environment variable associated with `key`.
    /// If there is no environment variable with the specified key or the value is not valid
    /// UTF-8, `None` is returned.
    pub fn get(key: &str) -> Option<String> {
        let key = key.as_bytes();
        let mut current_len: usize = 32;
        loop {
            let mut val: Vec<u8> = Vec::with_capacity(current_len);
            let real_len: i32;

            unsafe {
                val.set_len(current_len);
                real_len = sys::env_get(
                    slice_raw_ptr_or_null!(key),
                    key.len(),
                    slice_raw_ptr_or_null_mut!(&mut val),
                    val.len()
                );
            }

            if real_len < 0 {
                return None;
            }

            let real_len = real_len as usize;
            if real_len > current_len {
                current_len = real_len;
                continue;
            }

            unsafe {
                val.set_len(real_len);
            }

            return match String::from_utf8(val) {
                Ok(v) => Some(v),
                Err(_) => None
            };
        }
    }
}

pub mod runtime {
    use super::{err, sys};

    pub fn spec_major() -> i32 {
        unsafe { sys::runtime_spec_major() }
    }
    pub fn spec_minor() -> i32 {
        unsafe { sys::runtime_spec_minor() }
    }

    pub fn name_buf(out: &mut [u8]) -> Option<&mut str> {
        let ret = unsafe { sys::runtime_name(out.as_mut_ptr(), out.len()) };
        match ret {
            err::INVALID_ARGUMENT => None,
            len => {
                let out = &mut out[0..len as usize];
                Some(unsafe { ::std::str::from_utf8_unchecked_mut(out) })
            }
        }
    }
    pub fn name() -> String {
        const MAX_LEN: usize = 32;
        let mut out = String::with_capacity(MAX_LEN);
        {
            let mut out = unsafe { out.as_mut_vec() };
            let len = unsafe { sys::runtime_name(out.as_mut_ptr(), MAX_LEN) };
            unsafe { out.set_len(len as usize) };
        }
        out
    }

    pub fn msleep(ms: i32) {
        unsafe { sys::runtime_msleep(ms) }
    }
}

pub mod startup {
    use super::{err, sys};
    use std::str;

    pub fn arg_len() -> i32 {
        unsafe { sys::startup_arg_len() }
    }

    pub fn arg_os_at_buf(id: i32, out: &mut [u8]) -> Option<&mut [u8]> {
        let ret = unsafe { sys::startup_arg_at(id, out.as_mut_ptr(), out.len()) };
        match ret {
            err::INVALID_ARGUMENT => None,
            bytes_written => Some(&mut out[0..bytes_written as usize]),
        }
    }
    pub fn arg_os_at(id: i32, capacity: usize) -> Option<Vec<u8>> {
        let mut out = Vec::with_capacity(capacity);
        let ret = unsafe { sys::startup_arg_at(id, out.as_mut_ptr(), out.len()) };
        match ret {
            err::INVALID_ARGUMENT => None,
            bytes_written => {
                unsafe { out.set_len(bytes_written as usize) };
                Some(out)
            }
        }
    }

    pub fn arg_at_buf(id: i32, out: &mut [u8]) -> Option<&mut str> {
        arg_os_at_buf(id, out).map(|s| str::from_utf8_mut(s).expect("arg isn't UTF-8"))
    }
    pub fn arg_at(id: i32, capacity: usize) -> Option<String> {
        arg_os_at(id, capacity).map(|s| String::from_utf8(s).expect("arg isn't UTF-8"))
    }
}

#[derive(Debug)]
pub struct Resource(i32);

impl Resource {
    pub fn open(url: &str) -> Result<Resource, err::Error> {
        let res = unsafe { sys::resource_open(url.as_ptr(), url.len()) };
        err::Error::check(res).map(Resource)
    }

    pub unsafe fn from_raw(handle: i32) -> Resource {
        Resource(handle)
    }
}

impl Drop for Resource {
    fn drop(&mut self) {
        unsafe { sys::resource_close(self.0) };
    }
}

impl Read for Resource {
    fn read(&mut self, out: &mut [u8]) -> io::Result<usize> {
        let len = out.len();
        if len == 0 {
            Ok(0)
        } else {
            let ret = unsafe { sys::resource_read(self.0, out.as_mut_ptr(), len) };
            err::check_io(ret)
                .map(|b| b as usize)
                .map_err(io::Error::from)
        }
    }
}

impl Write for Resource {
    fn write(&mut self, data: &[u8]) -> io::Result<usize> {
        let len = data.len();
        if len == 0 {
            Ok(0)
        } else {
            let ret = unsafe { sys::resource_write(self.0, data.as_ptr(), len) };
            err::check_io(ret)
                .map(|b| b as usize)
                .map_err(io::Error::from)
        }
    }

    fn flush(&mut self) -> io::Result<()> {
        let ret: i32 = unsafe { sys::resource_flush(self.0) };
        if ret == 0 {
            return Ok(());
        }

        return Err(err::check_io(ret)
                   .map_err(io::Error::from)
                   .unwrap_err());
    }
}

pub mod time {
    extern crate chrono;

    use time::chrono::TimeZone;

    pub fn now() -> chrono::DateTime<chrono::Utc> {
        return chrono::Utc.timestamp(ts(), 0);
    }

    pub fn ts() -> i64 {
        unsafe { ::sys::time_now() }
    }
}

pub mod stdio {
    pub fn inp() -> ::Resource {
        unsafe { ::Resource::from_raw(::sys::io_get_stdin()) }
    }

    pub fn out() -> ::Resource {
        unsafe { ::Resource::from_raw(::sys::io_get_stdout()) }
    }

    pub fn err() -> ::Resource {
        unsafe { ::Resource::from_raw(::sys::io_get_stderr()) }
    }
}

pub mod random {
    pub fn i31() -> i32 {
        unsafe { ::sys::random_i32() }
    }

    pub fn i63() -> i64 {
        unsafe { ::sys::random_i64() }
    }
}
