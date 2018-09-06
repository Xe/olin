use raw;
use std::io;
use std::io::{Read, Write};

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
            raw::resource_open(&url[0], url.len())
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
            raw::resource_close(self.handle);
        }
    }
}

impl Read for Resource {
    fn read(&mut self, out: &mut [u8]) -> io::Result<usize> {
        let len = out.len();

        if len == 0 {
            return Ok(0);
        }

        let ret = unsafe { raw::resource_read(
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

        let ret = unsafe { raw::resource_write(
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
