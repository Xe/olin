package dagger

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/Xe/olin/internal/abi"
	"github.com/Xe/olin/internal/fileresolver"
	"github.com/perlin-network/life/exec"
)

// Process is a higher level wrapper around a set of files for dagger
// modules.
type Process struct {
	name  string
	files []abi.File
}

// Files returns the process' list of open files.
func (p Process) Files() []abi.File {
	return p.files
}

// NewProcess creates a new process.
func NewProcess(name string) *Process {
	return &Process{name: name}
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

// OpenFD opens a file descriptor for this Process with the given file url
// string and flags integer.
func (p *Process) OpenFD(furl string, flags uint32) int64 {
	fd, err := p.open(furl, flags)

	if err != nil {
		log.Printf("%s: OpenFD(%s, %x): %v", p.name, furl, flags, err)
		return int64(-1 * int64(err.(Error).Errno))
	}

	return int64(fd)
}

// CloseFD closes a file descriptor.
func (p *Process) CloseFD(fd int64) int64 {
	err := p.files[fd].Close()
	if err != nil {
		log.Printf("%s: CloseFD(%d): %v", p.name, fd, err)
		return -1
	}

	p.files[fd] = nil

	return 0
}

// WriteFD writes the given data to a file descriptor, returning -1 if it failed
// somehow.
func (p *Process) WriteFD(fd int64, data []byte) int64 {
	n, err := p.files[int(fd)].Write(data)
	if err != nil {
		log.Printf("%s: WriteFD(%d, []byte{%d}): %v", p.name, fd, len(data), err)
		return -1
	}

	return int64(n)
}

// SyncFD runs a file's sync operation and returns -1 if it failed.
func (p *Process) SyncFD(fd int64) int64 {
	err := p.files[fd].Flush()
	if err != nil {
		log.Printf("%s: Sync(%d) %v", p.name, fd, err)
		return -1
	}

	return 0
}

// ReadFD reads data from the given FD into the byte buffer.
func (p *Process) ReadFD(fd int64, buf []byte) int64 {
	n, err := p.files[fd].Read(buf)
	if err != nil {
		log.Printf("%s: ReadFD(%d, []byte{%d}): %v", p.name, fd, len(buf), err)
		return -1
	}

	return int64(n)
}

// ResolveFunc resolves dagger's ABI and importable functions.
func (p *Process) ResolveFunc(module, field string) exec.FunctionImport {
	switch module {
	case "dagger":
		switch field {
		case "open": // :: String -> Int32 -> Int64
			return func(vm *exec.VirtualMachine) int64 {
				f := vm.GetCurrentFrame()
				furlPtr := uint32(f.Locals[0])
				flags := uint32(f.Locals[1])
				furl := string(readMem(vm.Memory, furlPtr))

				return p.OpenFD(furl, flags)
			}
		case "close": // :: Int64 -> IO Int64
			return func(vm *exec.VirtualMachine) int64 {
				f := vm.GetCurrentFrame()
				fd := f.Locals[0]

				return p.CloseFD(fd)
			}
		case "write": // :: Int64 -> String -> IO Int64
			return func(vm *exec.VirtualMachine) int64 {
				f := vm.GetCurrentFrame()
				fd := f.Locals[0]
				ptr := f.Locals[1]
				len := f.Locals[2]
				mem := vm.Memory[int(ptr):int(ptr+len)]

				return p.WriteFD(fd, mem)
			}
		case "sync": // :: Int64 -> IO Int64
			return func(vm *exec.VirtualMachine) int64 {
				f := vm.GetCurrentFrame()
				fd := f.Locals[0]

				return p.SyncFD(fd)
			}
		case "read": // :: Int64 -> String -> IO Int64
			return func(vm *exec.VirtualMachine) int64 {
				f := vm.GetCurrentFrame()
				fd := f.Locals[0]
				ptr := int32(f.Locals[1])
				len := f.Locals[2]
				buf := make([]byte, int(len))
				ret := p.ReadFD(fd, buf)

				for i, d := range buf {
					vm.Memory[ptr+int32(i)] = d
				}

				return ret
			}
		}
	}

	return nil
}

// Open makes this Process track an arbitrary extra file.
func (p *Process) Open(f abi.File) {
	p.insertFile(f)
}

// ResolveGlobal does nothing, currently.
func (p *Process) ResolveGlobal(module, field string) int64 { return 0 }

func (p *Process) open(furl string, flags uint32) (int, error) {
	u, err := url.Parse(furl)
	if err != nil {
		return -1, makeError(ErrorBadURL, err)
	}

	var file abi.File
	switch u.Scheme {
	case "log":
		q := u.Query()
		file = fileresolver.Log(os.Stdout, q.Get("prefix"), log.LstdFlags)

	case "fd":
		fdNum, err := strconv.Atoi(u.Host)
		if err != nil {
			return -1, makeError(ErrorBadURLInput, err)
		}

		file = fileresolver.NewOSFile(uintptr(fdNum), u.Host)

	case "http", "https":
		file, _ = fileresolver.HTTP(&http.Client{}, u)

	default:
		return -1, makeError(ErrorUnknownScheme, fmt.Errorf("dagger: open: unknown scheme %s", u.Scheme))
	}

	fd := p.insertFile(file)

	return fd, nil
}
