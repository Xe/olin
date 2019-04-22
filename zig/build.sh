#!/bin/sh

set -e
set -x

zig build-exe -target wasm32-freestanding-none src/coi.zig
zig build-exe -target wasm32-freestanding-none --release-fast src/shaman.zig
zig build-exe -target wasm32-freestanding-none --release-fast src/triangle.zig
cwa -vm-stats -test -write-mem coi.mem.wasm coi.wasm
cwa -vm-stats triangle.wasm
