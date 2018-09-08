extern crate libcwa;

use libcwa::env;
use libcwa::log;
use libcwa::stdio;
use std::io::{Read, Write};
use std::string::String;

fn main() {}

#[no_mangle]
pub extern "C" fn cwa_main() -> i32 {
    log::info("booted");

    match friendly_main() {
        Ok(()) => 0,
        Err(e) => e as i32,
    }
}

#[derive(Debug)]
struct Context {
    pub method: String,
    pub request_uri: String,
    pub body: libcwa::Resource,
}

#[derive(Debug)]
struct Response {
    status: i32,
    body: String,
}

pub fn friendly_main() -> Result<(), i32> {
    let method = getenv("REQUEST_METHOD").unwrap();
    let request_uri = getenv("REQUEST_URI").unwrap();

    let fin = stdio::inp();
    let mut fout = stdio::out();

    let ctx = Context {
        method: method,
        request_uri: request_uri,
        body: fin,
    };

    let resp: Response = respond_to(ctx);
    let set: std::vec::Vec<u8> = serialize(resp);

    let len = fout.write(set.as_slice()).map_err(|e| {
        log::error(&format!("can't write resulting response: {:?}", e));
        1
    }).unwrap();

    if len != set.len() {
        log::warning("wasn't able to write entire response");
        log::warning(&format!("wanted: {}, got: {}", set.len(), len));
    }

    Ok(())
}

fn getenv(name: &str) -> Option<String> {
    let result = env::get(&name.as_bytes());
    if result.is_none()
    {
        ()
    }

    let result = result.unwrap().to_vec();
    let result: String = String::from_utf8(result).unwrap();

    Some(result)
}

fn respond_to(ctx: Context) -> Response {
    match ctx.request_uri.as_str() {
        "/cadey" => Response { status: 200, body: "you are awesome!".to_owned() },
        _        => Response { status: 404, body: "try /cadey".to_owned() }
    }
}

fn serialize(mut response: Response) -> Vec<u8> {
    let mut output = String::new();
    output.push_str("HTTP/1.1 ");
    output.push_str(&format!("{}", response.status));
    output.push_str("\nContent-Type: text/plain\nCetacean-Powered-By: Cadey~#1337\n\n");

    let mut output = output.into_bytes();

    output.append(unsafe { response.body.as_mut_vec() });

    output
}

