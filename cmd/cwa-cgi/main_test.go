package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/Xe/olin/abi/cwa"
	"github.com/pborman/uuid"
	"github.com/perlin-network/life/compiler"
	"github.com/perlin-network/life/exec"
	"within.website/ln"
	"within.website/ln/opname"
)

type vmServer struct {
	vm       *exec.VirtualMachine
	p        *cwa.Process
	lock     sync.Mutex
	mainFunc int
}

func (v *vmServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stdout := bytes.NewBuffer(nil)
	stdin := r.Body
	defer r.Body.Close()
	ctx := opname.With(r.Context(), "vmServer.ServeHTTP")
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	runID := uuid.New()

	f := ln.F{
		"main_func":    v.mainFunc,
		"process_name": v.p.Name(),
		"run_id":       runID,
		"method":       r.Method,
		"request_uri":  r.RequestURI,
	}

	v.lock.Lock()
	defer v.lock.Unlock()
	v.p.Stdin = stdin
	v.p.Stdout = stdout
	v.p.Setenv(map[string]string{
		"REQUEST_METHOD": r.Method,
		"REQUEST_URI":    r.RequestURI,
		"QUERY_STRING":   r.URL.Query().Encode(),
		"RUN_ID":         runID,
		"WORKER_ID":      uuid.New(),
	})

	t0 := time.Now()
	ret, err := v.vm.Run(v.mainFunc)
	if err != nil {
		http.Error(w, "internal server error: VM error, run ID: "+runID, http.StatusInternalServerError)
		go func() {
			time.Sleep(125 * time.Millisecond)
			ln.FatalErr(ctx, err, f)
		}()
		return
	}
	f["exec_dur"] = time.Since(t0)

	if ret != 0 {
		ln.Log(ctx, f, ln.F{
			"return_value": ret,
		})
		http.Error(w, fmt.Sprintf("internal server error: return code %d", ret), http.StatusInternalServerError)
		return
	}

	ctx = opname.With(ctx, "respond")
	resp, err := http.ReadResponse(bufio.NewReader(stdout), r)
	if err != nil {
		ln.Error(ctx, err, f)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// copy headers
	for k, v := range resp.Header {
		for _, val := range v {
			w.Header().Add(k, val)
		}
	}

	// copy status code
	w.WriteHeader(resp.StatusCode)
	f["status"] = resp.StatusCode

	// copy body
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		ln.Error(opname.With(ctx, "copy_body"), err, f)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	ln.Log(ctx, f, ln.Info("successful invocation"))
}

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

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if bytes.Contains(data, []byte{0}) {
		t.Fatalf("response body was garbage: %x", data)
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
	gp := compiler.SimpleGasPolicy{}
	vm, err := exec.NewVirtualMachine(data, cfg, p, &gp)
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
