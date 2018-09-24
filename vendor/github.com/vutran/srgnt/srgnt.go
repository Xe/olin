package srgnt

import (
	"os"
)

func CreateProgram(name string) Program {
	args := os.Args
	return Program{Name: name, Args: args}
}
