FROM rustlang/rust:nightly AS rust
WORKDIR /usr/src/olin
COPY ./cwa/ .
RUN rustup target add wasm32-unknown-unknown --toolchain nightly \
 && ./build.sh

FROM xena/go:1.10 AS build
RUN apk add --no-cache build-base
COPY . /root/go/src/github.com/Xe/olin
WORKDIR /root/go/src/github.com/Xe/olin
COPY --from=rust /usr/src/olin/cwagi/cwagi.wasm /root/go/src/github.com/Xe/olin/cmd/cwa-cgi/testdata/test.wasm
COPY --from=rust /usr/src/olin/tests/cwa-tests.wasm /root/go/src/github.com/Xe/olin/cmd/cwa/testdata/test.wasm
RUN go test ./cmd/... ./internal/...
RUN GOBIN=/ go install -tags heroku ./cmd/cwa-cgi

FROM xena/alpine
COPY ./run/run.sh /run.sh
COPY --from=build /root/go/src/github.com/Xe/olin/cmd/cwa-cgi/testdata/test.wasm /main.wasm
COPY --from=build /cwa-cgi /cwa-cgi
WORKDIR /
CMD ["/run.sh"]
