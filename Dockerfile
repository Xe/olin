FROM rustlang/rust:nightly AS rust
WORKDIR /usr/src/olin
COPY ./rust .
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

FROM xena/zig:0.5.0-21a85d4fb as zig
WORKDIR /olin
COPY ./zig .
RUN ./build.sh

FROM xena/go:1.13.5 AS go
RUN apk add --no-cache build-base
ENV GOPROXY https://cache.greedo.xeserv.us
WORKDIR /olin
COPY . .
COPY --from=rust-wasm-tools /olin/cwa-tests.wasm ./cmd/cwa/testdata/test.wasm
RUN GOARCH=wasm GOOS=js go build -o ./cmd/cwa/testdata/go.wasm ./abi/wasmgo/testdata/nothing.go \
 && go test -v ./cmd/... \
 && GOBIN=/usr/local/bin go install ./cmd/cwa-cgi \
 && GOBIN=/usr/local/bin go install ./cmd/cwa
COPY --from=zig /olin/*.wasm ./zig/
RUN cd zig && ./test.sh

FROM xena/alpine
COPY ./run/run.sh /run.sh
COPY --from=rust-wasm-tools /olin/*.wasm /wasm/
COPY --from=zig /olin/coi.wasm /wasm/coi-zig.wasm
COPY --from=zig /olin/shaman.wasm /wasm/shaman.wasm
COPY --from=zig /olin/cwagi.wasm /wasm/cwagi_zig.wasm
COPY --from=rust /usr/src/olin/bin/shaman.wasm /wasm/shaman_rust.wasm
COPY --from=go /usr/local/bin/cwa /usr/local/bin/cwa
COPY --from=go /usr/local/bin/cwa-cgi /usr/local/bin/cwa-cgi
COPY --from=go /olin/run/olin-policy-mode.el /olin/olin-policy-mode.el
COPY --from=go /olin/docs /olin/docs
COPY --from=go /olin/LICENSE /olin/LICENSE
COPY --from=go /olin/README.md /olin/README.md
WORKDIR /
CMD ["/run.sh"]
