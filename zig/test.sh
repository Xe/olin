#!/bin/sh

set -e
set -x

cwa -main-func _start -vm-stats allyourbase.wasm
cwa -vm-stats -test coi.wasm a b c d
cwa -vm-stats triangle.wasm
cwa -vm-stats httptest.wasm
uname -av | cwa -vm-stats cat.wasm

cwa exit0.wasm

set +e
cwa exit1.wasm
status=$?
if [ $status -ne 1 ]
then
    echo "expected exit status 1, got: $status"
    exit 1
fi
