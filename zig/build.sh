#!/bin/sh

set -e
set -x

zig build-exe -target wasm32-freestanding-none --release-small src/main.zig
cwa -vm-stats main
