package cwa

func (p *Process) EnvGet(keyPtr, keyLen, valPtr, valLen uint32) (int32, error) {
	val, ok := p.env[string(readMem(p.vm.Memory, keyPtr, keyLen))]
	if !ok {
		return 0, NotFoundError
	}

	if len(val) < int(valLen) {
		copy(p.vm.Memory[valPtr:], []byte(val))
	}

	//p.logger.Printf("mem: 0x%x-0x%x: %x", valPtr, valPtr+valLen, p.vm.Memory[valPtr:valPtr+valLen])

	return int32(len(val)), nil
}
