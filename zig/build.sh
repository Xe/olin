#!/bin/sh

set -e
set -x

ZIGFLAGS="-target wasm32-other-none --release-fast"

if [ -n "$HOME" -a "$HOME" = "/homeless-shelter" ]
then
    export HOME=$TMPDIR;
fi

zig build-lib $ZIGFLAGS src/allyourargs.zig
zig build-lib $ZIGFLAGS src/allyourbase.zig
zig build-lib $ZIGFLAGS src/allyourlogs.zig
zig build-lib $ZIGFLAGS src/runtime_name.zig
zig build-lib $ZIGFLAGS src/cat.zig
zig build-lib $ZIGFLAGS src/httptest.zig
zig build-lib $ZIGFLAGS src/shaman.zig
zig build-lib $ZIGFLAGS src/cwagi.zig
zig build-lib $ZIGFLAGS src/triangle.zig
zig build-lib $ZIGFLAGS src/coi.zig
zig build-lib $ZIGFLAGS src/exit0.zig
zig build-lib $ZIGFLAGS src/exit1.zig

rm *.h *.o

