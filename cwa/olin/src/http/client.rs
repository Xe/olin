extern crate http;
extern crate httparse;
extern crate std;

use std::io::{self, ErrorKind, Read, Write};
use std::string::String;
use std::vec::Vec;

pub fn transport<'a>(
    req: http::Request<&'a mut Vec<u8>>,
    _resp_body: &'a mut Vec<u8>,
) -> Result<(http::Response<()>), io::Error> {
    let mut fout: crate::Resource = crate::Resource::open("https://").map_err(|e| {
        crate::log::error(&format!("http: couldn't open {:?}", e));
        io::Error::new(ErrorKind::Other, crate::err::Error::Unknown)
    })?;

    let req = serialize_req(req);

    fout.write(req.as_slice()).map_err(|e| {
        crate::log::error(&format!("http: couldn't write: {:?}", e));
        e
    })?;

    fout.flush().map_err(|e| {
        crate::log::error(&format!("http: couldn't flush: {:?}", e));
        e
    })?;

    let mut resp_code: u16 = 200;
    let mut response = http::Response::builder();

    let mut resp_data: Vec<u8> = Vec::new();
    resp_data.resize(4096, 0u8);

    let mut data_pos = 0;
    loop {
        if data_pos == resp_data.len() {
            let new_len = resp_data.len() * 2;
            resp_data.resize(new_len, 0u8);
        }

        data_pos += fout.read(&mut resp_data[data_pos..]).map_err(|e| {
            crate::log::error(&format!("http: couldn't read: {:?}", e));
            e
        })?;

        // These must be defined here because they contain pointers into resp_data
        let mut headers = [httparse::EMPTY_HEADER; 64];
        let mut resp = httparse::Response::new(&mut headers);

        let res = resp.parse(&resp_data).map_err(|e| {
            crate::log::error(&format!("http: couldn't parse response: {:?}", e));
            io::Error::new(ErrorKind::Other, e)
        })?;

        // No need to repeat any work until parsing is complete
        if res.is_partial() {
            continue;
        }

        match resp.code {
            None => {}
            Some(code) => resp_code = code,
        }

        for hdr in resp.headers {
            if hdr.name != "" {
                response.header(
                    hdr.name,
                    std::string::String::from_utf8_lossy(hdr.value).into_owned(),
                );
            }
        }
        break;
    }

    Ok(response.status(resp_code).body(()).unwrap())
}

fn serialize_req<'a>(mut req: http::Request<&'a mut std::vec::Vec<u8>>) -> Vec<u8> {
    let mut output = String::new();
    let mut body = req.body_mut().clone();
    let method = req.method().as_str();
    let host = req.uri().host();
    let path = req.uri().path();

    output.push_str(method);
    output.push_str(" ");
    output.push_str(path);
    output.push_str(" HTTP/1.1\r\n");
    output.push_str(&format!("Host: {}\n", host.unwrap()));

    for (key, value) in req.headers().iter() {
        output.push_str(&format!("{}: {}\n", key.as_str(), value.to_str().unwrap()));
    }

    output.push_str("\n\n");

    let mut output = output.into_bytes();

    output.append(&mut body);

    output
}
