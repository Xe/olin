#!/bin/sh

cwa-cgi -addr ":$PORT" -pool-size 1 -max-pool-size 16 -main-func cwa_main /wasm/cwagi.wasm inside:docker
