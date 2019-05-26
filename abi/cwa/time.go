package cwa

import "time"

func (p *Process) TimeNow() int64 {
	return time.Now().UTC().Unix()
}
