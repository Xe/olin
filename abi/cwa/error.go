package cwa

import "strconv"

//go:generate stringer -type=Error

// Error is an individual error as defined by the CommonWA spec.
type Error int

// CommonWA errors as defined by the spec at https://github.com/CommonWA/cwa-spec/blob/master/errors.md
const (
	ErrNone Error = (iota * -1)
	UnknownError
	InvalidArgumentError
	PermissionDeniedError
	NotFoundError
)

func (e Error) Error() string {
	return "cwa: error " + strconv.Itoa(int(e)) + " " + e.String()
}

// ErrorCode extracts the code from an error.
func ErrorCode(err error) int {
	val, _ := err.(Error)
	return int(val)
}
