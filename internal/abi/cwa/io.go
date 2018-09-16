package cwa

import (
	"math/rand"

	"github.com/Xe/olin/internal/fileresolver"
)

func (p *Process) IOGetStdin() int32 {
	fid := rand.Int31()
	p.files[fid] = fileresolver.Reader(p.Stdin, "stdin")
	return fid
}

func (p *Process) IOGetStdout() int32 {
	fid := rand.Int31()
	p.files[fid] = fileresolver.Writer(p.Stdout, "stdout")
	return fid
}

func (p *Process) IOGetStderr() int32 {
	fid := rand.Int31()
	p.files[fid] = fileresolver.Writer(p.Stderr, "stderr")
	return fid
}
