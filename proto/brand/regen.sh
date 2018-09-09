#!/bin/sh

protoc --proto_path=$GOPATH/src:. \
       --go_out=. \
       brand.proto
