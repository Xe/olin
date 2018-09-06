package fileresolver

import (
	"io"
	"log"

	"github.com/Xe/olin/internal/abi"
)

// Log returns a file that redirects all logging calls to a standard library logger
// with the given prefix and flags.
func Log(writer io.Writer, prefix string, flag int) abi.File {
	l := log.New(writer, prefix, flag)
	return logFile{
		l:      l,
		prefix: prefix,
	}
}

type logFile struct {
	prefix string
	l      *log.Logger
}

func (logFile) Read(p []byte) (int, error) { return 0, nil }
func (logFile) Sync() error                { return nil }
func (logFile) Close() error               { return nil }
func (l logFile) Name() string             { return l.prefix }

func (l logFile) Write(p []byte) (int, error) {
	res := len(string(p))
	l.l.Println(string(p))

	return res, nil
}
