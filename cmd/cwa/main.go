package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Xe/olin/internal/abi/cwa"
	"github.com/perlin-network/life/exec"
	"github.com/perlin-network/life/compiler"
)

var (
	minPages   = flag.Int("min-pages", 32, "number of memory pages to open (default is 2 MB)")
	mainFunc   = flag.String("main-func", "cwa_main", "main function to call (because rust is broken)")
	jitEnabled = flag.Bool("jit-enabled", false, "enable jit?")
	doTest     = flag.Bool("test", false, "unit testing?")
	vmStats    = flag.Bool("vm-stats", false, "dump VM statistics?")
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

	getenvironment := func(data []string, getkeyval func(item string) (key, val string)) map[string]string {
		items := make(map[string]string)
		for _, item := range data {
			key, val := getkeyval(item)
			items[key] = val
		}
		return items
	}
	environment := getenvironment(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = splits[1]
		return
	})

	argv := flag.Args()
	if len(argv) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	fname := argv[0]

	data, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}

	if *doTest {
		environment = map[string]string{
			"MAGIC_CONCH": "yes",
		}
	}

	p := cwa.NewProcess(fname, argv, environment)

	if *doTest {
		p.Stdin = bytes.NewBuffer([]byte("cwa test environment"))
	}

	cfg := exec.VMConfig{
		EnableJIT:          *jitEnabled,
		DefaultMemoryPages: *minPages,
	}

	gp := &compiler.SimpleGasPolicy{GasPerInstruction: 1}
	vm, err := exec.NewVirtualMachine(data, cfg, p, gp)
	if err != nil {
		log.Fatalf("%s: %v", fname, err)
	}

	if *doTest {
		log.Printf("loading function %s", *mainFunc)
	}

	main, ok := vm.GetFunctionExport(*mainFunc)
	if !ok {
		log.Fatalf("%s: no main function exported", fname)
	}

	if *doTest {
		log.Printf("executing %s (%d)", *mainFunc, main)
	}

	t0 := time.Now()
	ret, err := vm.RunWithGasLimit(main, *gas)
	if err != nil {
		log.Fatalf("%s: vm error: %v", fname, err)
	}
	if *vmStats || *doTest {
		log.Printf("execution time: %s", time.Since(t0))
	}

	if ret != 0 {
		log.Fatalf("%s: exit status %d", fname, ret)
	}

	if *vmStats {
		log.Printf("memory pages: %d", len(vm.Memory)/65536)
	}

	os.Exit(int(ret))
}
