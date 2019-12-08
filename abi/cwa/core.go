package cwa

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/perlin-network/life/exec"
	"within.website/olin/abi"
)

// NewProcess creates a new process with the given name, arguments and environment.
func NewProcess(name string, argv []string, env map[string]string) *Process {
	return &Process{
		name:   name,
		hc:     &http.Client{},
		Logger: log.New(os.Stdout, name+": ", log.LstdFlags),
		env:    env,
		argv:   argv,

		FileHandles: map[int32]abi.File{},
		Stdin:       bytes.NewBuffer([]byte("")),
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
	}
}

// Process is an individual CommonWA process. It is the collection of resources
// and other macguffins that the child module ends up requiring.
type Process struct {
	name string

	hc           *http.Client
	Logger       *log.Logger
	env          map[string]string
	vm           *exec.VirtualMachine
	argv         []string
	syscallCount int64

	FileHandles    map[int32]abi.File
	Stdin          io.Reader
	Stdout, Stderr io.Writer
}

// Setenv updates a process' environment. This does not inform the process. It
// is up to the running process to detect these values have changed.
func (p *Process) Setenv(m map[string]string) { p.env = m }

// Name returns this process' name.
func (p *Process) Name() string { return p.name }

// SetVM sets the VM associated with this process.
func (p *Process) SetVM(vm *exec.VirtualMachine) { p.vm = vm }

// Open does nothing
func (Process) Open(abi.File) {}

func (p Process) SyscallCount() int64 {
	return p.syscallCount
}

// Files returns the set of open files in use by this process.
func (p *Process) Files() []abi.File {
	var result []abi.File

	for _, fi := range p.FileHandles {
		result = append(result, fi)
	}

	return result
}

// ResolveGlobal does nothing, currently.
func (p *Process) ResolveGlobal(module, field string) int64 { return 0 }

// ResolveFunc resolves the CommonWA ABI and importable functions.
func (p *Process) ResolveFunc(module, field string) exec.FunctionImport {
	p.syscallCount++
	switch module {
	case "cwa", "env":
		switch field {
		case "log_write":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				f := vm.GetCurrentFrame()
				level := int32(f.Locals[0])
				msgPtr := uint32(f.Locals[1])
				msgLen := uint32(f.Locals[2])

				p.LogWrite(level, msgPtr, msgLen)

				return 0
			}
		case "env_get":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				f := vm.GetCurrentFrame()
				keyPtr := uint32(f.Locals[0])
				keyLen := uint32(f.Locals[1])
				valPtr := uint32(f.Locals[2])
				valLen := uint32(f.Locals[3])

				result, err := p.EnvGet(keyPtr, keyLen, valPtr, valLen)
				if err != nil {
					return int64(ErrorCode(err))
				}

				return int64(result)
			}
		case "runtime_spec_major":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				return int64(p.specMajor())
			}
		case "runtime_spec_minor":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				return int64(p.specMinor())
			}
		case "runtime_name":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				f := vm.GetCurrentFrame()
				namePtr := uint32(f.Locals[0])
				nameLen := uint32(f.Locals[1])

				return int64(p.RuntimeName(namePtr, nameLen))
			}
		case "runtime_msleep":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				f := vm.GetCurrentFrame()
				ms := int32(f.Locals[0])

				p.msleep(ms)

				return 0
			}
		case "startup_arg_len":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				return int64(p.ArgLen())
			}
		case "startup_arg_at":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				f := vm.GetCurrentFrame()
				i := int32(f.Locals[0])
				outPtr := uint32(f.Locals[1])
				outLen := uint32(f.Locals[2])

				result, err := p.ArgAt(i, outPtr, outLen)
				if err != nil {
					return int64(ErrorCode(err))
				}

				return int64(result)
			}
		case "resource_open":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				f := vm.GetCurrentFrame()
				urlPtr := uint32(f.Locals[0])
				urlLen := uint32(f.Locals[1])

				result, err := p.ResourceOpen(urlPtr, urlLen)
				if err != nil {
					return int64(ErrorCode(err))
				}

				return int64(result)
			}
		case "resource_read":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				f := vm.GetCurrentFrame()
				fid := int32(f.Locals[0])
				dataPtr := uint32(f.Locals[1])
				dataLen := uint32(f.Locals[2])

				result, err := p.ResourceRead(fid, dataPtr, dataLen)
				if err != nil {
					return int64(ErrorCode(err))
				}

				return int64(result)
			}
		case "resource_write":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				f := vm.GetCurrentFrame()
				fid := int32(f.Locals[0])
				dataPtr := uint32(f.Locals[1])
				dataLen := uint32(f.Locals[2])

				result, err := p.ResourceWrite(fid, dataPtr, dataLen)
				if err != nil {
					return int64(ErrorCode(err))
				}

				return int64(result)
			}
		case "resource_close":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				f := vm.GetCurrentFrame()
				fid := int32(f.Locals[0])

				err := p.ResourceClose(fid)
				if err != nil {
					return int64(ErrorCode(err))
				}

				return 0
			}
		case "resource_flush":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				f := vm.GetCurrentFrame()
				fid := int32(f.Locals[0])

				err := p.ResourceFlush(fid)
				if err != nil {
					return int64(ErrorCode(err))
				}

				return 0
			}
		case "time_now":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				return p.TimeNow()
			}
		case "io_get_stdin":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				return int64(p.IOGetStdin())
			}
		case "io_get_stdout":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				return int64(p.IOGetStdout())
			}
		case "io_get_stderr":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				return int64(p.IOGetStderr())
			}
		case "random_i32":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				return int64(p.RandI32())
			}
		case "random_i64":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				return p.RandI64()
			}
		}
	}

	return nil
}
