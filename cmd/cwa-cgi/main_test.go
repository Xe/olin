package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Xe/olin/internal/abi/cwa"
	"github.com/perlin-network/life/exec"
)

func TestHTTPRequest(t *testing.T) {
	const fname = "./testdata/test.wasm"

	ts := httptest.NewServer(loadWasm(t, fname, "cwa_main"))
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/cadey", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("wanted status %d, got: %d", http.StatusOK, resp.StatusCode)
	}
}
func loadWasm(t *testing.T, fname, main string) *vmServer {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatal(err)
	}

	p := cwa.NewProcess(fname, []string{"test.wasm"}, map[string]string{})

	cfg := exec.VMConfig{
		EnableJIT:          false,
		DefaultMemoryPages: 32,
	}
	vm, err := exec.NewVirtualMachine(data, cfg, p)
	if err != nil {
		t.Fatal(err)
	}

	mn, ok := vm.GetFunctionExport(main)
	if !ok {
		t.Fatal(err)
	}

	vs := &vmServer{
		vm:       vm,
		p:        p,
		mainFunc: mn,
	}

	return vs
}
