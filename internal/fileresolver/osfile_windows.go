// +build windows

package fileresolver

import "github.com/Xe/olin/internal/abi"

// OSFile is a kludge for Windows
type OSFile struct{}

func NewOSFile(_ uintptr, _ string) abi.File {
	return Null()
}
