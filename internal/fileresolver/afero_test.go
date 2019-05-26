package fileresolver

import (
	"testing"

	"github.com/Xe/olin/abi"
)

func TestAferoFileIsABIFile(t *testing.T) {
	var _ abi.File = AferoFile{}
}
