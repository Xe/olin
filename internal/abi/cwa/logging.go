package cwa

import "strconv"

func (p *Process) log(level int32, msgPtr, msgLen uint32) {
	msg := string(readMem(p.vm.Memory, msgPtr, msgLen))
	p.logger.Printf("%s: %s", levelToString(level), msg)
}

// Log levels
const (
	LogLevelError   = 1
	LogLevelWarning = 3
	LogLevelInfo    = 6
)

func levelToString(level int32) string {
	switch level {
	case LogLevelError:
		return "error"
	case LogLevelWarning:
		return "warning"
	case LogLevelInfo:
		return "info"
	}

	return "level " + strconv.Itoa(int(level))
}
