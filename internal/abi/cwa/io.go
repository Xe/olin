package cwa

import (
	"math/rand"

	"within.website/olin/internal/fileresolver"
)

func (p *Process) IOGetStdin() int32 {
	fid := rand.Int31()
	p.FileHandles[fid] = fileresolver.Reader(p.Stdin, "stdin")
	return fid
}

func (p *Process) IOGetStdout() int32 {
	fid := rand.Int31()
	p.FileHandles[fid] = fileresolver.Writer(p.Stdout, "stdout")
	return fid
}

func (p *Process) IOGetStderr() int32 {
	fid := rand.Int31()
	p.FileHandles[fid] = fileresolver.Writer(p.Stderr, "stderr")
	return fid
}
