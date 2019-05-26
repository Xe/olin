package fileresolver

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/Xe/olin/abi"
)

// HTTP creates a new HTTP transport that pretends to be a file. A process can
// (and should) have as many of these open as they need active HTTP connections.
// Users of this package are suggested to use the same underlying HTTP client as
// much as makes sense.
//
// To use this file:
//
//    - write your request body to the file
//    - flush the file (this blocks for the duration of the HTTP request)
//    - read the response body from the file
//    - close the file when you are done reading the response
//
// A new HTTP file will be required for every request on the guest side, but it
// will be for the best.
func HTTP(cl *http.Client, u *url.URL) (abi.File, error) {
	return &httpFile{
		cl:   cl,
		u:    u,
		req:  bytes.NewBuffer(nil),
		resp: bytes.NewBuffer(nil),
	}, nil

}

type httpFile struct {
	cl        *http.Client
	u         *url.URL
	req, resp *bytes.Buffer
	respGot   bool
}

func (h httpFile) Name() string { return h.u.String() }

func (h *httpFile) Write(p []byte) (int, error) {
	n, err := h.req.Write(p)
	if err != nil {
		return -1, err
	}

	return n, nil
}

func (h *httpFile) Read(p []byte) (int, error) {
	if !h.respGot {
		return -1, errors.New("no response data yet")
	}
	return h.resp.Read(p)
}

func (h *httpFile) Flush() error {
	buf := bufio.NewReader(h.req)
	mreq, err := http.ReadRequest(buf)
	if err != nil {
		return err
	}

	mreq.URL.Scheme = "http"
	mreq.URL.Host = mreq.Host
	mreq.URL.Path = mreq.RequestURI

	req, err := http.NewRequest(mreq.Method, mreq.URL.String(), mreq.Body)
	if err != nil {
		return err
	}

	req.Header = mreq.Header

	resp, err := h.cl.Do(req)
	if err != nil {
		return err
	}

	resp.ProtoMajor = 1
	resp.ProtoMinor = 1
	resp.Proto = "HTTP/1.1"

	err = resp.Write(h.resp)
	if err != nil {
		return err
	}
	h.respGot = true

	log.Printf("%s: %d bytes waiting in resp", h.u, h.resp.Len())

	return nil
}

func (h *httpFile) Close() error {
	h.resp.Reset()
	h.req.Reset()

	return nil
}
