package dagger

import (
	"sync"

	"github.com/Xe/olin/internal/abi"
)

type Process struct {
	sync.Mutex
	abi.ABI

	files []abi.File
}
