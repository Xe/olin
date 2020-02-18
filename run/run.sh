#!/bin/sh

cwa-cgi -addr ":$PORT" -main-func _start -bin /wasm/cwagi_zig.wasm inside:docker
