package fileresolver

import (
	"errors"
	"io"

	"within.website/olin/internal/abi"
)

// Errors exposed by partial.go
var (
	ErrNotImplemented = errors.New("fileresolver: not implemeted")
)

type readerFile struct {
	io.Reader
	name string
}

func (readerFile) Write([]byte) (int, error) { return -1, ErrNotImplemented }
func (readerFile) Flush() error              { return nil }
func (readerFile) Close() error              { return nil }
func (r readerFile) Name() string            { return r.name }

// Reader wraps an io.Reader as an abi.File.
func Reader(r io.Reader, name string) abi.File {
	return readerFile{Reader: r, name: name}
}

type writerFile struct {
	io.Writer
	name string
}

func (writerFile) Read([]byte) (int, error) { return -1, ErrNotImplemented }
func (writerFile) Flush() error             { return nil }
func (writerFile) Close() error             { return nil }
func (w writerFile) Name() string           { return w.name }

// Writer wraps an io.Writer as an abi.File.
func Writer(w io.Writer, name string) abi.File {
	return writerFile{Writer: w, name: name}
}
