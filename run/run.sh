#!/bin/sh

/cwa-cgi -addr ":$PORT" -main-func cwa_main /main.wasm inside:docker
