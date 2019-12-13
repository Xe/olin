/// Main modules to perform the low-level operations
/// described in the Olin specs. In particular the basic
/// "actions": open, read, write, close, flush on files.
///
/// Modules:
/// * sys
/// * err
/// * log
/// * env
/// * runtime
/// * startup
/// * time
/// * stdio
/// * random
///
/// The ```Resource``` struct defines a high level wrapper
/// to the more low level IO calls in sys.
extern crate chrono;

use std::io::{self, Read, Write};

pub mod http;
pub mod panic;

// Replace the allocator

extern crate wee_alloc;

// Use `wee_alloc` as the global allocator.
#[global_allocator]
static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;

/// system-level operations, interface and FFI access:
/// * setup and low-level operations
/// * bindings to environment variables
/// * bindings for system calls
pub mod sys {
    extern "C" {
        // https://github.com/CommonWA/cwa-spec/blob/master/ns/log.md#write
        pub fn log_write(level: i32, text_ptr: *const u8, text_len: usize);

        // https://github.com/CommonWA/cwa-spec/blob/master/ns/env.md#get
        pub fn env_get(
            key_ptr: *const u8,
            key_len: usize,
            value_buf_ptr: *mut u8,
            value_buf_len: usize,
        ) -> i32;

        pub fn runtime_exit(status: i32) -> !;

        // https://github.com/CommonWA/cwa-spec/blob/master/ns/runtime.md#spec_major
        pub fn runtime_spec_major() -> i32;
        // https://github.com/CommonWA/cwa-spec/blob/master/ns/runtime.md#spec_minor
        pub fn runtime_spec_minor() -> i32;
        // https://github.com/CommonWA/cwa-spec/blob/master/ns/runtime.md#name
        pub fn runtime_name(out_ptr: *mut u8, out_len: usize) -> i32;
        // https://github.com/CommonWA/cwa-spec/blob/master/ns/runtime.md#msleep
        pub fn runtime_msleep(ms: i32);

        // https://github.com/CommonWA/cwa-spec/blob/master/ns/startup.md#arg_len
        pub fn startup_arg_len() -> i32;
        // https://github.com/CommonWA/cwa-spec/blob/master/ns/startup.md#art_at
        pub fn startup_arg_at(id: i32, out_ptr: *mut u8, out_len: usize) -> i32;

        // https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md#open
        pub fn resource_open(url_ptr: *const u8, url_len: usize) -> i32;
        // https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md#read
        pub fn resource_read(id: i32, data_ptr: *mut u8, data_len: usize) -> i32;
        // https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md#write
        pub fn resource_write(id: i32, data_ptr: *const u8, data_len: usize) -> i32;
        // https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md#close
        pub fn resource_close(id: i32);
        // https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md#flush
        pub fn resource_flush(id: i32) -> i32;

        // https://github.com/Xe/cwa-spec/blob/Xe/time/ns/time.md#now
        // This is special-snowflaked at the moment because there's a disagreement
        // about how to do this in the spec at https://github.com/CommonWA/cwa-spec/pull/8
        pub fn time_now() -> i64;

        // https://github.com/CommonWA/cwa-spec/blob/master/ns/random.md#i32
        pub fn random_i32() -> i32;
        // https://github.com/CommonWA/cwa-spec/blob/master/ns/random.md#i64
        pub fn random_i64() -> i64;

        // This api is a nonce.
        // https://github.com/Xe/olin/blob/master/docs/cwa-spec/ns/io.md#io_get_stdin
        pub fn io_get_stdin() -> i32;
        // https://github.com/Xe/olin/blob/master/docs/cwa-spec/ns/io.md#io_get_stdout
        pub fn io_get_stdout() -> i32;
        // https://github.com/Xe/olin/blob/master/docs/cwa-spec/ns/io.md#io_get_stderr
        pub fn io_get_stderr() -> i32;
    }
}

mod err {
    use std::error;
    use std::fmt;
    use std::io;

