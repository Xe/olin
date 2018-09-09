extern crate http;
extern crate httparse;
extern crate std;

use std::io::{self, ErrorKind, Read, Write};
use std::string::String;
use std::vec::Vec;

pub fn transport<'a>(req: http::Request<&'a mut std::vec::Vec<u8>>, resp_body: &'a mut std::vec::Vec<u8>) -> Result<(http::Response<()>), io::Error> {
    let mut fout: ::Resource =
        ::Resource::open("https://").map_err(|e| {
            ::log::error(&format!("http: couldn't open {:?}", e));
            io::Error::new(ErrorKind::Other, ::err::Error::Unknown)
        })?;

    let req = serialize_req(req);

    fout.write(req.as_slice()).map_err(|e| {
        ::log::error(&format!("http: couldn't write: {:?}", e));
        e
    })?;

    fout.flush().map_err(|e| {
        ::log::error(&format!("http: couldn't flush: {:?}", e));
        e
    })?;

    let mut headers = [httparse::EMPTY_HEADER; 64];
    let mut resp = httparse::Response::new(&mut headers);
    let mut resp_code: u16 = 200;
    let mut response = http::Response::builder();

    while true {
        let mut resp_data: [u8; 2048] = [0u8; 2048];

        fout.read(&mut resp_data).map_err(|e| {
            ::log::error(&format!("http: couldn't read: {:?}", e));
            e
        })?;

        let res = &mut resp.parse(&resp_data).map_err(|e| {
            ::log::error(&format!("http: couldn't parse response: {:?}", e));
            io::Error::new(ErrorKind::Other, e)
        })?;

        if !res.is_partial() {
            break;
        } else {
            match resp.code {
                None => {},
                Some(code) => resp_code = code,
            }
        }
    }

    for x in 0..resp.headers.len() {
        let hdr = resp.headers[x];

        if hdr.name != "" {
            let response = response.header(hdr.name, std::string::String::from_utf8_lossy(hdr.value).into_owned());
        }
    }

    Ok(response.status(resp_code).body(()).unwrap())
}

fn serialize_req<'a>(mut req: http::Request<&'a mut std::vec::Vec<u8>>) -> Vec<u8> {
    let mut output = String::new();
    let mut body = req.body_mut().clone();
    let method = req.method().as_str();
    let path = req.uri().path();

    output.push_str(method);
    output.push_str(" ");
    output.push_str(path);
    output.push_str(" HTTP/1.1\r\n");

    for (key, value) in req.headers().iter() {
        output.push_str(&format!("{:?}: {:?}", key, value));
    }

    output.push_str("\n\n");

    let mut output = output.into_bytes();

    output.append(&mut body);

    output
}
