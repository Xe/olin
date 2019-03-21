package cwagi

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/Xe/olin/internal/abi/cwa"
	"github.com/pborman/uuid"
	"github.com/perlin-network/life/compiler"
	"github.com/perlin-network/life/exec"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"within.website/ln"
	"within.website/ln/opname"
)

var (
	ramUse = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "ram_use",
		Help: "The amount of ram in use per VM",
	}, []string{"vm"})

	executionTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "execution_time",
		Help: "The amount of time spent executing for the VM",
	}, []string{"vm", "path"})

	gasUsed = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "gas_used",
		Help: "The amount of VM gas used",
	}, []string{"vm", "path"})
)

// NewVM creates a new virtual machine with the given WebAssembly binary code and name.
func NewVM(data []byte, argv []string, name, mainFunc string) (*VMServer, error) {
	myID := uuid.New()

	p := cwa.NewProcess(name+"+"+myID, argv, map[string]string{})

	cfg := exec.VMConfig{
		EnableJIT:          false,
		DefaultMemoryPages: 32, // 2 MB
	}
	gp := compiler.SimpleGasPolicy{GasPerInstruction: 1}
	vm, err := exec.NewVirtualMachine(data, cfg, p, &gp)
	if err != nil {
		return nil, err
	}

	main, ok := vm.GetFunctionExport(mainFunc)
	if !ok {
		return nil, errors.New("cwagi: need main function to be exported")
	}

	return &VMServer{
		VM:       vm,
		P:        p,
		mainFunc: main,
		myID:     myID,
	}, nil
}

// VMServer wraps a single WebAssembly virtual machine around a HTTP style handler.
type VMServer struct {
	VM       *exec.VirtualMachine
	P        *cwa.Process
	lock     sync.Mutex
	mainFunc int
	myID     string
}

func (v *VMServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	stdout := bytes.NewBuffer(nil)
	stdin := r.Body
	defer r.Body.Close()
	ctx := opname.With(r.Context(), "vmServer.ServeHTTP")
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	runID := uuid.New()

	f := ln.F{
		"main_func":    v.mainFunc,
		"process_name": v.P.Name(),
		"run_id":       runID,
		"method":       r.Method,
		"request_uri":  r.RequestURI,
	}

	v.lock.Lock()
	defer v.lock.Unlock()
	v.P.Stdin = stdin
	v.P.Stdout = stdout
	v.P.Setenv(map[string]string{
		"REQUEST_METHOD": r.Method,
		"REQUEST_URI":    r.RequestURI,
		"QUERY_STRING":   r.URL.Query().Encode(),
		"RUN_ID":         runID,
		"WORKER_ID":      v.myID,
	})

	t0 := time.Now()
	ret, err := v.VM.Run(v.mainFunc)
	if err != nil {
		http.Error(w, "internal server error: VM error, run ID: "+runID, http.StatusInternalServerError)
		go func() {
			time.Sleep(125 * time.Millisecond)
			ln.FatalErr(ctx, err, f)
		}()
		return
	}
	dur := time.Since(t0)
	f["exec_dur"] = dur
	executionTime.With(prometheus.Labels{"vm": v.myID, "path": r.URL.Path}).Observe(float64(dur))
	gasUsed.With(prometheus.Labels{"vm": v.myID, "path": r.URL.Path}).Observe(float64(v.VM.Gas))

	exitStatus.With(prometheus.Labels{"status": fmt.Sprint(ret)}).Inc()

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

	if rand.Int63()%50 == 0 {
		ln.Log(ctx, f, ln.Info("successful invocation"))
	}
}
