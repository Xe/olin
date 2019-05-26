package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/Xe/olin/abi/dagger"
	"github.com/perlin-network/life/compiler"
	"github.com/perlin-network/life/exec"
	"github.com/spf13/afero"
)

var (
	mainFunc = flag.String("main-func", "dagger_main", "main function to call (because rust is broken)")
	vmStats  = flag.Bool("vm-stats", false, "dump VM statistics?")
	gas      = flag.Int64("gas", 65536*64, "number of instructions the VM can perform")
	writeMem = flag.String("write-mem", "", "write memory heap to the given file on exit")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "  %s <file.wasm>\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	rand.Seed(time.Now().UnixNano())
}

func main() {
	flag.Parse()

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

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	p := dagger.NewProcess(fname, afero.NewBasePathFs(afero.NewOsFs(), cwd))
	p.Stdin = os.Stdin
	p.Stdout = os.Stdout
	p.Stderr = os.Stderr

	cfg := exec.VMConfig{
		GasLimit: uint64(*gas),
	}

	gp := &compiler.SimpleGasPolicy{GasPerInstruction: 1}
	vm, err := exec.NewVirtualMachine(data, cfg, p, gp)
	vmInitTime := time.Since(t0)
	if err != nil {
		log.Fatalf("%s: %v", fname, err)
	}

	main, ok := vm.GetFunctionExport(*mainFunc)
	if !ok {
		log.Fatalf("%s: no main function exported", fname)
	}

	t0 = time.Now()
	ret, err := vm.RunWithGasLimit(main, int(*gas))
	if err != nil {
		log.Fatalf("runtime error for VM %s: %v", fname, err)
	}
	vmRunTime := time.Since(t0)

	if err != nil {
		log.Fatalf("%s: vm error: %v", fname, err)
	}
	if *vmStats {
		log.Printf("reading file time: %s", readingFileTime)
		log.Printf("vm init time:      %s", vmInitTime)
		log.Printf("vm gas limit:      %v", *gas)
		log.Printf("vm gas used:       %v", vm.Gas)
		log.Printf("vm gas percentage: %v", float64(float64(vm.Gas)/float64(*gas))*100)
		log.Printf("vm syscalls:       %d", p.SyscallCount())
		log.Printf("execution time:    %s", vmRunTime)
		log.Printf("exit status:       %d", ret)
		log.Printf("memory pages:      %d", len(vm.Memory)/65536)
	}

	if fname := *writeMem; fname != "" {
		log.Printf("writing memory to %s (%d bytes)", fname, len(vm.Memory))
		err := ioutil.WriteFile(fname, vm.Memory, 0600)
		if err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(int(ret))
}
