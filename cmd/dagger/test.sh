#!/bin/sh

set -e
set -x

go run main.go -vm-stats -jail file://$(pwd)/testdata -main-func dagger_main ./testdata/test.wasm
