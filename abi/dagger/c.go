package dagger

func readMem(inp []byte, off uint32) []byte {
	var result []byte

	mem := inp[int(off):]
	for _, bt := range mem {
		if bt == 0 {
			return result
		}

		result = append(result, bt)
	}

	return result
}
