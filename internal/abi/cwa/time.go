package cwa

import "time"

func (p *Process) timeNow() int64 {
	return time.Now().UTC().Unix()
}
