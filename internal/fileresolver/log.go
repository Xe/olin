package fileresolver

import (
	"io"
	"log"

	"github.com/Xe/olin/internal/abi"
)

// Log returns a file that redirects all write calls to a standard library logger
// with the given prefix and flags.
//
// For more information, please see the spec here: https://github.com/CommonWA/cwa-spec/blob/master/schemes/log.md
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
func (logFile) Flush() error               { return nil }
func (logFile) Close() error               { return nil }
func (l logFile) Name() string             { return l.prefix }

func (l logFile) Write(p []byte) (int, error) {
	res := len(string(p))
	l.l.Println(string(p))

	return res, nil
}
