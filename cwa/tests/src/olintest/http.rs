extern crate http;
extern crate olin;

use olin::log;
use std::vec::Vec;

pub extern "C" fn test() -> Result<(), i32> {
    log::info("running olintest::http tests");

    let mut resp_body = Vec::<u8>::new();
    let mut req_body = Vec::<u8>::new();
    let req = http::Request::builder()
        .uri("https://printerfacts.herokuapp.com")
        .header("User-Agent", "my-awesome-agent/1.0")
        .body(&mut req_body)
        .map_err(|e| {
            log::error(&format!("request error: {:?}", e));
            1
        })?;

    let resp = olin::http::client::transport(req, &mut resp_body).map_err(|e| {
        log::error(&format!("transport error: {:?}", e));
        1
    })?;

    log::info(&format!("status: {:?}", resp.status()));

    log::info("olintest::http tests passed");
    Ok(())
}
