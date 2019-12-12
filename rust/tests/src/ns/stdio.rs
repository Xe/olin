extern crate olin;

use olin::{log, stdio};
use std::io::{Read, Write};

pub extern "C" fn test() -> Result<(), i32> {
    log::info("running ns::stdio tests");

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

    log::info("ns::stdio tests passed");
    Ok(())
}
