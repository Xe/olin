package dagger

import "fmt"

// Error is a common error type for dagger operations.
type Error struct {
	Errno      Errno
	Underlying error
}

func (e Error) Error() string {
	return fmt.Sprintf("dagger: error code %d (%s): %v", e.Errno, e.Errno.String(), e.Underlying)
}

func makeError(no Errno, err error) Error {
	return Error{
		Errno:      no,
		Underlying: err,
	}
}

// Errno is the error number for an error.
type Errno int

//go:generate stringer -type=Errno

// Error numbers
const (
	ErrorNone Errno = iota
	ErrorBadURL
	ErrorBadURLInput
	ErrorUnknownScheme
)
