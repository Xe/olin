package wasmgo

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"log"
	"time"

	"github.com/perlin-network/life/exec"
)

func (w *wasmGo) getInt32(addr int32) int32 {
	mem := w.vm.Memory[addr : addr+4]

	var result int32
	err := binary.Read(bytes.NewReader(mem), binary.LittleEndian, &result)
	if err != nil {
		panic(err)
	}

	return result
}

func (w *wasmGo) setInt32(addr int32, val int32) {
	data := make([]byte, 0, 4)
	buf := bytes.NewBuffer(data)
	err := binary.Write(buf, binary.LittleEndian, val)
	if err != nil {
		panic(err)
	}

	_, err = w.writeMem(addr, buf.Bytes())
	if err != nil {
		panic(err)
	}
}

func (w *wasmGo) getInt64(addr int32) int64 {
	mem := w.vm.Memory[addr : addr+8]

	var result int64
	err := binary.Read(bytes.NewReader(mem), binary.LittleEndian, &result)
	if err != nil {
		panic(err)
	}

	return result
}

func (w *wasmGo) setInt64(addr int32, val int64) {
	data := make([]byte, 0, 8)
	buf := bytes.NewBuffer(data)
	err := binary.Write(buf, binary.LittleEndian, val)
	if err != nil {
		panic(err)
	}

	_, err = w.writeMem(addr, buf.Bytes())
	if err != nil {
		panic(err)
	}
}

// goRuntimeWasmExit implements the go runtime function runtime.wasmExit. It uses
// the Go abi.
//
// This has the effective type of:
//
//     func (w *wasmGo) goRuntimeWasmExit(code int32)
func (w *wasmGo) goRuntimeWasmExit(sp int32) {
	exitCode := w.getInt32(sp + 8)
	w.Exited = true
	w.StatusCode = exitCode
}

// goRuntimeWasmWrite wraps the write function.
//
// This has the effective type of:
//
//     func (w *wasmGo) goRuntimeWasmWrite(fd int64, ptr int64, n int32)
func (w *wasmGo) goRuntimeWasmWrite(sp int32) {
	fd := w.getInt64(sp + 8)
	ptr := w.getInt64(sp + 16)
	n := w.getInt32(sp + 24)

	cnt, err := w.child.Files()[fd].Write(w.vm.Memory[ptr : ptr+int64(n)])
	if err != nil {
		log.Printf("goRuntimeWasmWrite(%d (%s), %x, %d): %v", fd, w.child.Files()[fd].Name(), ptr, n, err)
	}

	if int32(cnt) != n {
		log.Printf("goRuntimeWasmWrite(%d (%s), %x, %d): cnt: %d, err: nil", fd, w.child.Files()[fd].Name(), ptr, n, cnt)
	}
}

// goRuntimeNanotime implements the go runtime function runtime.nanotime. It uses
// the Go abi.
//
// This has the effective type of:
//
//     func (w *wasmGo) goRuntimeNanotime() int64
func (w *wasmGo) goRuntimeNanotime(sp int32) {
	now := time.Now().UnixNano()
	w.setInt64(sp+8, int64(now))
}

// goRuntimeWalltime implements the go runtime function runtime.walltime. It uses
// the Go abi.
//
// This has the effective type of:
//
//     func (w *wasmGo) goRuntimeWalltime() (int64, int32)
func (w *wasmGo) goRuntimeWalltime(sp int32) {
	now := time.Now()
	w.setInt64(sp+8, now.Unix())
	w.setInt32(sp+16, int32(now.Nanosecond()))
}

// goRuntimeScheduleCallback implements the go runtime function runtime.scheduleCallback.
// It uses the Go abi and is not implemented yet.
func (w *wasmGo) goRuntimeScheduleCallback(sp int32) {
	w.setInt32(sp+16, -1)
}

func (w *wasmGo) goRuntimeClearScheduledCallback(sp int32) {}

// goLoadSlice loads a Go slice out of wasm memory. It uses the the Go abi.
func (w *wasmGo) goLoadSlice(sp int32) []byte {
	arr := w.getInt64(sp)
	len := w.getInt64(sp + 8)
	return w.vm.Memory[arr : arr+len]
}

func (w *wasmGo) goRuntimeGetRandomData(sp int32) {
	data := w.goLoadSlice(sp + 8)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}

	w.writeMem(sp+8, data)
}

func (w *wasmGo) notImplemented(module, field string) goABIFunc {
	return func(sp int32) {
		w.vm.PrintStackTrace()
		log.Panicf("not implemented: %s %s", module, field)
	}
}

type goABIFunc func(int32)

var (
	goABIFuncs = []string{
		"debug",
		"runtime.wasmExit",
		"runtime.wasmWrite",
		"runtime.nanotime",
		"runtime.walltime",
		"runtime.scheduleCallback",
		"runtime.clearScheduledCallback",
		"runtime.getRandomData",
		"syscall/js.stringVal",
		"syscall/js.valueGet",
		"syscall/js.valueCall",
		"syscall/js.valueNew",
		"syscall/js.valuePrepareString",
		"syscall/js.valueLoadString",
	}
)

func (w *wasmGo) runGoABI(doer func(int32)) exec.FunctionImport {
	return func(vm *exec.VirtualMachine) int64 {
		if w.vm == nil {
			w.vm = vm
		}

		f := vm.GetCurrentFrame()
		sp := int32(f.Locals[0])
		doer(sp)
		return 0
	}
}

// ResolveGlobal does nothing, currently.
func (w *wasmGo) ResolveGlobal(module, field string) int64 { return 0 }

func (w *wasmGo) ResolveFunc(module, field string) exec.FunctionImport {
	val := w.child.ResolveFunc(module, field)
	if val != nil {
		return val
	}

	switch module {
	case "go":
		switch field {
		case "debug":
			return func(vm *exec.VirtualMachine) int64 {
				f := vm.GetCurrentFrame()
				log.Printf("debug: %x %d", f.Locals[0], f.Locals[0])
				return 0
			}

		case "runtime.wasmExit":
			return w.runGoABI(w.goRuntimeWasmExit)
		case "runtime.wasmWrite":
			return w.runGoABI(w.goRuntimeWasmWrite)
		case "runtime.nanotime":
			return w.runGoABI(w.goRuntimeNanotime)
		case "runtime.walltime":
			return w.runGoABI(w.goRuntimeWalltime)
		case "runtime.scheduleCallback":
			return w.runGoABI(w.goRuntimeScheduleCallback)
		case "runtime.clearScheduledCallback":
			return w.runGoABI(w.goRuntimeClearScheduledCallback)
		case "runtime.getRandomData":
			return w.runGoABI(w.goRuntimeGetRandomData)
		case "github.com/Xe/olin/dagger.openFD":
			return w.runGoABI(w.daggerOpenFD)
		default:
			log.Printf("unknown module+field %s %s, using shim", module, field)
			return w.runGoABI(w.notImplemented(module, field))
		}
	default:
		log.Panicf("unknown module+field %s %s", module, field)
	}

	panic("not implemented")
}
