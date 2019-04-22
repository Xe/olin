package cwa

func (p *Process) EnvGet(keyPtr, keyLen, valPtr, valLen uint32) (int32, error) {
	key := string(readMem(p.vm.Memory, keyPtr, keyLen))
	val, ok := p.env[key]
	if !ok {
		return 0, NotFoundError
	}

	if len(val) == int(valLen) {
		for i, by := range []byte(val) {
			p.vm.Memory[valPtr+uint32(i)] = by
		}
	}

	//p.logger.Printf("mem: 0x%x-0x%x: %x", valPtr, valPtr+valLen, p.vm.Memory[valPtr:valPtr+valLen])
	//p.logger.Printf("getenv: 0x%x %d (%q), 0x%x %d returns (%q) %d", keyPtr, keyLen, key, valPtr, valLen, val, len(val))

	return int32(len(val)), nil
}
