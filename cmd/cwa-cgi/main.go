package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Xe/ln"
	"github.com/Xe/ln/opname"
	"github.com/Xe/olin/internal/abi/cwa"
	"github.com/pborman/uuid"
	"github.com/perlin-network/life/exec"
)

var (
	minPages   = flag.Int("min-pages", 32, "number of memory pages to open (default is 2 MB)")
	mainFunc   = flag.String("main-func", "main", "main function to call (because rust takes over the name main)")
	jitEnabled = flag.Bool("jit-enabled", false, "enable jit?")
	gas        = flag.Int("gas", 65536*64, "number of instructions the VM can perform per handler invocation")
	addr       = flag.String("addr", ":8400", "TCP host:port to listen on")
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
	ctx := opname.With(context.Background(), "main")
	f := ln.F{
		"min_pages":   *minPages,
		"main_func":   *mainFunc,
		"jit_enabled": *jitEnabled,
		"addr":        *addr,
	}

	argv := flag.Args()
	if len(argv) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	fname := argv[0]

	data, err := ioutil.ReadFile(fname)
	if err != nil {
		ln.FatalErr(ctx, err, f)
	}

	p := cwa.NewProcess(fname, argv, map[string]string{})

	cfg := exec.VMConfig{
		EnableJIT:          *jitEnabled,
		DefaultMemoryPages: *minPages,
	}
	vm, err := exec.NewVirtualMachine(data, cfg, p)
	if err != nil {
		ln.FatalErr(ctx, err, f)
	}

	log.Printf("loading function %s", *mainFunc)
	main, ok := vm.GetFunctionExport(*mainFunc)
	if !ok {
		ln.FatalErr(ctx, err, f)
	}

	vs := &vmServer{
		vm:       vm,
		p:        p,
		mainFunc: main,
	}

	mux := http.NewServeMux()
	mux.Handle("/", vs)

	ln.Log(opname.With(ctx, "listenAndServe"), f)
	ln.FatalErr(ctx, http.ListenAndServe(*addr, mux), f)
}

type vmServer struct {
	vm       *exec.VirtualMachine
	p        *cwa.Process
	lock     sync.Mutex
	mainFunc int
}

func (v *vmServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stdout := bytes.NewBuffer(nil)
	stdin := r.Body
	defer r.Body.Close()
	ctx := opname.With(r.Context(), "vmServer.ServeHTTP")
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	runID := uuid.New()

	f := ln.F{
		"gas_limit":    *gas,
		"main_func":    v.mainFunc,
		"process_name": v.p.Name(),
		"run_id":       runID,
		"method":       r.Method,
		"request_uri":  r.RequestURI,
	}

	v.lock.Lock()
	defer v.lock.Unlock()
	v.p.Stdin = stdin
	v.p.Stdout = stdout
	v.p.Setenv(map[string]string{
		"REQUEST_METHOD": r.Method,
		"REQUEST_URI":    r.RequestURI,
		"QUERY_STRING":   r.URL.Query().Encode(),
		"RUN_ID":         runID,
	})

	t0 := time.Now()
	ret, err := v.vm.RunWithGasLimit(v.mainFunc, *gas)
	if err != nil {
		http.Error(w, "internal server error: VM error, run ID: "+runID, http.StatusInternalServerError)
		go func() {
			time.Sleep(125 * time.Millisecond)
			ln.FatalErr(ctx, err, f)
		}()
		return
	}
	f["exec_dur"] = time.Since(t0)

	if ret != 0 {
		ln.Log(ctx, f, ln.F{
			"return_value": ret,
		})
		http.Error(w, fmt.Sprintf("internal server error: return code %d", ret), http.StatusInternalServerError)
		return
	}

	ctx = opname.With(ctx, "respond")
	resp, err := http.ReadResponse(bufio.NewReader(stdout), r)
	if err != nil {
		ln.Error(ctx, err, f)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// copy headers
	for k, v := range resp.Header {
		for _, val := range v {
			w.Header().Add(k, val)
		}
	}

	// copy status code
	w.WriteHeader(resp.StatusCode)
	f["status"] = resp.StatusCode

	// copy body
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		ln.Error(opname.With(ctx, "copy_body"), err, f)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	ln.Log(ctx, f, ln.Info("successful invocation"))
}
