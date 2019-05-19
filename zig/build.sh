#!/bin/sh

set -e
set -x

zig build-lib -target wasm32-freestanding-none src/coi.zig
zig build-lib -target wasm32-freestanding-none --release-fast src/shaman.zig
zig build-lib -target wasm32-freestanding-none --release-fast src/cwagi.zig
zig build-lib -target wasm32-freestanding-none --release-fast src/triangle.zig
cwa -vm-stats -test coi.wasm a b c d
zig build-lib -target wasm32-freestanding-none --release-fast src/coi.zig
cwa -vm-stats -test coi.wasm a b c d
cwa -vm-stats triangle.wasm

zig build-lib -target wasm32-freestanding-none --release-fast src/dagger_test.zig

