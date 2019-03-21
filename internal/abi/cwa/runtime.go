package cwa

import (
	"time"
)

func (p *Process) specMajor() int32 {
	return SpecMajor
}

func (p *Process) specMinor() int32 {
	return SpecMinor
}

func (p *Process) RuntimeName(namePtr, nameLen uint32) int32 {
	if len(RuntimeName) < int(nameLen) {
		for i, by := range []byte(RuntimeName) {
			p.vm.Memory[namePtr+uint32(i)] = by
		}
	}
	return int32(len(RuntimeName))
}

func (p *Process) msleep(ms int32) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
