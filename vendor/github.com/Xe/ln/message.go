package ln

import "fmt"

func withLevel(level, format string, args ...interface{}) Fer {
	return F{level: fmt.Sprintf(format, args...), "level": level}
}

// Info adds an informational connotation to this log line. This is the 
// "default" state for things that can be ignored.
func Info(format string, args ...interface{}) Fer {
	return withLevel("info", format, args...)
}

// Debug adds a debugging connotation to this log line. This may be ignored
// or aggressively sampled in order to save ram.
func Debug(format string, args ...interface{}) Fer {
	return withLevel("debug", format, args...)
}