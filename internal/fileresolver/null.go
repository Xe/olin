package fileresolver

import "github.com/Xe/olin/internal/abi"

// Null is a file that emulates /dev/null on Linux.
func Null() abi.File {
	return nullFile{}
}

type nullFile struct{}

func (nullFile) Write(p []byte) (int, error) { return len(p), nil }
func (nullFile) Read(p []byte) (int, error)  { return 0, nil }
func (nullFile) Sync() error                 { return nil }
func (nullFile) Close() error                { return nil }
func (nullFile) Name() string                { return "null" }
