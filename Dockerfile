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
 && wasm-gc ./shaman.wasm \
 && wasm-gc ./cwa-tests.wasm \
 && du -hs ./*.wasm

FROM xena/go:1.11.1 AS go
RUN apk add --no-cache build-base
COPY . /root/go/src/github.com/Xe/olin
WORKDIR /root/go/src/github.com/Xe/olin
COPY --from=rust-wasm-tools /olin/cwagi.wasm ./cmd/cwa-cgi/testdata/test.wasm
COPY --from=rust-wasm-tools /olin/cwa-tests.wasm ./cmd/cwa/testdata/test.wasm
RUN go test -v ./cmd/... ./internal/...
RUN GOBIN=/usr/local/bin go install -tags heroku ./cmd/cwa-cgi
RUN GOBIN=/usr/local/bin go install ./cmd/cwa

FROM xena/alpine
COPY ./run/run.sh /run.sh
COPY --from=rust-wasm-tools /olin/*.wasm /wasm/
COPY --from=go /usr/local/bin/cwa /usr/local/bin/cwa
COPY --from=go /usr/local/bin/cwa-cgi /usr/local/bin/cwa-cgi
WORKDIR /
CMD ["/run.sh"]
