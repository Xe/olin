#!/bin/sh

set -e
set -x

zig build-lib -target wasm32-freestanding-none src/coi.zig
zig build-lib -target wasm32-freestanding-none --release-fast src/shaman.zig
zig build-lib -target wasm32-freestanding-none --release-fast src/cwagi.zig
zig build-lib -target wasm32-freestanding-none --release-fast src/triangle.zig
zig build-lib -target wasm32-freestanding-none --release-fast src/cwa_hello.zig
zig build-lib -target wasm32-freestanding-none --release-fast src/dagger_test.zig
zig build-lib -target wasm32-freestanding-none --release-fast src/dagger_hello.zig
rm *.h
