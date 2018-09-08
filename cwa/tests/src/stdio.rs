extern crate libcwa;

use libcwa::*;
use std::io::{self, Read, Write};

pub fn test() -> Result<(), i32> {
    log::info("stdout");
    {
        let mut fout = stdio::out();
        fout.write(b"Hi there\n");
    }

    log::info("stderr");
    {
        let mut fout = stdio::err();
        fout.write(b"Hi there\n");
    }

    log::info("stdin");
    {
        let mut fin = stdio::inp();
        let mut resp: [u8; 16] = [0u8; 16];
        fin.read(&mut resp).map_err(|e| {
            log::error(&format!("can't read from stdin: {:?}", e));
        });
    }

    log::info("stdio tests passed");
    Ok(())
}
