package main

import (
	"context"
	"os"
	"os/exec"
	"testing"
)

func TestCWASpec(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := exec.CommandContext(ctx, "/bin/sh", "./test.sh")
	p.Stdout = os.Stdout
	p.Stderr = os.Stderr
	err := p.Start()
	if err != nil {
		t.Fatal(err)
	}
	err = p.Wait()
	if err != nil {
		t.Fatal(err)
	}
}
