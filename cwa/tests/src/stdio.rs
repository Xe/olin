extern crate libcwa;

use libcwa::*;
use std::io::{self, Read, Write};

pub fn test() -> Result<(), i32> {
    log::info("stdout");
    {
        let mut fout = stdio::out();
        fout.write(b"Hi there\n").map_err(|e| {
            log::error(&format!("can't write to stdout: {:?}", e));
            1
        });
    }

    log::info("stderr");
    {
        let mut fout = stdio::err();
        fout.write(b"Hi there\n").map_err(|e| {
            log::error(&format!("can't write to stderr: {:?}", e));
            1
        });
    }

    log::info("stdin");
    {
        let mut fin = stdio::inp();
        let mut resp: [u8; 16] = [0u8; 16];
        fin.read(&mut resp).map_err(|e| {
            log::error(&format!("can't read from stdin: {:?}", e));
            1
        });
    }

    log::info("stdio tests passed");
    Ok(())
}
