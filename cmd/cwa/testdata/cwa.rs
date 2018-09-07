#![allow(warnings)]

use std::io;
use std::io::{Read, Write};

//#[link(name = "cwa", wasm_import_module = "cwa")]
extern "C" {
    pub fn log_write(level: i32, text_ptr: *const u8, text_len: usize);
    pub fn env_get(
        key_ptr: *const u8, key_len: usize,
        value_buf_ptr: *mut u8, value_buf_len: usize
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
}

#[repr(i32)]
pub enum Level {
    Error = 1,
    Warning = 3,
    Info = 6
}

/// Writes a line of log with the specified level to host logger.
pub fn log(level: Level, text: &str) {
    let text = text.as_bytes();

    unsafe {
        log_write(
            level as i32,
            text.as_ptr(),
            text.len()
        );
    }
}

pub struct Resource {
    handle: i32
}

impl Resource {
    pub fn open(url: &str) -> Option<Resource> {
        let url = url.as_bytes();
        if url.len() == 0 {
            return None;
        }

        let handle = unsafe {
            resource_open(&url[0], url.len())
        };
        if handle < 0 {
            None
        } else {
            Some(Resource {
                handle: handle
            })
        }
    }

    pub unsafe fn from_raw(handle: i32) -> Resource {
        // TODO: Deal with invalid handles (< 0)
        Resource {
            handle: handle
        }
    }
}

impl Drop for Resource {
    fn drop(&mut self) {
        unsafe {
            resource_close(self.handle);
        }
    }
}

impl Read for Resource {
    fn read(&mut self, out: &mut [u8]) -> io::Result<usize> {
        let len = out.len();

        if len == 0 {
            return Ok(0);
        }

        let ret = unsafe { resource_read(
            self.handle,
            &mut out[0],
            len
        ) };

        if ret < 0 {
            Err(io::Error::from(io::ErrorKind::Other))
        } else {
            Ok(ret as usize)
        }
    }
}

impl Write for Resource {
    fn write(&mut self, data: &[u8]) -> io::Result<usize> {
        let len = data.len();

        if len == 0 {
            return Ok(0);
        }

        let ret = unsafe { resource_write(
            self.handle,
            &data[0],
            len
        ) };
        if ret < 0 {
            Err(io::Error::from(io::ErrorKind::Other))
        } else {
            Ok(ret as usize)
        }
    }

    fn flush(&mut self) -> io::Result<()> {
        // TODO
        Ok(())
    }
}
