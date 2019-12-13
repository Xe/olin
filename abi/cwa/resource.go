package cwa

import (
	"io"
	"log"
	"math/rand"
	"net/url"
	"os"

	"within.website/olin/abi"
	"within.website/olin/fileresolver"
)

func (p *Process) checkPolicy(uri string) bool {
	if p.Policy == nil {
		return true // yolo
	}

	for _, allow := range p.Policy.Allowed {
		if allow.MatchString(uri) {
			return true
		}
	}

	for _, disallow := range p.Policy.Disallowed {
		if disallow.MatchString(uri) {
			return false
		}
	}

	return false
}

func (p *Process) ResourceOpen(urlPtr, urlLen uint32) (int32, error) {
	u := string(readMem(p.vm.Memory, urlPtr, urlLen))
	uu, err := url.Parse(u)
	if err != nil {
		p.Logger.Printf("can't parse url %s: %v, returning:  %v", u, err, InvalidArgumentError)
		return 0, InvalidArgumentError
	}

	if !p.checkPolicy(u) {
		p.Logger.Printf("vibe check failed: %s forbidden by policy", u)
		p.vm.Exited = true
		p.vm.ReturnValue = -1
		return 0, PermissionDeniedError
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
		file, err = fileresolver.HTTP(p.HC, uu)
		if err != nil {
			p.Logger.Printf("can't resource_open(%q): %v", u, err)
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

	//p.Logger.Printf("writing %d bytes to %d (%s)", dataLen, fid, f.Name())

	n, err := f.Write(mem)
	if err != nil {
		p.Logger.Printf("write error for fid %d (%s): %v", fid, f.Name(), err)
		if err == io.EOF {
			return -5, EndOfFileError
		}
		return 0, UnknownError
	}

	return int32(n), nil
}

func (p *Process) ResourceRead(fid int32, dataPtr, dataLen uint32) (int32, error) {
	f, ok := p.FileHandles[fid]
	if !ok {
		return 0, InvalidArgumentError
	}

	//p.Logger.Printf("reading %d bytes from %d (%s)", dataLen, fid, f.Name())

	outp := make([]byte, int(dataLen))
	n, err := f.Read(outp)
	if err != nil {
		p.Logger.Printf("read error for fid %d (%s): %v", fid, f.Name(), err)
		if err == io.EOF {
			return 0, EndOfFileError
		}
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
		p.Logger.Printf("close error for fid %d (%s): %v", fid, f.Name(), err)
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
		p.Logger.Printf("flush error for fid %d (%s): %v", fid, f.Name(), err)
		return UnknownError
	}

	return nil
}
