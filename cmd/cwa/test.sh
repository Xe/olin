#!/bin/sh

go run main.go -vm-stats -test -main-func cwa_main ./testdata/test.wasm foo bar
