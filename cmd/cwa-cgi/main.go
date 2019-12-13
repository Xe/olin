package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/http/cgi"
	"os"
	"os/exec"

	"github.com/google/uuid"
	"github.com/povilasv/prommod"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"within.website/ln"
	"within.website/ln/opname"
)

var (
	mainFunc = flag.String("main-func", "_start", "main function to call (because rust takes over the name main)")
	addr     = flag.String("addr", ":8400", "TCP host:port to listen on")
	bin      = flag.String("bin", "cwagi.wasm", "cgi handler to run")
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
		"main_func": *mainFunc,
		"addr":      *addr,
	}

	_ = prometheus.Register(prommod.NewCollector("cwa-cgi"))

	cwaPath, err := exec.LookPath("cwa")
	if err != nil {
		ln.FatalErr(ctx, err)
	}

	h := &cgi.Handler{
		Path: cwaPath,
		Env:  []string{"RUN_ID=" + uuid.New().String(), "WORKER_ID=" + uuid.New().String()},
		Args: []string{*bin},
		Dir:  "/tmp",
	}

	mux := http.NewServeMux()
	mux.Handle("/", h)
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/reboot", func(w http.ResponseWriter, r *http.Request) {
		os.Exit(0)
	})

	ln.Log(opname.With(ctx, "listenAndServe"), f)
	ln.FatalErr(ctx, http.ListenAndServe(*addr, mux), f)
}
