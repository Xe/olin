// +build js,wasm go1.11

package dagger

import (
	"errors"
	"strconv"
)

func OpenFile(furl string, flags int32) int32 {
	return openFD(furl, flags)
}

func openFD(furl string, flags int32) int32

type file struct {
	fd int
}

func read(fd int, buf []byte) int

func (f file) Read(buf []byte) (int, error) {
	n := read(f.fd, buf)
	if n < 0 {
		return n, errors.New("dagger: error code " + strconv.Itoa(n*-1))
	}

	return n, nil
}
