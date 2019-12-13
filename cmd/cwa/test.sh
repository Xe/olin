#!/bin/sh

set -e
set -x

go run main.go -vm-stats -test ./testdata/test.wasm foo bar
go run main.go -vm-stats -test -go ./testdata/go.wasm foo bar ||:
