package cwa

func (p *Process) ArgLen() int32 {
	return int32(len(p.argv))
}

func (p *Process) ArgAt(i int32, outPtr, outLen uint32) (int32, error) {
	if i > int32(len(p.argv)) {
		return 0, InvalidArgumentError
	}

	arg := p.argv[i]
	if len(arg) < int(outLen) {
		for i, by := range []byte(arg) {
			p.vm.Memory[outPtr+uint32(i)] = by
		}
	}

	return int32(len(arg)), nil
}
