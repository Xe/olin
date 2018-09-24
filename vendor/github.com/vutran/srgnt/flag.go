package srgnt

import (
	"flag"
	"time"
)

func (p *Program) AddBoolFlag(name string, value bool, usage string) {
	flag.Bool(name, value, usage)
}

func (p *Program) AddDurationFlag(name string, value time.Duration, usage string) {
	flag.Duration(name, value, usage)
}

func (p *Program) AddFloat64Flag(name string, value float64, usage string) {
	flag.Float64(name, value, usage)
}

func (p *Program) AddIntFlag(name string, value int, usage string) {
	flag.Int(name, value, usage)
}

func (p *Program) AddInt64Flag(name string, value int64, usage string) {
	flag.Int64(name, value, usage)
}

func (p *Program) AddStringFlag(name string, value string, usage string) {
	flag.String(name, value, usage)
}

func (p *Program) AddUint(name string, value uint, usage string) {
	flag.Uint(name, value, usage)
}

func (p *Program) AddUint64(name string, value uint64, usage string) {
	flag.Uint64(name, value, usage)
}
