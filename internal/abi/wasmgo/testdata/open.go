// +build ignore

package main

import "os"

func main() {
	fout, err := os.Open("log://")
	if err != nil {
		panic(err)
	}

	_, err = fout.Write([]byte("memes"))
	if err != nil {
		panic(err)
	}

	err = fout.Close()
	if err != nil {
		panic(err)
	}
}
