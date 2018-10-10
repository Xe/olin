#!/bin/sh

cargo +nightly -Z unstable-options build --target wasm32-unknown-unknown --out-dir ../bin --release
