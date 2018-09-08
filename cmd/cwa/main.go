package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/Xe/olin/internal/abi/cwa"
	"github.com/perlin-network/life/exec"
)

var (
	minPages   = flag.Int("min-pages", 32, "number of memory pages to open (default is 2 MB)")
	mainFunc   = flag.String("main-func", "main", "main function to call (because rust is broken)")
	jitEnabled = flag.Bool("jit-enabled", false, "enable jit?")
	gas        = flag.Int("gas", 65536*64, "number of instructions the VM can perform")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "  %s <file.wasm>\n\n", os.Args[0])
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

	p := cwa.NewProcess(fname, argv, map[string]string{
		"MAGIC_CONCH": "yes",
	})

	cfg := exec.VMConfig{
		EnableJIT:          *jitEnabled,
		DefaultMemoryPages: *minPages,
	}
	vm, err := exec.NewVirtualMachine(data, cfg, p)
	if err != nil {
		log.Fatalf("%s: %v", fname, err)
	}

	log.Printf("loading function %s", *mainFunc)
	main, ok := vm.GetFunctionExport(*mainFunc)
	if !ok {
		log.Fatalf("%s: no main function exported", fname)
	}

	log.Printf("executing %s (%d)", *mainFunc, main)
	t0 := time.Now()
	ret, err := vm.RunWithGasLimit(main, *gas)
	if err != nil {
		log.Fatalf("%s: vm error: %v", fname, err)
	}
	log.Printf("execution time: %s", time.Since(t0))

	if ret != 0 {
		log.Fatalf("%s: exit status %d", fname, ret)
	}

	log.Printf("memory pages: %d", len(vm.Memory)/65536)

	os.Exit(int(ret))
}
