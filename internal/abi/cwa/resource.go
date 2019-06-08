package cwa

import (
	"log"
	"math/rand"
	"net/url"
	"os"

	"within.website/olin/internal/abi"
	"within.website/olin/internal/fileresolver"
)

func (p *Process) ResourceOpen(urlPtr, urlLen uint32) (int32, error) {
	u := string(readMem(p.vm.Memory, urlPtr, urlLen))
	uu, err := url.Parse(u)
	if err != nil {
		p.logger.Printf("can't parse url %s: %v, returning:  %v", u, err, InvalidArgumentError)
		return 0, InvalidArgumentError
	}

	q := uu.Query()
	var file abi.File
	switch uu.Scheme {
	case "log":
		prefix := q.Get("prefix")
		file = fileresolver.Log(os.Stdout, p.name+": "+prefix, log.LstdFlags)
	case "random":
		file = fileresolver.Random()
	case "null":
		file = fileresolver.Null()
	case "zero":
		file = fileresolver.Zero()
	case "http", "https":
		var err error
		file, err = fileresolver.HTTP(p.hc, uu)
		if err != nil {
			p.logger.Printf("can't resource_open(%q): %v", u, err)
			return 0, UnknownError
		}
	default:
		return 0, NotFoundError
	}

	fid := rand.Int31()
	p.FileHandles[fid] = file

	return fid, nil
}

func (p *Process) ResourceWrite(fid int32, dataPtr, dataLen uint32) (int32, error) {
	mem := p.vm.Memory[dataPtr : dataPtr+dataLen]

	f, ok := p.FileHandles[fid]
	if !ok {
		return 0, InvalidArgumentError
	}

	//p.logger.Printf("writing %d bytes to %d (%s)", dataLen, fid, f.Name())

	n, err := f.Write(mem)
	if err != nil {
		p.logger.Printf("write error for fid %d (%s): %v", fid, f.Name(), err)
		return 0, UnknownError
	}

	return int32(n), nil
}

func (p *Process) ResourceRead(fid int32, dataPtr, dataLen uint32) (int32, error) {
	f, ok := p.FileHandles[fid]
	if !ok {
		return 0, InvalidArgumentError
	}

	//p.logger.Printf("reading %d bytes from %d (%s)", dataLen, fid, f.Name())

	outp := make([]byte, int(dataLen))
	n, err := f.Read(outp)
	if err != nil {
		p.logger.Printf("read error for fid %d (%s): %v", fid, f.Name(), err)
		return 0, UnknownError
	}

	for i, d := range outp {
		p.vm.Memory[dataPtr+uint32(i)] = d
	}

	return int32(n), nil
}

func (p *Process) ResourceClose(fid int32) error {
	f, ok := p.FileHandles[fid]
	if !ok {
		return InvalidArgumentError
	}

	err := f.Close()
	if err != nil {
		p.logger.Printf("close error for fid %d (%s): %v", fid, f.Name(), err)
		return UnknownError
	}

	delete(p.FileHandles, fid)

	return nil
}

func (p *Process) ResourceFlush(fid int32) error {
	f, ok := p.FileHandles[fid]
	if !ok {
		return InvalidArgumentError
	}

	err := f.Flush()
	if err != nil {
		p.logger.Printf("flush error for fid %d (%s): %v", fid, f.Name(), err)
		return UnknownError
	}

	return nil
}
