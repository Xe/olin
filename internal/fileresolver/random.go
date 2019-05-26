package fileresolver

import (
	"crypto/rand"

	"github.com/Xe/olin/abi"
)

// Random returns a file that reads cryptographically random data.
func Random() abi.File {
	return Reader(rand.Reader, "random")
}
