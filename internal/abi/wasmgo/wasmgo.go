package wasmgo

import (
	"log"
	"math"
	"syscall"
	"time"

	"github.com/Xe/olin/internal/abi/cwa"
	"github.com/Xe/olin/internal/fileresolver"
	"github.com/perlin-network/life/exec"
)

// WasmGo is a wrapper for the Go runtime.
type WasmGo struct {
	*cwa.Process

	BootTime   time.Time
	Exited     bool
	StatusCode int32
	Memory     *ArrayBuffer

	vm *exec.VirtualMachine

	values []interface{}
	refs   map[interface{}]int
}

// New builds a new instance of the Go glue for interop.
//
// See https://github.com/Xe/olin/issues/69 for more information about
// compatibility with go_wasm_js_exec.js.
func New(name string, argv []string, env map[string]string) *WasmGo {
	goObj := map[string]interface{}{
		"_makeFuncWrapper": func(this *Object, args []interface{}) interface{} {
			return &FuncWrapper{id: args[0]}
		},

		// This object seems to be the key to how callbacks/timeouts work in wasm+js.
		// When an event has been processed by the JavaScript side of the wasm+js Go
		// support and is ready for Go to handle it, it will be put here. I wonder if
		// we can have a channel feed into this. This may need the Go object to be a
		// bit more elaborate than just a map of strings to void pointers though.
		"_pendingEvent": nil,
	}

	w := &WasmGo{
		Memory:   &ArrayBuffer{},
		Process:  cwa.NewProcess(name, argv, env),
		BootTime: time.Now(),
		refs:     make(map[interface{}]int),
	}

	w.Process.FileHandles[0] = fileresolver.NewOSFile(uintptr(syscall.Stdin), "stdin")
	w.Process.FileHandles[1] = fileresolver.NewOSFile(uintptr(syscall.Stdout), "stdout")
	w.Process.FileHandles[2] = fileresolver.NewOSFile(uintptr(syscall.Stderr), "stderr")

	w.values = []interface{}{
		math.NaN(),
		float64(0),
		nil,
		true,
		false,
		&Object{Props: map[string]interface{}{
			"Object": &Object{
				New: func(args []interface{}) interface{} {
					panic("new Object")
				},
			},
			"Array": &Object{
				New: func(args []interface{}) interface{} {
					panic("new Array")
				},
			},
			"Int8Array":    typedArrayClass,
			"Int16Array":   typedArrayClass,
			"Int32Array":   typedArrayClass,
			"Uint8Array":   typedArrayClass,
			"Uint16Array":  typedArrayClass,
			"Uint32Array":  typedArrayClass,
			"Float32Array": typedArrayClass,
			"Float64Array": typedArrayClass,
			"process":      &Object{},
			"fs": &Object{Props: map[string]interface{}{
				"constants": &Object{Props: map[string]interface{}{
					"O_WRONLY": syscall.O_WRONLY,
					"O_RDWR":   syscall.O_RDWR,
					"O_CREAT":  syscall.O_CREAT,
					"O_TRUNC":  syscall.O_TRUNC,
					"O_APPEND": syscall.O_APPEND,
					"O_EXCL":   syscall.O_EXCL,
				}},
				// open(path, flags, mode, callback) {
				"open": func(this *Object, args []interface{}) interface{} {
					log.Printf("%#v", args[0])
					path := args[0].(string)
					_ = int(args[1].(float64)) // flags are ignored
					_ = int(args[2].(float64)) // mode is ignored
					callback := args[3].(*FuncWrapper)

					fd, err := w.OpenFile(path)
					if err != nil {
						panic(err)
					}

					goObj["_pendingEvent"] = &Object{Props: map[string]interface{}{
						"id":   callback.id,
						"this": nil,
						"args": &[]interface{}{
							nil,
							float64(fd),
						},
					}}

					return nil
				},
				// write(fd, buf, offset, length, position, callback) {
				"write": func(this *Object, args []interface{}) interface{} {
					fd := int32(args[0].(float64))
					buffer := args[1].(*TypedArray)
					offset := int(args[2].(float64))
					length := int(args[3].(float64))
					b := buffer.contents()[offset : offset+length]
					callback := args[5].(*FuncWrapper)

					if args[4] != nil {
						// TODO: pwrite(2)/lseek(2) support
						panic("no pwrite support :(")
					} else {
						if f, ok := w.Process.FileHandles[fd]; !ok {
							log.Panicf("can't find fd %d", fd)
						} else {
							_, err := f.Write(b)
							if err != nil {
								panic(err)
							}
						}
					}

					goObj["_pendingEvent"] = &Object{Props: map[string]interface{}{
						"id":   callback.id,
						"this": nil,
						"args": &[]interface{}{
							nil,
							length,
						},
					}}
					return nil
				},
			}},
		}}, // global
		&Object{Props: map[string]interface{}{
			"buffer": w.Memory,
		}}, // memory
		&Object{Props: goObj}, // go
	}

	return w
}

// TODO(Xe): replace with copy() or obviate the need to copy to begin with?
// while this code is being run, the WebAssembly memory is considered locked
// because this is a single-threaded environment at the moment.
func (w *WasmGo) writeMem(ptr int32, data []byte) (int, error) {
	ctr := 0
	for i, d := range data {
		w.vm.Memory[ptr+int32(i)] = d
		ctr++
	}

	return ctr, nil
}
