#!/bin/sh

set -e
set -x

zig build-lib -target wasm32-other-none --release-fast src/allyourbase.zig
zig build-lib -target wasm32-other-none --release-fast src/cat.zig
zig build-lib -target wasm32-other-none --release-fast src/httptest.zig
zig build-lib -target wasm32-other-none --release-fast src/shaman.zig
zig build-lib -target wasm32-other-none --release-fast src/cwagi.zig
zig build-lib -target wasm32-other-none --release-fast src/triangle.zig
zig build-lib -target wasm32-other-none --release-fast src/coi.zig
zig build-lib -target wasm32-other-none --release-fast src/exit0.zig
zig build-lib -target wasm32-other-none --release-fast src/exit1.zig

rm *.h *.o

