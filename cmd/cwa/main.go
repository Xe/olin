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
	"github.com/Xe/olin/internal/abi/wasmgo"
	"github.com/perlin-network/life/compiler"
	"github.com/perlin-network/life/exec"
)

var (
	minPages   = flag.Int("min-pages", 32, "number of memory pages to open (default is 2 MB)")
	mainFunc   = flag.String("main-func", "cwa_main", "main function to call (because rust is broken)")
	jitEnabled = flag.Bool("jit-enabled", false, "enable jit?")
	doTest     = flag.Bool("test", false, "unit testing?")
	vmStats    = flag.Bool("vm-stats", false, "dump VM statistics?")
	gas        = flag.Int("gas", 65536*64, "number of instructions the VM can perform")
	goMode     = flag.Bool("go", false, "run in Go mode?")
	writeMem   = flag.String("write-mem", "", "write memory heap to the given file on exit")
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

	t0 := time.Now()
	data, err := ioutil.ReadFile(fname)
	readingFileTime := time.Since(t0)
	if err != nil {
		log.Fatal(err)
	}

	if *doTest {
		environment = map[string]string{
			"MAGIC_CONCH": "yes",
		}
	}

	p := wasmgo.New(fname, argv, environment)

	if *doTest {
		p.Stdin = bytes.NewBuffer([]byte("cwa test environment"))
	}

	cfg := exec.VMConfig{
		EnableJIT:          *jitEnabled,
		DefaultMemoryPages: *minPages,
	}

	gp := &compiler.SimpleGasPolicy{GasPerInstruction: 1}
	t0 = time.Now()
	vm, err := exec.NewVirtualMachine(data, cfg, p, gp)
	vmInitTime := time.Since(t0)
	if err != nil {
		log.Fatalf("%s: %v", fname, err)
	}

	var ret int32
	if *goMode {
		ret, err = runGo(p, vm)
	} else {
		ret, err = runCWA(p.Process, vm)
	}
	vmRunTime := time.Since(t0)
	if err != nil {
		log.Fatalf("%s: vm error: %v", fname, err)
	}
	if *vmStats || *doTest {
		log.Printf("reading file time: %s", readingFileTime)
		log.Printf("vm init time:      %s", vmInitTime)
		log.Printf("vm gas limit:      %v", *gas)
		log.Printf("vm gas used:       %v", vm.Gas)
		log.Printf("vm gas percentage: %v", float64(float64(vm.Gas)/float64(*gas))*100)
		log.Printf("vm syscalls:       %d", p.SyscallCount())
		log.Printf("execution time:    %s", vmRunTime)
	}

	if ret != 0 {
		log.Fatalf("%s: exit status %d", fname, ret)
	}

	if *vmStats {
		log.Printf("memory pages:      %d", len(vm.Memory)/65536)
	}

	if fname := *writeMem; fname != "" {
		err := ioutil.WriteFile(fname, vm.Memory, 0600)
		if err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(int(ret))
}

func runCWA(p *cwa.Process, vm *exec.VirtualMachine) (int32, error) {
	if *doTest {
		log.Printf("loading function %s", *mainFunc)
	}

	main, ok := vm.GetFunctionExport(*mainFunc)
	if !ok {
		return -1, fmt.Errorf("%s: no main function exported", p.Name())
	}

	if *doTest {
		log.Printf("executing %s (%d)", *mainFunc, main)
	}

	ret, err := vm.RunWithGasLimit(main, *gas)
	if err != nil {
		return 1, err
	}

	return int32(ret), nil
}

func runGo(w *wasmgo.WasmGo, vm *exec.VirtualMachine) (int32, error) {
	log.Printf("starting wasmgo...")
	w.Memory.Data = vm.Memory

	run, ok := vm.GetFunctionExport("run")
	if !ok {
		panic("function not found: run")
	}

	resume, ok := vm.GetFunctionExport("resume")
	if !ok {
		panic("function not found: resume")
	}

	if _, err := vm.RunWithGasLimit(run, *gas, 0, 0); err != nil {
		return 1, err
	}

	for !vm.Exited {
		if _, err := vm.RunWithGasLimit(resume, *gas); err != nil {
			return w.StatusCode, err
		}
	}

	return w.StatusCode, nil
}
