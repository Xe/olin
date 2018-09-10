extern crate httparse;
extern crate olin;

use olin::Resource;
use olin::*;
use std::io::{Read, Write};

pub extern "C" fn test() -> Result<(), i32> {
    log::info("running scheme::http tests");

    let reqd = "GET /404 HTTP/1.1\r\nHost: printerfacts.herokuapp.com\r\nUser-Agent: Bit-banging it in rust\r\n\r\n";
    let mut headers = [httparse::EMPTY_HEADER; 16];
    let mut req = httparse::Request::new(&mut headers);
    log::info("validating HTTP request");
    req.parse(reqd.as_bytes()).map_err(|e| {
        log::error(&format!("can't parse request: {:?}", e));
        1
    });

    log::info("opening https://printerfacts.herokuapp.com");
    let mut fout: Resource = Resource::open("https://printerfacts.herokuapp.com").map_err(|e| {
        log::error(&format!("couldn't open: {:?}", e));
        1
    })?;

    log::info("writing HTTP request");
    fout.write(reqd.as_bytes()).map_err(|e| {
        log::error(&format!("can't write request: {:?}", e));
        1
    });

    log::info("fetching response");
    fout.flush().map_err(|e| {
        log::error(&format!("can't send request to remote server: {:?}", e));
        1
    });

    log::info("reading response");
    let mut resp_data = [0u8; 2048];
    fout.read(&mut resp_data).map_err(|e| {
        log::error(&format!("can't read response: {:?}", e));
        1
    });

    log::info("parsing response");
    let mut headers = [httparse::EMPTY_HEADER; 16];
    let mut resp = httparse::Response::new(&mut headers);
    resp.parse(&resp_data).map_err(|e| {
        log::error(&format!("can't parse response: {:?}", e));
        1
    });

    log::info(&format!(
        "version: {:?}, code: {:?}, reason: {:?}",
        resp.version, resp.code, resp.reason
    ));

    log::info("scheme::http tests passed");
    Ok(())
}
