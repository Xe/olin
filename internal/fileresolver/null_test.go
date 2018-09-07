package fileresolver

import (
	"fmt"
	"testing"
)

func TestNull(t *testing.T) {
	assert := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	nu := Null()

	val := make([]byte, 512)
	n, err := nu.Write(val)
	assert(err)
	if n != len(val) {
		assert(fmt.Errorf("wanted n to be len(val) (%d), got: %d", len(val), n))
	}

	n, err = nu.Read(val)
	assert(err)
	if n != 0 {
		assert(fmt.Errorf("wanted n to be 0, got: %d", n))
	}

	assert(nu.Sync())
	assert(nu.Close())
}
