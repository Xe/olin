#!/bin/sh

/cwa-cgi -addr ":$PORT" -pool-size 16 -max-pool-size 128 -main-func cwa_main /main.wasm inside:docker
