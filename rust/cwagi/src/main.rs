#![no_main]
#![feature(start)]

extern crate olin;

use olin::env;
use olin::log;
use olin::runtime;
use olin::startup;
use olin::stdio;
use std::io::Write;
use std::string::String;

olin::entrypoint!();

fn main() -> Result<(), std::io::Error> {
    olin::runtime::exit(match friendly_main() {
        Ok(()) => 0,
        Err(e) => e as i32,
    });
}

#[no_mangle]
pub extern "C" fn cwa_main() -> i32 {
    olin::panic::set_hook();

    match friendly_main() {
        Ok(()) => 0,
        Err(e) => e as i32,
    }
}

#[derive(Debug)]
struct Context<'a> {
    pub method: String,
    pub request_uri: String,
    pub body: &'a olin::Resource,
}

#[derive(Debug)]
struct Response {
    status: i32,
    body: String,
}

pub fn friendly_main() -> Result<(), i32> {
    let method = env::get("REQUEST_METHOD").map_err(|e| {
        log::error(&format!("error getting REQUEST_METHOD: {:?}", e));
        1
    })?;
    let request_uri = env::get("PATH_INFO").map_err(|e| {
        log::error(&format!("error getting REQUEST_URI: {:?}", e));
        1
    })?;

    let fin = stdio::inp();
    let mut fout = stdio::out();

    let ctx: Context = Context {
        method,
        request_uri,
        body: &fin,
    };

    let resp: Response = respond_to(&ctx);
    let set: std::vec::Vec<u8> = serialize(&resp);

    let len = fout.write(set.as_slice()).map_err(|e| {
        log::error(&format!("can't write resulting response: {:?}", e));
        1
    })?;

    if len != set.len() {
        log::warning("wasn't able to write entire response");
        log::warning(&format!("wanted: {}, got: {}", set.len(), len));
    }

    Ok(())
}

fn runtime_info(ctx: &Context) -> Response {
    let mut result = String::new();
    result.push_str(
        "Hello, I am served from Rust compiled to wasm32-unknown-unknown.\n
I know the following about the environment I am running in:\n",
    );

    let minor: i32 = runtime::spec_minor();
    let major: i32 = runtime::spec_major();
    let rt_name: String = runtime::name();

    result.push_str(&format!(
        " - I am running in {}, which implements version {}.{} of the\n   CommonWebAssembly API.\n",
        rt_name, major, minor
    ));

    result.push_str(&format!(
        " - I think the current time is {}\n",
        olin::time::now()
    ));

    let run_id: String = env::get("RUN_ID")
        .map_err(|e| {
            log::error(&format!("error getting RUN_ID: {:?}", e));
            1
        }).unwrap();
    let worker_id: String = env::get("WORKER_ID")
        .map_err(|e| {
            log::error(&format!("error getting WORKER_ID: {:?}", e));
            1
        }).unwrap();

    result.push_str(&format!(
        " - RUN_ID:    {}\n - WORKER_ID: {}\n",
        run_id, worker_id,
    ));
    result.push_str(&format!(
        " - Method: {}\n - Request URI: {}\n",
        ctx.method, ctx.request_uri
    ));

    let argc: i32 = startup::arg_len();

    result.push_str(&format!(" - argc: {}\n", argc));

    for x in 0..argc {
        let mut arg_val = [0u8; 64];
        let arg = startup::arg_at_buf(x, &mut arg_val)
            .ok_or_else(|| {
                log::error(&format!("arg {} missing", x));
                panic!("arg missing");
            }).unwrap();
        result.push_str(&format!(" - arg {}: {}\n", x, arg));
    }

    result.push_str(
        "\n
Here is my source code: 
  https://github.com/Xe/olin/blob/master/cwa/cwagi/src/main.rs

If you would like to learn more about this project, please 
take a look at the following links:
 - https://github.com/Xe/olin
 - https://christine.website/blog/olin-1-why-09-1-2018
 - https://christine.website/blog/olin-2-the-future-09-5-2018

If you know of a Rust HTTP library that lets users hijack
the transport layer (and can offer code examples, please, I'm
still new to Rust) so that it could be adapted to Olin's
native HTTP support, here is its test: 
  https://github.com/Xe/olin/blob/master/cwa/tests/src/scheme/http.rs

If you would like to, please feel free to load test the route
/cadey. That route does a minimal number of system calls, and
should make for the best benchmarking results. There is a more
detailed metrics view at /metrics. Olin does a level of automatic
load balancing that takes advantage of the Go runtime, so in
theory the limits of this program approach the limits of how fast
the code running in the WebAssembly environment is.

And this is just the beginning.

Have a good day and be well, creator.",
    );

    Response {
        status: 200,
        body: result,
    }
}

fn respond_to(ctx: &Context) -> Response {
    match ctx.request_uri.as_str() {
        "/cadey" => Response {
            status: 200,
            body: String::from("you are awesome!"),
        },
        _ => runtime_info(ctx),
    }
}

fn serialize(response: &Response) -> Vec<u8> {
    let mut output = String::new();
    output.push_str(&format!(
        "Status: {}\nContent-Type: text/plain\nCetacean-Powered-By: Cadey~#1337\nContent-Length: {}\n\n",
        response.status,
        response.body.len(),
    ));

    output.push_str(&response.body);
    output.into_bytes()
}
