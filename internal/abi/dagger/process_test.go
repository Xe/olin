package dagger

import (
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"syscall"
	"testing"

	"github.com/Xe/olin/internal/abi"
	"github.com/perlin-network/life/exec"
)

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

type testProcessFragment func(t *testing.T, p *Process)

func TestProcess(t *testing.T) {
	cases := []testProcessFragment{
		processInsertFile,
		processOpenFile,
		processCantOpenUnknownScheme,
		processTestOpenWasm,
		processTestHelloWorld,
		processTestHTTPFile,
	}

	for _, cs := range cases {
		t.Run(getFunctionName(cs), func(t *testing.T) {
			cs(t, &Process{})
		})
	}
}

func wantError(t *testing.T, err error) Error {
	t.Helper()

	merr, ok := err.(Error)
	if !ok {
		t.Fatalf("err is %T, not Error", err)
	}

	return merr
}

func ensureErrorCode(t *testing.T, err error, code Errno) {
	t.Helper()

	merr := wantError(t, err)
	if merr.Errno != code {
		t.Fatalf("wanted error code to be %d (%s), got: %d (%s)", code, code, merr.Errno, merr.Errno)
	}
}

func processInsertFile(t *testing.T, p *Process) {
	o := abi.NewOSFile(uintptr(syscall.Stdout), "stdout")
	i := p.insertFile(o)
	i2 := p.insertFile(o)

	if i == i2 {
		t.Fatalf("expected to get different results when adding two files to the process, got: %d %d", i, i2)
	}
}

func processOpenFile(t *testing.T, p *Process) {
	_, err := p.open("fd://1", 0)
	if err != nil {
		t.Fatal(err)
	}
}

func processCantOpenUnknownScheme(t *testing.T, p *Process) {
	_, err := p.open("unknown+scheme://some.host/files/i/guess", 0)
	if err == nil {
		t.Fatal("expected error")
	}

	ensureErrorCode(t, err, ErrorUnknownScheme)
}

func openAndRunWasmMain(t *testing.T, p *Process, fname string, isok func(int64) bool) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("%s: %v", fname, err)
	}

	cfg := exec.VMConfig{}
	vm, err := exec.NewVirtualMachine(data, cfg, p)
	if err != nil {
		t.Fatalf("%s: %v", fname, err)
	}

	main, ok := vm.GetFunctionExport("main")
	if !ok {
		t.Fatalf("%s: no main function exported", fname)
	}

	ret, err := vm.Run(main)
	if err != nil {
		t.Fatalf("%s: vm error: %v", fname, err)
	}

	if !isok(ret) {
		t.Fatalf("%s returned %d which is not ok", fname, ret)
	}
}

func processTestOpenWasm(t *testing.T, p *Process) {
	openAndRunWasmMain(t, p, "./testdata/open.wasm", func(i int64) bool { return i == 0 })
}

func processTestCloseWasm(t *testing.T, p *Process) {
	openAndRunWasmMain(t, p, "./testdata/close.wasm", func(i int64) bool { return i == 0 })
}

func processTestHelloWorld(t *testing.T, p *Process) {
	openAndRunWasmMain(t, p, "./testdata/helloworld.wasm", func(i int64) bool { return i == 0 })
}

func processTestHTTPFile(t *testing.T, p *Process) {
	hs := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("got request from wasm :D")
			http.Error(w, "success", http.StatusOK)
		}),
		Addr: "127.0.0.1:30405",
	}
	go func() {
		err := hs.ListenAndServe()
		if err != nil {
			t.Skip("travis is broken for this test")
			t.Fatal(err)
		}
	}()
	defer hs.Close()

	openAndRunWasmMain(t, p, "./testdata/http.wasm", func(i int64) bool { return i == 0 })
}
