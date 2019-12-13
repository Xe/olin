#!/bin/sh

set -e
set -x

cwa -vm-stats -policy allyourbase.policy allyourbase.wasm
cwa -vm-stats -policy coi.policy -test coi.wasm a b c d
cwa -vm-stats -policy triangle.policy triangle.wasm
cwa -vm-stats -policy httptest.policy httptest.wasm
uname -av | cwa -policy cat.policy -vm-stats cat.wasm

cwa -policy exit.policy exit0.wasm

set +e
cwa -policy exit.policy exit1.wasm
status=$?
if [ $status -ne 1 ]
then
    echo "expected exit status 1, got: $status"
    exit 1
fi

cwa -policy ../policy/testdata/gitea.policy httptest.wasm
status=$?
if [ $status -ne 1 ]
then
    echo "expected exit status 1, got: $status"
    exit 1
fi
