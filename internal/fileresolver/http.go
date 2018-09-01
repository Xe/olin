package fileresolver

import (
	"bufio"
	"bytes"
	"net/http"
	"net/url"

	"github.com/Xe/olin/internal/abi"
)

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
	return h.resp.Read(p)
}

func (h *httpFile) Sync() error {
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

	err = resp.Write(h.resp)
	if err != nil {
		return err
	}

	return nil
}

func (h *httpFile) Close() error {
	h.resp.Reset()
	h.req.Reset()

	return nil
}
