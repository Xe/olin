package fileresolver

import (
	"testing"

	"github.com/Xe/olin/internal/abi"
)

func TestAferoFileIsABIFile(t *testing.T) {
	var _ abi.File = AferoFile{}
}
