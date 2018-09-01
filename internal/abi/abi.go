// Package abi is the common interface for all ABI implementations
package abi

import (
	"io"
	"net/url"

	"github.com/perlin-network/life/exec"
)

// Namer is something that has a name.
type Namer interface {
	Name() string
}

// ABI is the interface that is provided to a webaseembly module. Instances
// of this interface will live for the lifetime of webassembly modules.
//
// This is low level intentionally.
type ABI interface {
	exec.ImportResolver
	Namer

	Files() []File
	Open(File)
}

// File is the most common denominator for most of what you usefully need out of
// files to make useful programs.
type File interface {
	io.Closer
	io.Reader
	io.Writer
	Namer

	// Sync isn't required for all backends, but it is used for some, such as HTTP.
	Sync() error
}

// FileOpener opens a given file by URL.
type FileOpener interface {
	Open(furl *url.URL, flags int) (File, error)
}
