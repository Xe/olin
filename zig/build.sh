#!/bin/sh

set -e
set -x

zig build-exe -target wasm32-freestanding-none --name coi.wasm src/coi.zig
cwa -vm-stats coi.wasm
