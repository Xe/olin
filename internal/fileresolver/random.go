package fileresolver

import (
	"crypto/rand"

	"within.website/olin/internal/abi"
)

// Random returns a file that reads cryptographically random data.
func Random() abi.File {
	return Reader(rand.Reader, "random")
}
