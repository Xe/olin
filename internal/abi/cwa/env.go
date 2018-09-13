package cwa

func (p *Process) envGet(keyPtr, keyLen, valPtr, valLen uint32) (int32, error) {
	key := string(readMem(p.vm.Memory, keyPtr, keyLen))
	val, ok := p.env[key]
	if !ok {
		return 0, NotFoundError
	}

	if len(val) < int(valLen) {
		for i, by := range []byte(val) {
			p.vm.Memory[valPtr+uint32(i)] = by
		}
	}

	return int32(len(val)), nil
}
