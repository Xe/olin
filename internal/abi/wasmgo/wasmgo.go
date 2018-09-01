package wasmgo

import (
	"log"
	"time"

	"github.com/Xe/olin/internal/abi"
	"github.com/perlin-network/life/exec"
)

type wasmGo struct {
	child abi.ABI

	BootTime   time.Time
	Exited     bool
	StatusCode int32
	Callbacks  map[int32]time.Time

	vm *exec.VirtualMachine
}

func (w *wasmGo) writeMem(ptr int32, data []byte) (int, error) {
	ctr := 0
	for i, d := range data {
		w.vm.Memory[ptr+int32(i)] = d
		log.Printf("memory[%x] = %x", ptr+int32(i), w.vm.Memory[ptr+int32(i)])
		ctr++
	}

	log.Printf("wasmGo.writeMem(%x, byte slice of len %d): %d", ptr, len(data), ctr)
	return ctr, nil
}
