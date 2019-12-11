#!/bin/sh

set -e
set -x

zig build-lib -target wasm32-freestanding-none src/coi.zig
zig build-lib -target wasm32-freestanding-none --release-fast src/cat.zig
zig build-lib -target wasm32-freestanding-none --release-fast src/httptest.zig
zig build-lib -target wasm32-freestanding-none --release-fast src/shaman.zig
zig build-lib -target wasm32-freestanding-none --release-fast src/cwagi.zig
zig build-lib -target wasm32-freestanding-none --release-fast src/triangle.zig
cwa -vm-stats -test coi.wasm a b c d
zig build-lib -target wasm32-freestanding-none --release-fast src/coi.zig
cwa -vm-stats -test coi.wasm a b c d
zig build-lib -target wasm32-freestanding-none --release-fast src/exit0.zig
zig build-lib -target wasm32-freestanding-none --release-fast src/exit1.zig
cwa -vm-stats triangle.wasm

cwa exit0.wasm

set +e
cwa exit1.wasm
status=$?
if [ $status -ne 1 ]
then
  echo "expected exit status 1, got: $status"
  exit 1
fi

rm *.h *.o
