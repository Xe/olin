#!/bin/sh

set -e
set -x

rustc --target wasm32-unknown-unknown -o test.wasm test.rs
wasm-gc test.wasm
