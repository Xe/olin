package dagger

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/Xe/olin/internal/abi"
	"github.com/perlin-network/life/exec"
)

// Process is a higher level wrapper around a set of files for dagger
// modules.
type Process struct {
	abi.ABI

	name  string
	files []abi.File
}

// insertFile adds a file to the set of files and returns its descriptor.
func (p *Process) insertFile(file abi.File) int {
	for i, f := range p.files {
		if f == nil {
			p.files[i] = file
			return i
		}
	}

	i := len(p.files)
	p.files = append(p.files, file)
	return i
}

// Name returns this process's name.
func (p *Process) Name() string { return p.name }

func (p *Process) ResolveFunc(module, field string) exec.FunctionImport {
	switch module {
	case "dagger":
		switch field {
		case "open":
			return func(vm *exec.VirtualMachine) int64 {
				f := vm.GetCurrentFrame()
				furlPtr := uint32(f.Locals[0])
				flags := uint32(f.Locals[1])
				furl := string(readMem(vm.Memory, furlPtr))

				fd, err := p.open(furl, flags)
				if err != nil {
					return int64(-1 * int64(err.(Error).Errno))
				}

				return int64(fd)
			}
		case "write":
			return func(vm *exec.VirtualMachine) int64 {
				f := vm.GetCurrentFrame()
				fd := f.Locals[0]
				ptr := f.Locals[1]
				len := f.Locals[2]

				mem := vm.Memory[int(ptr):int(ptr+len)]

				n, err := p.files[int(fd)].Write(mem)
				if err != nil {
					return -1
				}

				return int64(n)
			}
		}
	}

	return nil
}

func (p *Process) ResolveGlobal(module, field string) int64 { return 0 }

func (p *Process) open(furl string, flags uint32) (int, error) {
	u, err := url.Parse(furl)
	if err != nil {
		return -1, makeError(ErrorBadURL, err)
	}

	var file abi.File
	switch u.Scheme {
	case "fd":
		fdNum, err := strconv.Atoi(u.Host)
		if err != nil {
			return -1, makeError(ErrorBadURLInput, err)
		}

		file = abi.NewOSFile(uintptr(fdNum), u.Host)

	default:
		return -1, makeError(ErrorUnknownScheme, fmt.Errorf("dagger: open: unknown scheme %s", u.Scheme))
	}

	fd := p.insertFile(file)

	return fd, nil
}

func (p *Process) close(fd int) int {
	return 0
}
