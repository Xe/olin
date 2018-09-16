FROM rustlang/rust:nightly AS rust
WORKDIR /usr/src/olin
COPY ./cwa/ .
RUN rustup target add wasm32-unknown-unknown --toolchain nightly \
 && ./build.sh

FROM xena/rust-wasm-tools AS rust-wasm-tools
RUN mkdir -p /olin
WORKDIR /olin
COPY --from=rust /usr/src/olin/cwagi/cwagi.wasm /olin/cwagi.wasm
COPY --from=rust /usr/src/olin/tests/cwa-tests.wasm /olin/tests.wasm
RUN wasm-gc ./cwagi.wasm \
 && wasm-gc ./tests.wasm \
 && du -hs ./*.wasm

FROM xena/go:1.10 AS go
RUN apk add --no-cache build-base
COPY . /root/go/src/github.com/Xe/olin
WORKDIR /root/go/src/github.com/Xe/olin
COPY --from=rust-wasm-tools /olin/cwagi.wasm ./cmd/cwa-cgi/testdata/test.wasm
COPY --from=rust-wasm-tools /olin/tests.wasm ./cmd/cwa/testdata/test.wasm
RUN go test ./cmd/... ./internal/...
RUN GOBIN=/usr/local/bin go install -tags heroku ./cmd/cwa-cgi

FROM xena/alpine
COPY ./run/run.sh /run.sh
COPY --from=rust-wasm-tools /olin/cwagi.wasm /main.wasm
COPY --from=go /usr/local/bin/cwa-cgi /cwa-cgi
WORKDIR /
CMD ["/run.sh"]
