// +build !linux

package fileresolver

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHTTP(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "hello", http.StatusOK)
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	f, err := HTTP(&http.Client{}, u)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.URL = u

	err = req.Write(f)
	if err != nil {
		t.Fatal(err)
	}

	err = f.Flush()
	if err != nil {
		t.Fatal(err)
	}

	buf := bufio.NewReader(f)
	resp, err := http.ReadResponse(buf, req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "hello\n" {
		t.Fatalf("wanted \"hello\", got: %q", string(data))
	}
}
