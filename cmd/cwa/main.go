package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Xe/olin/internal/abi/cwa"
	"github.com/perlin-network/life/exec"
)

var (
	jitEnabled = flag.Bool("jit-enabled", false, "enable jit?")
	gas        = flag.Int("gas", 65536*64, "number of instructions the VM can perform")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "  %s <file.wasm>", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	argv := flag.Args()
	if len(argv) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	fname := argv[0]
	argv = argv[1:]

	data, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}

	p := cwa.NewProcess(fname, argv, map[string]string{})

	cfg := exec.VMConfig{
		EnableJIT: *jitEnabled,
	}
	vm, err := exec.NewVirtualMachine(data, cfg, p)
	if err != nil {
		log.Fatalf("%s: %v", fname, err)
	}

	main, ok := vm.GetFunctionExport("main")
	if !ok {
		log.Fatalf("%s: no main function exported", fname)
	}

	ret, err := vm.RunWithGasLimit(main, *gas)
	if err != nil {
		log.Fatalf("%s: vm error: %v", fname, err)
	}

	if ret != 0 {
		log.Fatalf("%s: exit status %d", fname, ret)
	}

	os.Exit(int(ret))
}
