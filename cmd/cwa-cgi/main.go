package main

import (
	"context"
	"expvar"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Xe/olin/internal/cwagi"
	"within.website/ln"
	"within.website/ln/opname"
)

var (
	mainFunc    = flag.String("main-func", "main", "main function to call (because rust takes over the name main)")
	addr        = flag.String("addr", ":8400", "TCP host:port to listen on")
	poolSize    = flag.Int("pool-size", 1, "initial worker pool size")
	maxPoolSize = flag.Int("max-pool-size", 32, "maximum worker pool size")
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

	vp := cwagi.NewPool(data, fname, *mainFunc, *poolSize, *maxPoolSize)
	defer vp.Close()

	mux := http.NewServeMux()
	mux.Handle("/", vp)
	mux.Handle("/expvar", expvar.Handler())
	mux.HandleFunc("/reboot", func(w http.ResponseWriter, r *http.Request) {
		os.Exit(0)
	})

	ln.Log(opname.With(ctx, "listenAndServe"), f)
	ln.FatalErr(ctx, http.ListenAndServe(*addr, mux), f)
}
