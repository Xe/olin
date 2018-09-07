package fileresolver

import (
	"fmt"
	"testing"
)

func TestZero(t *testing.T) {
	assert := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	ze := Zero()

	val := make([]byte, 512)
	n, err := ze.Write(val)
	assert(err)
	if n != len(val) {
		assert(fmt.Errorf("wanted n to be len(val) (%d), got: %d", len(val), n))
	}

	n, err = ze.Read(val)
	assert(err)
	if n != len(val) {
		assert(fmt.Errorf("wanted n to be len(val) (%d), got: %d", len(val), n))
	}

	assert(ze.Sync())
	assert(ze.Close())
}
