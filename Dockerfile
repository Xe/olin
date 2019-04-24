FROM rustlang/rust:nightly AS rust
WORKDIR /usr/src/olin
COPY ./cwa/ .
RUN rustup target add wasm32-unknown-unknown --toolchain nightly \
 && ./build.sh

FROM xena/rust-wasm-tools AS rust-wasm-tools
RUN mkdir -p /olin
WORKDIR /olin
COPY --from=rust /usr/src/olin/bin/*.wasm /olin/
RUN wasm-gc ./olinfetch.wasm \
 && wasm-gc ./cwagi.wasm \
 && wasm-gc ./cwa-tests.wasm \
 && du -hs ./*.wasm

FROM xena/go:1.12.1 AS go
RUN apk add --no-cache build-base
ENV GOPROXY https://cache.greedo.xeserv.us
WORKDIR /olin
COPY . .
COPY --from=rust-wasm-tools /olin/cwagi.wasm ./cmd/cwa-cgi/testdata/test.wasm
COPY --from=rust-wasm-tools /olin/cwa-tests.wasm ./cmd/cwa/testdata/test.wasm
RUN GOARCH=wasm GOOS=js go build -o ./cmd/cwa/testdata/go.wasm ./internal/abi/wasmgo/testdata/nothing.go
RUN go test -v ./cmd/... ./internal/...
RUN GOBIN=/usr/local/bin go install ./cmd/cwa-cgi
RUN GOBIN=/usr/local/bin go install ./cmd/cwa

FROM xena/zig:0.4.0-0f8fc3b9 AS zig
WORKDIR /olin
COPY ./zig .
COPY --from=go /usr/local/bin/cwa /usr/local/bin/cwa
RUN ./build.sh

FROM xena/alpine
COPY ./run/run.sh /run.sh
COPY --from=rust-wasm-tools /olin/*.wasm /wasm/
COPY --from=zig /olin/coi.wasm /wasm/coi-zig.wasm
COPY --from=zig /olin/shaman.wasm /wasm/shaman.wasm
COPY --from=zig /olin/cwagi.wasm /wasm/cwagi_zig.wasm
COPY --from=rust /usr/src/olin/bin/shaman.wasm /wasm/shaman_rust.wasm
COPY --from=go /usr/local/bin/cwa /usr/local/bin/cwa
COPY --from=go /usr/local/bin/cwa-cgi /usr/local/bin/cwa-cgi
WORKDIR /
CMD ["/run.sh"]
