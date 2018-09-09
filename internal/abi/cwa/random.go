package cwa

import (
	"math/rand"
)

func (p *Process) randI32() int32 {
	return rand.Int31()
}

func (p *Process) randI64() int64 {
	return rand.Int63()
}
