package fileresolver

import (
	"bytes"
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	assert := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}

	buf := bytes.NewBuffer(nil)
	lf := Log(buf, "fun", 0)

	val := []byte("looking forward to the weekend")
	n, err := lf.Write(val)
	if n != len(val) {
		assert(fmt.Errorf("wanted n to be len(val) (%d), got: %d", len(val), n))
	}
	assert(err)

	if buf.String() != "fun: looking forward to the weekend\n" {
		t.Logf("buf: %q", buf)
		t.Fatalf("wanted buf to match expected message")
	}

	_, err = lf.Read(val)
	assert(err)
	assert(lf.Flush())
	assert(lf.Close())
}
