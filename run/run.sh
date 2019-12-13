#!/bin/sh

cwa-cgi -addr ":$PORT" -main-func cwa_main -bin /wasm/cwagi_zig.wasm inside:docker
