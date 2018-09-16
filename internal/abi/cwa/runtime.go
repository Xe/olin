package cwa

import (
	"time"

	"github.com/Xe/olin/internal/names"
)

func (p *Process) specMajor() int32 {
	return names.CommonWASpecMajor
}

func (p *Process) specMinor() int32 {
	return names.CommonWASpecMinor
}

func (p *Process) RuntimeName(namePtr, nameLen uint32) int32 {
	if len(names.CommonWARuntimeName) < int(nameLen) {
		for i, by := range []byte(names.CommonWARuntimeName) {
			p.vm.Memory[namePtr+uint32(i)] = by
		}
	}
	return int32(len(names.CommonWARuntimeName))
}

func (p *Process) msleep(ms int32) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