    // https://github.com/CommonWA/cwa-spec/blob/master/errors.md
    pub const UNKNOWN: i32 = -1;
    pub const INVALID_ARGUMENT: i32 = -2;
    pub const PERMISSION_DENIED: i32 = -3;
    pub const NOT_FOUND: i32 = -4;
    pub const EOF: i32 = -5;

    /// An error abstraction, all of the following values are copied from the spec at:
    /// https://github.com/CommonWA/cwa-spec/blob/master/errors.md
    #[repr(i32)]
    #[derive(Debug)]
    pub enum Error {
        Unknown = UNKNOWN,
        InvalidArgument = INVALID_ARGUMENT,
        PermissionDenied = PERMISSION_DENIED,
        NotFound = NOT_FOUND,
        EOF = EOF,
    }

    impl Error {
        pub fn check(n: i32) -> Result<i32, Error> {
            match n {
                n if n >= 0 => Ok(n),
                INVALID_ARGUMENT => Err(Error::InvalidArgument),
                PERMISSION_DENIED => Err(Error::PermissionDenied),
                NOT_FOUND => Err(Error::NotFound),
                EOF => Err(Error::EOF),
                _ => Err(Error::Unknown),
            }
        }
    }

    /// pretty-print the error using Debug-derived code
    /// XXX(Xe): is this a mistake?
    impl self::fmt::Display for Error {
        fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
            write!(f, "{:?}", self)
        }
    }

    impl self::error::Error for Error {}

    pub fn check_io(error: i32) -> Result<i32, io::ErrorKind> {
        match error {
            n if n >= 0 => Ok(n),
            INVALID_ARGUMENT => Err(io::ErrorKind::InvalidInput),
            PERMISSION_DENIED => Err(io::ErrorKind::PermissionDenied),
            NOT_FOUND => Err(io::ErrorKind::NotFound),
            EOF => Err(io::ErrorKind::UnexpectedEof),
            _ => Err(io::ErrorKind::Other),
        }
    }
}

pub mod log {
    use super::sys;

    /// See Level enum defined in https://github.com/CommonWA/cwa-spec/blob/master/ns/log.md#write
    #[repr(i32)]
    #[derive(Debug)]
    pub enum Level {
        Error = 1,
        Warning = 3,
        Info = 6,
    }

    /// Writes a line of text with the specified level to the host logger.
    pub fn write(level: Level, text: &str) {
        let text = text.as_bytes();
        unsafe { sys::log_write(level as i32, text.as_ptr(), text.len()) }
    }

    /// Convenience wrapper for the error level.
    pub fn error(text: &str) {
        write(Level::Error, text)
    }

    /// Convenience wrapper for the warning level.
    pub fn warning(text: &str) {
        write(Level::Warning, text)
    }

    /// Convenience wrapper for the info level.
    pub fn info(text: &str) {
        write(Level::Info, text)
    }
}

/// Access to environment variables
pub mod env {
    use super::{err, sys};

    /// https://github.com/CommonWA/cwa-spec/blob/master/ns/env.md#get
    #[derive(Debug)]
    pub enum Error {
        NotFound,
        TooSmall(u32),
    }

    pub fn get_buf<'a>(key: &[u8], value: &'a mut [u8]) -> Result<&'a mut [u8], Error> {
        let ret = unsafe { sys::env_get(key.as_ptr(), key.len(), value.as_mut_ptr(), value.len()) };
        match ret {
            err::NOT_FOUND => Err(Error::NotFound),
            n if (n as usize) <= value.len() => {
                Ok(unsafe { value.get_unchecked_mut(0..n as usize) })
            }
            n => Err(Error::TooSmall(n as u32)),
        }
    }

    /// Return the value of an environment variable or the reason why we can't.
    pub fn get(key: &str) -> Result<String, Error> {
        const ENVVAR_LEN: usize = 4096; // at work we don't see longer than 128 usually
        let key = key.as_bytes();
        let mut val = [0u8; ENVVAR_LEN];
        let val = get_buf(&key, &mut val).map_err(|e| {
            crate::log::error(&format!("can't get envvar: {:?}", e));
            e
        })?;

        let mut s: String = ::std::string::String::from("");
        s.push_str(::std::str::from_utf8(&val).unwrap());
        Ok(s)
    }
}

