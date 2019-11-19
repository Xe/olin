package cwa

func readMem(inp []byte, off uint32, max uint32) []byte {
	var result []byte

	mem := inp[int(off):]
	for i, bt := range mem {
		if uint32(i) == max {
			return result
		}
		if bt == 0 {
			return result
		}

		result = append(result, bt)
	}

	return result
}
