package wasmgo

import (
	"time"

	"github.com/Xe/olin/internal/abi/dagger"
	"github.com/perlin-network/life/exec"
)

// TODO(Xe): upgrade this to a CWA process
type wasmGo struct {
	child *dagger.Process

	BootTime   time.Time
	Exited     bool
	StatusCode int32
	Callbacks  map[int32]time.Time

	vm *exec.VirtualMachine

	values []interface{}
	refs   map[interface{}]int
}

// TODO(Xe): replace with copy() or obviate the need to copy to begin with?
// while this code is being run, the WebAssembly memory is considered locked
// because this is a single-threaded environment at the moment.
func (w *wasmGo) writeMem(ptr int32, data []byte) (int, error) {
	ctr := 0
	for i, d := range data {
		w.vm.Memory[ptr+int32(i)] = d
		ctr++
	}

	return ctr, nil
}
