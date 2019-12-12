#!/bin/sh

set -e

mkdir -p bin

(cd olin && ./build.sh)
(cd cwagi && ./build.sh) &
(cd tests && ./build.sh) &

wait
