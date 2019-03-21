package wasmgo

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"log"
	"math"
	"reflect"
	"runtime"
	"time"

	"github.com/perlin-network/life/exec"
)

func (w *wasmGo) loadString(addr int32) string {
	saddr := w.getInt64(addr)
	leng := w.getInt64(addr + 4)
	return string(w.vm.Memory[saddr : saddr+leng])
}

// TODO properly endian
func (w *wasmGo) getUint8(addr int32) uint8 {
	return uint8(w.vm.Memory[addr])
}

// TODO properly endian
func (w *wasmGo) setUint8(addr int32, val uint8) {
	w.vm.Memory[addr] = byte(val)
}

func (w *wasmGo) getFloat64(addr int32) float64 {
	mem := w.vm.Memory[addr : addr+8]
	var result float64
	err := binary.Read(bytes.NewReader(mem), binary.LittleEndian, &result)
	if err != nil {
		panic(err)
	}

	return result
}

func (w *wasmGo) setFloat64(addr int32, val float64) {
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

func (w *wasmGo) getUint32(addr int32) uint32 {
	mem := w.vm.Memory[addr : addr+4]
	var result uint32
	err := binary.Read(bytes.NewReader(mem), binary.LittleEndian, &result)
	if err != nil {
		panic(err)
	}

	return result
}

func (w *wasmGo) setUint32(addr int32, val uint32) {
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

func (w *wasmGo) loadValue(addr int32) interface{} {
	f := w.getFloat64(addr)
	if f == 0 {
		return Undefined
	}

	if !math.IsNaN(f) {
		return f
	}

	id := w.getUint32(addr)
	return w.values[id]
}

func (w *wasmGo) storeValue(addr int32, v interface{}) {
	const nanHead = 0x7FF80000
	if i, ok := v.(int); ok {
		v = float64(i)
	}
	if i, ok := v.(uint); ok {
		v = float64(i)
	}
	if v, ok := v.(float64); ok {
		if math.IsNaN(v) {
			w.setUint32(addr+4, nanHead)
			w.setUint32(addr, 0)
			return
		}
		if v == 0 {
			w.setUint32(addr+4, nanHead)
			w.setUint32(addr, 1)
			return
		}
		w.setFloat64(addr, v)
		return
	}

	switch v {
	case Undefined:
		w.setFloat64(addr, 0)
		return
	case nil:
		w.setUint32(addr+4, nanHead)
		w.setUint32(addr, 2)
		return
	case true:
		w.setUint32(addr+4, nanHead)
		w.setUint32(addr, 3)
		return
	case false:
		w.setUint32(addr+4, nanHead)
		w.setUint32(addr, 4)
		return
	}

	ref, ok := w.refs[v]
	if !ok {
		ref = len(w.values)
		w.values = append(w.values, v)
		w.refs[v] = ref
	}

	typeFlag := 0
	switch v.(type) {
	case string:
		typeFlag = 1
		// TODO symbol
		// TODO function
	}
	w.setUint32(addr+4, uint32(nanHead|typeFlag))
	w.setUint32(addr, uint32(ref))
}

func (w *wasmGo) loadSliceOfValues(addr int32) []interface{} {
	array := int(w.getInt64(addr))
	leng := int(w.getInt64(addr + 8))
	result := make([]interface{}, leng)

	for i := range result {
		result[i] = w.loadValue(int32(array + i*8))
	}

	return result
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

	cnt, err := w.Files()[fd].Write(w.vm.Memory[ptr : ptr+int64(n)])
	if err != nil {
		log.Printf("goRuntimeWasmWrite(%d (%s), %x, %d): %v", fd, w.Files()[fd].Name(), ptr, n, err)
	}

	if int32(cnt) != n {
		log.Printf("goRuntimeWasmWrite(%d (%s), %x, %d): cnt: %d, err: nil", fd, w.Files()[fd].Name(), ptr, n, cnt)
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

// goSyscallJSValueGet implements the go runtime function syscall/js.valueGet. It uses
// the Go abi.
//
// This has the effective type of:
//
//     func (w *wasmGo) goSyscallJSValueGet(ref, string) ref
func (w *wasmGo) goSyscallJSValueGet(sp int32) {
	name := string(w.goLoadSlice(sp + 16))
	result, ok := w.loadValue(sp + 8).(*Object).Props[name]
	if !ok {
		log.Panicf("ref doesn't have %s", name)
	}

	w.storeValue(sp+32, result)
}

// goSyscallJSValueSet implements the go runtime function syscall/js.valueSet. It uses
// the Go abi.
//
// This has the effective type of:
//
//     func (w *wasmGo) goSyscallJSValueSet(v ref, p string, x ref)
func (w *wasmGo) goSyscallJSValueSet(sp int32) {
	w.loadValue(sp + 8).(*Object).Props[w.loadString(sp+16)] = w.loadValue(sp + 32)
}

// goSyscallJSValueIndex implements the go runtime function syscall/js.valueIndex. It
// uses the Go abi.
//
// This has the effective type of:
//
//     func (w *wasmGo) goSyscallJSValueIndex(v ref, i int) ref
func (w *wasmGo) goSyscallJSValueIndex(sp int32) {
	result := (*w.loadValue(sp + 8).(*[]interface{}))[w.getInt64(sp+16)]
	w.storeValue(sp+24, result)
}

// goSyscallJSValueCall implements the go runtime function syscall/js.valueCall. It
// uses the Go abi.
//
// This has the effective type of:
//
//     func (w *wasmGo) goSyscallJSValueCall(v ref, m string, args []ref) (ref, bool)
func (w *wasmGo) goSyscallJSValueCall(sp int32) {
	// TODO error handling
	v := w.loadValue(sp + 8).(*Object)
	name := w.loadString(sp + 16)
	m, ok := v.Props[name]
	if !ok {
		panic("missing method: " + name) // TODO
	}
	args := w.loadSliceOfValues(sp + 32)
	result := m.(func(*Object, []interface{}) interface{})(v, args)
	// TODO getsp
	w.storeValue(sp+56, result)
	w.setUint8(sp+64, 1)
}

// goSyscallJSValueNew implements the go runtime function syscall/js.valueNew. It
// uses the Go abi.
//
// This has the effective type of:
//
//     func (w *wasmGo) goSyscallJSValueNew(v ref, args []ref) (ref, bool)
func (w *wasmGo) goSyscallJSValueNew(sp int32) {
	// TODO error handling
	v := w.loadValue(sp + 8)
	args := w.loadSliceOfValues(sp + 16)
	result := v.(*Object).New(args)
	// TODO getsp
	w.storeValue(sp+40, result)
	w.setUint8(sp+48, 1)
}

// goSyscallJSValueLength implements the go runtime function syscall/js.valueLength.
// It uses the Go abi.
//
// This has the effective type of:
//
//     func (w *wasmGo) goSyscallJSValueLength(v ref) int
func (w *wasmGo) goSyscallJSValueLength(sp int32) {
	array := w.loadValue(sp + 8).(*[]interface{})
	w.setInt64(sp+16, int64(len(*array)))
}

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

	_, err = w.writeMem(sp+8, data)
	if err != nil {
		panic(err)
	}
}

func (w *wasmGo) notImplemented(module, field string) goABIFunc {
	return func(sp int32) {
		w.vm.PrintStackTrace()
		log.Panicf("not implemented: %s %s", module, field)
	}
}

type goABIFunc func(int32)

var (
	// ImplementedGoABIFuncs is not used at runtime. This is used for documentation only.
	ImplementedGoABIFuncs = []string{
		"debug",                         // implemented
		"runtime.wasmExit",              // implemented
		"runtime.wasmWrite",             // implemented
		"runtime.nanotime",              // implemented
		"runtime.walltime",              // implemented
		"runtime.scheduleTimeoutEvent",  // stub
		"runtime.clearTimeoutEvent",     // stub
		"runtime.getRandomData",         // implemented
		"syscall/js.stringVal",          // stub
		"syscall/js.valueGet",           // implemented
		"syscall/js.valueSet",           // implemented
		"syscall/js.valueIndex",         // implemented
		"syscall/js.valueSetIndex",      // stub
		"syscall/js.valueCall",          // implemented
		"syscall/js.valueInvoke",        // stub
		"syscall/js.valueNew",           // implemented
		"syscall/js.valueLength",        // implemented
		"syscall/js.valuePrepareString", // stub
		"syscall/js.valueLoadString",    // stub
		"syscall/js.valueInstanceOf",    // stub
	}
)

func (w *wasmGo) runGoABI(doer func(int32)) exec.FunctionImport {
	return func(vm *exec.VirtualMachine) int64 {
		if w.vm == nil {
			w.vm = vm
		}

		f := vm.GetCurrentFrame()
		sp := int32(f.Locals[0])

		//log.Printf("%s(%d)", getFunctionName(doer), sp)
		doer(sp)

		return 0
	}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// ResolveGlobal does nothing, currently.
func (w *wasmGo) ResolveGlobal(module, field string) int64 { return 0 }

func (w *wasmGo) ResolveFunc(module, field string) exec.FunctionImport {
	log.Printf("wasmgo: resolving %s %s", module, field)

	val := w.ResolveFunc(module, field)
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
		case "runtime.getRandomData":
			return w.runGoABI(w.goRuntimeGetRandomData)
		case "syscall/js.valueGet":
			return w.runGoABI(w.goSyscallJSValueGet)
		case "syscall/js.valueSet":
			return w.runGoABI(w.goSyscallJSValueSet)
		case "syscall/js.valueIndex":
			return w.runGoABI(w.goSyscallJSValueIndex)
		case "syscall/js.valueCall":
			return w.runGoABI(w.goSyscallJSValueCall)
		case "syscall/js.valueNew":
			return w.runGoABI(w.goSyscallJSValueNew)
		case "syscall/js.valueLength":
			return w.runGoABI(w.goSyscallJSValueLength)
		default:
			log.Printf("unknown module+field %s %s, using shim", module, field)
			return w.runGoABI(w.notImplemented(module, field))
		}
	default:
		log.Panicf("unknown module+field %s %s", module, field)
	}

	panic("not implemented")
}
