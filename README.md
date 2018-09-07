# olin

[![Build Status](https://travis-ci.org/Xe/olin.svg?branch=master)](https://travis-ci.org/Xe/olin) [![Go Report Card](https://goreportcard.com/badge/github.com/Xe/olin)](https://goreportcard.com/report/github.com/Xe/olin) [![GoDoc](https://godoc.org/github.com/Xe/olin?status.svg)](https://godoc.org/github.com/Xe/olin) ![powered by WebAssembly](https://img.shields.io/badge/powered%20by-WebAssembly-orange.svg)

Olin is an environment to run and operate functions as a service projects using
event sourcing and webassembly under the hood. Your handler code shouldn't need
to care that there is an event queue involved. Your handler should just do what
your handler needs to do.

## Background

Very frequently, I end up needing to write applications that basically end up
waiting forever to make sure things get put in the right place and then the 
right code runs as a response. I then have to make sure these things get put
in the right places and then that the right versions of things are running for
each of the relevant services. This doesn't scale very well, not to mention is 
hard to secure. This leads to a lot of duplicate infrastructure over time and 
as things grow. Not to mention adding in tracing, metrics and log aggreation.

I would like to change this.

I would like to make a perscriptive environment kinda like [Google Cloud Functions][gcf]
or [AWS Lambda][lambda] backed by a durable message queue and with handlers
compiled to webassembly to ensure forward compatibility. As such, the ABI 
involved will be versioned, documented and tested. Multiple ABI's will eventually
need to be maintained in parallel, so it might be good to get used to that early
on.

I expect this project to last decades. I want the binary modules I upload today
to be still working in 5 years, assuming its dependencies outside of the module 
still work. 

## Blogposts

Before asking questions or using Olin in this early state, I ask you please read
these blogposts outlining the point of this project and some other bikeshedding
I have put out on the topic:

- https://christine.website/blog/olin-1-why-09-1-2018
- https://christine.website/blog/olin-2-the-future-09-5-2018

These will explain a lot that hasn't been fit into here yet.

## ABI's Supported

### Common WebAssembly

Olin includes support for binaries linked against the [Common WebAssembly](https://github.com/CommonWA/cwa-spec)
specification. Please see the tests in `cmd/cwa` for more information. Currently
the Common WebAssembly is fairly basic, but at the same time there are currently
the most tests targeting Common Webassembly.
The tests for the Common WebAssembly spec can be found [here](https://github.com/Xe/olin/blob/master/cmd/cwa/testdata/test.rs).

### Dagger

> The dagger of light that renders your self-importance a decisive death

Dagger is currently in use for testing purposes. It defines five simple system 
calls (`open`, `read`, `write`, `sync` and `close`) and allows the user to chain
them as they wish. `open` returns the file descriptor that is going to be the 
first argument of all of the other functions.

To use these functions from C, import them like such from the `dagger` module:

```c
extern int open(const char *furl, int flags);
extern int close(int fd);
extern int read(int fd, void *buf, int nbyte);
extern int write(int fd, void *buf, int nbyte);
extern int sync(int fd);
```

### Go

Olin also includes support for running webassembly modules created by [Go 1.11's webassembly support](https://golang.org/wiki/WebAssembly).
It uses [the `wasmgo` ABI][wasmgo] package in order to do things. Right now
this is incredibly basic, but should be extendable to more things in the future.

As an example:

```go
// +build js,wasm ignore
// hello_world.go

package main

func main() {
	println("Hello, world!")
}
```

when compiled like this:

```console
$ GOARCH=wasm GOOS=js go1.11 build -o hello_world.wasm hello_world.go
```

produces the following output when run with the testing shim:

```
=== RUN   TestWasmGo/github.com/Xe/olin/internal/abi/wasmgo.testHelloWorld
Hello, world!
--- PASS: TestWasmGo (1.66s)
    --- PASS: TestWasmGo/github.com/Xe/olin/internal/abi/wasmgo.testHelloWorld (1.66s)
```

Currently Go binaries cannot interface with the Dagger ABI. There is [an issue](https://github.com/Xe/olin/issues/5)
open to track the solution to this.

Future posts will include more detail about using Go on top of Olin. 

Under the hood, the Olin implementation of the Go ABI currently uses Dagger.

## Project Meta

To follow the project, check it on GitHub [here][olin]. To talk about it on Slack,
join the [Go community Slack][goslack] and join `#olin`. 

[gcf]: https://cloud.google.com/functions/
[lambda]: https://aws.amazon.com/lambda/
[syscall]: https://en.wikipedia.org/wiki/System_call
[olin]: https://github.com/Xe/olin
[goslack]: https://invite.slack.golangbridge.org
[wasmgo]: https://github.com/Xe/olin/tree/master/internal/abi/wasmgo
