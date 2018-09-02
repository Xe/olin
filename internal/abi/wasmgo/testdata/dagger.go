// +build js,wasm ignore

package main

import (
	"fmt"

	"github.com/Xe/olin/dagger"
)

func main() {
	fd := dagger.OpenFile("fd://1", 0)
	println(fmt.Sprint(fd))
}
