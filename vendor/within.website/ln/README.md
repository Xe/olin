# ln: The Natural Logger for Go

![Release](https://img.shields.io/github/release/Xe/ln.svg) [![License](https://img.shields.io/github/license/Xe/ln.svg)](https://github.com/Xe/ln/blob/master/LICENSE) [![GoDoc](https://godoc.org/within.website/ln?status.svg)](https://godoc.org/within.website/ln) [![Build Status](https://travis-ci.org/Xe/ln.svg?branch=master)](https://travis-ci.org/Xe/ln) [![Go Report Card](https://goreportcard.com/badge/github.com/Xe/ln)](https://goreportcard.com/report/github.com/Xe/ln)

`ln` provides a simple interface to contextually aware structured logging. 
The design of `ln` centers around the idea of key-value pairs, which
can be interpreted on the fly, but "Filters" to do things such as
aggregated metrics, and report said metrics to, say Librato, or
statsd.

"Filters" are like WSGI, or Rack Middleware. They are run "top down"
and can abort an emitted log's output at any time, or continue to let
it through the chain. However, the interface is slightly different
than that. Rather than encapsulating the chain with partial function
application, we utilize a simpler method, namely, each plugin defines
an `Apply` function, which takes as an argument the log event, and
performs the work of the plugin, only if the Plugin "Applies" to this
log event.

If `Apply` returns `false`, the iteration through the rest of the
filters is aborted, and the log is dropped from further processing.

## Examples

```go
type User struct {
	ID       int
	Username string
}

func (u User) F() ln.F {
	return ln.F{
		"user_id":       u.ID,
		"user_username": u.Username,
	}
}

// in some function
u, err := CreateUser(ctx, "Foobar")
if err != nil {
  ln.Error(ctx, err) // errors can also be Fers
  return
}

ln.Log(ctx, ln.Info("created new user"), u) // -> user_id=123 user_username=Foobar msg="created new user"
```

### Current Status: Known Stable

## Copyright

(c) 2015-2019, Andrew Gwozdziewycz, Christine Dodrill, BSD Licensed. 
See LICENSE for more info.

This library was forked from [apg's ln](https://github.com/apg/ln).
This library is not the official ln (though it is mostly compatible 
with it), but it is at least what I think a logger should be.

