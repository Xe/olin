extern crate libcwa;

use libcwa::env;
use libcwa::log;
use libcwa::startup;
use libcwa::stdio;
use libcwa::runtime;
use std::io::Write;
use std::string::String;

fn main() {}

#[no_mangle]
pub extern "C" fn cwa_main() -> i32 {
    libcwa::panic::set_hook();

    match friendly_main() {
        Ok(()) => 0,
        Err(e) => e as i32,
    }
}

#[derive(Debug)]
struct Context<'a> {
    pub method: String,
    pub request_uri: String,
    pub body: &'a libcwa::Resource,
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

    let ctx: Context = Context {
        method: method,
        request_uri: request_uri,
        body: &fin,
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

fn runtime_info(ctx: Context) -> Response {
    let mut result = String::new();
    result.push_str("Hello, I am served from Rust compiled to wasm32-unknown-unknown.\n\n");
    result.push_str("I know the following about the environment I am running in:\n");

    let minor: i32 = runtime::spec_minor();
    let major: i32 = runtime::spec_major();
    let rt_name: String = runtime::name();

    result.push_str(&format!(" - I am running in {}, which implements version {}.{} of the\n   CommonWebAssembly API.\n", rt_name, major, minor));

    result.push_str(&format!(" - I think the current time is {}\n", libcwa::time::now()));

    result.push_str(&format!(" - RUN_ID:    {:?}\n - WORKER_ID: {:?}\n", getenv("RUN_ID"), getenv("WORKER_ID")));
    result.push_str(&format!(" - Method: {}\n - Request URI: {}\n", ctx.method, ctx.request_uri));

    let argc: i32 = startup::arg_len();

    result.push_str(&format!(" - argc: {}\n", argc));

    for x in 0..argc {
        let mut arg_val = [0u8; 64];
        let arg = startup::arg_at_buf(x, &mut arg_val).ok_or_else(|| {
            log::error(&format!("arg {} missing", x));
            panic!("arg missing");
        }).unwrap();
        result.push_str(&format!(" - arg {}: {}\n", x, arg));
    }

    result.push_str("\nHere is my source code: \n  https://github.com/Xe/olin/blob/master/cwa/cwagi/src/main.rs\n\n");

    result.push_str("If you would like to learn more about this project, please \ntake a look at the following links:\n");
    result.push_str(" - https://github.com/Xe/olin\n");
    result.push_str(" - https://christine.website/blog/olin-1-why-09-1-2018\n");
    result.push_str(" - https://christine.website/blog/olin-2-the-future-09-5-2018\n\n");

    result.push_str("If you know of a Rust HTTP library that lets users hijack\n");
    result.push_str("the transport layer (and can offer code examples, please, I'm\n");
    result.push_str("still new to Rust) so that it could be adapted to Olin's\n");
    result.push_str("native HTTP support, here is its test: \n  https://github.com/Xe/olin/blob/master/cwa/tests/src/scheme/http.rs\n\n");

    result.push_str("If you would like to, please feel free to load test the route\n");
    result.push_str("/cadey. That route does a minimal number of system calls, and\n");
    result.push_str("should make for the best benchmarking results. There is a more\n");
    result.push_str("detailed metrics view at /expvar. Olin does a level of automatic\n");
    result.push_str("load balancing that takes advantage of the Go runtime, so in\n");
    result.push_str("theory the limits of this program approach the limits of how fast\n");
    result.push_str("the code running in the WebAssembly environment is.\n\n");

    result.push_str("And this is just the beginning.\n\n");

    result.push_str("Have a good day and be well, creator.");

    Response {
        status: 200,
        body: result,
    }
}

fn respond_to(ctx: Context) -> Response {
    match ctx.request_uri.as_str() {
        "/cadey" => Response { status: 200, body: "you are awesome!".to_owned() },
        _        => runtime_info(ctx),
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

