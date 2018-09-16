package cwa

func (p *Process) EnvGet(keyPtr, keyLen, valPtr, valLen uint32) (int32, error) {
	//p.logger.Printf("getenv: 0x%x:%d, 0x%x:%d", keyPtr, keyLen, valPtr, valLen)

	key := string(readMem(p.vm.Memory, keyPtr, keyLen))
	val, ok := p.env[key]
	if !ok {
		return 0, NotFoundError
	}

	if len(val) < int(valLen) {
		copy(p.vm.Memory[valPtr:], []byte(val))
	}

	return int32(len(val)), nil
}
