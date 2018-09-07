package fileresolver

import (
	"syscall"

	"github.com/Xe/olin/internal/abi"
)

// OSFile implements File backed by a raw operating system file.
//
// This doesn't buffer I/O.
type OSFile struct {
	fd   uintptr
	name string
}

// NewOSFile creates a file backed by an OS file descriptor without buffering.
// This is kinda dangerous. Only do this if you know what you are doing.
func NewOSFile(fd uintptr, name string) abi.File {
	return OSFile{fd: fd, name: name}
}

// Sync does nothing.
func (OSFile) Sync() error { return nil }

// Name returns the file's name.
func (o OSFile) Name() string { return o.name }

// Read reads data from the OS file.
func (o OSFile) Read(p []byte) (int, error) {
	return syscall.Read(int(o.fd), p)
}

// Write writes data to the OS file.
func (o OSFile) Write(p []byte) (int, error) {
	return syscall.Write(int(o.fd), p)
}

// Close closes this OS file.
func (o OSFile) Close() error {
	return syscall.Close(int(o.fd))
}
