package cwa

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Xe/olin/internal/abi"
	"github.com/perlin-network/life/exec"
)

// NewProcess creates a new process with the given name, arguments and environment.
func NewProcess(name string, argv []string, env map[string]string) *Process {
	return &Process{
		name:   name,
		hc:     &http.Client{},
		logger: log.New(os.Stdout, name+": ", log.LstdFlags),
		env:    env,
		argv:   argv,
		files:  map[int32]abi.File{},

		Stdin:  bytes.NewBuffer(nil),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Process is an individual CommonWA process. It is the collection of resources
// and other macguffins that the child module ends up requiring.
type Process struct {
	name string

	hc     *http.Client
	logger *log.Logger
	env    map[string]string
	vm     *exec.VirtualMachine
	argv   []string
	files  map[int32]abi.File

	Stdin          io.Reader
	Stdout, Stderr io.Writer
}

// Name returns this process' name.
func (p *Process) Name() string { return p.name }

// SetVM sets the VM associated with this process.
func (p *Process) SetVM(vm *exec.VirtualMachine) { p.vm = vm }

// Open does nothing
func (Process) Open(abi.File) {}

// Files returns the set of open files in use by this process.
func (p *Process) Files() []abi.File {
	var result []abi.File

	for _, fi := range p.files {
		result = append(result, fi)
	}

	return result
}

// ResolveGlobal does nothing, currently.
func (p *Process) ResolveGlobal(module, field string) int64 { return 0 }

// ResolveFunc resolves the CommonWA ABI and importable functions.
func (p *Process) ResolveFunc(module, field string) exec.FunctionImport {
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

				p.log(level, msgPtr, msgLen)

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

				result, err := p.envGet(keyPtr, keyLen, valPtr, valLen)
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

				return int64(p.runtimeName(namePtr, nameLen))
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

				return int64(p.argLen())
			}
		case "startup_arg_at":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				f := vm.GetCurrentFrame()
				i := int32(f.Locals[0])
				outPtr := uint32(f.Locals[1])
				outLen := uint32(f.Locals[2])

				result, err := p.argAt(i, outPtr, outLen)
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

				result, err := p.open(urlPtr, urlLen)
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

				result, err := p.read(fid, dataPtr, dataLen)
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

				result, err := p.write(fid, dataPtr, dataLen)
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

				err := p.close(fid)
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

				err := p.flush(fid)
				if err != nil {
					return int64(ErrorCode(err))
				}

				return 0
			}
		case "time_now":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				return p.timeNow()
			}
		case "io_get_stdin":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				return int64(p.ioGetStdin())
			}
		case "io_get_stdout":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				return int64(p.ioGetStdout())
			}
		case "io_get_stderr":
			return func(vm *exec.VirtualMachine) int64 {
				p.SetVM(vm)

				return int64(p.ioGetStderr())
			}
		default:
			log.Panicf("unknown import %s::%s", module, field)
		}
	default:
		log.Panicf("unknown module %s", module)
	}

	return nil
}
