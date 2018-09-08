#!/bin/sh

cargo +nightly -Z unstable-options build --target wasm32-unknown-unknown --out-dir . --release
wasm-gc cwagi.wasm
cp cwagi.wasm ../../cmd/cwa-cgi/testdata/test.wasm
