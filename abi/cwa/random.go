package cwa

import (
	"math/rand"
)

func (p *Process) RandI32() int32 {
	return rand.Int31()
}

func (p *Process) RandI64() int64 {
	return rand.Int63()
}