pub mod runtime {
    use super::{err, sys};

    pub fn exit(status: i32) -> ! {
        unsafe { sys::runtime_exit(status) }
    }

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
                let out = unsafe { out.get_unchecked_mut(0..len as usize) };
                Some(unsafe { ::std::str::from_utf8_unchecked_mut(out) })
            }
        }
    }

    pub fn name() -> String {
        const MAX_LEN: usize = 32;
        let mut out = Vec::with_capacity(MAX_LEN);
        {
            let len = unsafe { sys::runtime_name(out.as_mut_ptr(), MAX_LEN) };
            unsafe { out.set_len(len as usize) };
        }
        unsafe { String::from_utf8_unchecked(out) }
    }

    pub fn msleep(ms: i32) {
        unsafe { sys::runtime_msleep(ms) }
    }
}

/// used to fetch argv/argc
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
            bytes_written => Some(unsafe { out.get_unchecked_mut(0..bytes_written as usize) }),
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

///
/// Resources are maybe readable, maybe writable, maybe flushable and closable
/// streams of bytes. A resource is backed by implementation in the containing
/// runtime based on the URL scheme used in the `open` call.
///
/// Resources do not always have to be backed by things that exist. Resources
/// may contain the transient results of API calls. Once the data in a Resource
/// has been read, that may be the only time that data can ever be read.
///
/// Please take this into consideration when consuming the results of Resources.
///
/// The implementation of this type uses the functions at https://github.com/CommonWA/cwa-spec/blob/master/ns/resource.md
///
#[derive(Debug)]
pub struct Resource(i32);

impl Resource {
    pub fn open(url: &str) -> Result<Resource, err::Error> {
        let res = unsafe { sys::resource_open(url.as_ptr(), url.len()) };
        err::Error::check(res).map(Resource)
    }

    /// For the few times you _actually_ do know the file descriptor.
    ///
    /// The CommonWA spec doesn't mandate any properties about file descriptors.
    /// The implementation in Olin uses random file descriptors, meaning that
    /// unless the stdio module is used, this function should probably never be
    /// called unless you really know what you are doing.
    ///
    /// Passing an invalid file descriptor to this function will result in a
    /// handle that returns errors on every call to it. This behavior may be
    /// useful.
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
            Ok(())
        } else {
            err::check_io(ret).map(|_| ()).map_err(io::Error::from)
        }
    }
}

pub mod time {
    use super::sys;
    use chrono::{self, TimeZone};

    pub fn now() -> chrono::DateTime<chrono::Utc> {
        chrono::Utc.timestamp(ts(), 0)
    }

    pub fn ts() -> i64 {
        unsafe { sys::time_now() }
    }
}

/// Bindings from semantic I/O targets to Resources
pub mod stdio {
    use super::Resource;
    pub fn inp() -> Resource {
        unsafe { Resource::from_raw(crate::sys::io_get_stdin()) }
    }

    pub fn out() -> Resource {
        unsafe { Resource::from_raw(crate::sys::io_get_stdout()) }
    }

    pub fn err() -> Resource {
        unsafe { Resource::from_raw(crate::sys::io_get_stderr()) }
    }
}

pub mod random {
    use super::sys;
    pub fn i31() -> i32 {
        unsafe { sys::random_i32() }
    }

    pub fn i63() -> i64 {
        unsafe { sys::random_i64() }
    }
}

#[macro_export]
macro_rules! entrypoint {
    () => {
        #[no_mangle]
        #[start]
        extern "C" fn _start() {
            olin::panic::set_hook();

            if let Err(e) = main() {
                olin::log::error(format!("Application error: {:}", e).as_str());
                olin::runtime::exit(1);
            }

            olin::runtime::exit(0);
        }
    }
}
