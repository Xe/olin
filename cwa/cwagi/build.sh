#!/bin/sh

#cargo +nightly install wasm-gc
#cargo +nightly install wasm-snip
cargo +nightly -Z unstable-options build --target wasm32-unknown-unknown --out-dir ../bin --release
#wasm-gc cwagi.wasm
#wasm-snip cwagi.wasm -o cwagi.wasm main
