package wasmgo

import (
	"io/ioutil"
	"syscall"
	"testing"
	"time"

	"github.com/Xe/olin/internal/abi"

	"github.com/Xe/olin/internal/abi/dagger"
	"github.com/perlin-network/life/exec"
)

func (w *wasmGo) ensureExitStatus(t *testing.T, exited bool, status int32) {
	t.Helper()

	if w.Exited != exited {
		t.Errorf("exited: wanted: %v, got: %v", exited, w.Exited)
	}

	if w.StatusCode != status {
		t.Errorf("status code: wanted: %d, got: %d", status, w.StatusCode)
	}
}

type testWasmGoFragment func(t *testing.T, w *wasmGo)

func TestWasmGo(t *testing.T) {
	cases := []testWasmGoFragment{
		testNothing,
		testDagger,
	}

	for _, cs := range cases {
		t.Run(getFunctionName(cs), func(t *testing.T) {
			p := dagger.NewProcess(getFunctionName(cs))
			p.Open(abi.NewOSFile(uintptr(syscall.Stdin), "stdin"))
			p.Open(abi.NewOSFile(uintptr(syscall.Stdout), "stdout"))
			p.Open(abi.NewOSFile(uintptr(syscall.Stderr), "stderr"))

			w := &wasmGo{
				child:     p,
				BootTime:  time.Now(),
				Callbacks: map[int32]time.Time{},
			}
			cs(t, w)
		})
	}
}

func testNothing(t *testing.T, w *wasmGo) {
	openAndRunWasmRun(t, w, "./testdata/nothing.wasm")
	defer w.ensureExitStatus(t, true, 0)
}

func testDagger(t *testing.T, w *wasmGo) {
	openAndRunWasmRun(t, w, "./testdata/dagger.wasm")
	defer w.ensureExitStatus(t, true, 0)
}

func openAndRunWasmRun(t *testing.T, w *wasmGo, fname string) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("%s: %v", fname, err)
	}

	cfg := exec.VMConfig{}
	vm, err := exec.NewVirtualMachine(data, cfg, w)
	if err != nil {
		t.Fatalf("%s: %v", fname, err)
	}

	main, ok := vm.GetFunctionExport("run")
	if !ok {
		t.Fatalf("%s: no main function exported", fname)
	}

	_, err = vm.Run(main, 0, 0)
	if err != nil {
		t.Fatalf("%s: vm error: %v", fname, err)
	}
}
