package wasmgo

import (
	"math"
	"syscall"
	"time"

	"github.com/Xe/olin/abi/cwa"
	"github.com/Xe/olin/internal/fileresolver"
	"github.com/perlin-network/life/exec"
)

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

func New(name string, argv []string, env map[string]string) *WasmGo {
	goObj := map[string]interface{}{
		"_makeFuncWrapper": func(this *Object, args []interface{}) interface{} {
			return &FuncWrapper{id: args[0]}
		},
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
				"write": func(this *Object, args []interface{}) interface{} {
					fd := int(args[0].(float64))
					buffer := args[1].(*TypedArray)
					offset := int(args[2].(float64))
					length := int(args[3].(float64))
					b := buffer.contents()[offset : offset+length]
					callback := args[5].(*FuncWrapper)

					if args[4] != nil {
						position := int64(args[4].(float64))
						_, err := syscall.Pwrite(fd, b, position)
						if err != nil {
							panic(err)
						}
					} else {
						_, err := syscall.Write(fd, b)
						if err != nil {
							panic(err)
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
