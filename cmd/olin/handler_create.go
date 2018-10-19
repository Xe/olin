package main

import (
	"bytes"
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/Xe/olin/rpc/archway"
	"github.com/Xe/olin/rpc/brand"
	"github.com/go-interpreter/wagon/wasm"
	"github.com/golang/protobuf/proto"
)

func CreateHandler(fl *flag.FlagSet) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tp := fl.Lookup("topic")
	archwayURL := fl.Lookup("url")
	env := fl.Lookup("env")
	cli := archway.NewInteropProtobufClient(archwayURL.Value.String(), http.DefaultClient)

	log.Printf("%v", fl.Args())
	if fl.NArg() != 2 || tp.Value.String() == "" {
		log.Fatal("usage: olin <-topic topic.name> handler_create <file.wasm>")
	}

	envValues, err := url.ParseQuery(env.Value.String())
	if err != nil {
		log.Fatalf("can't parse environment: %s: %v", env.Value.String(), err)
	}

	fname := fl.Arg(1)
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("can't read wasm file: %v", err)
	}

	mod, err := wasm.DecodeModule(bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("can't decode wasm file: %v", err)
	}

	if sc := mod.Custom("olin-settings"); sc != nil {
		var b brand.Brand
		err := proto.Unmarshal(sc.Data, &b)
		if err != nil {
			log.Fatalf("can't unmarshal settings: %v", err)
		}

		log.Println("custom settings:")
		log.Printf("JIT enabled: %v", b.Opts.EnableJit)
		log.Printf("Default ram pages: %d", b.Opts.DefaultPages)
		log.Printf("Max ram pages: %d", b.Opts.MaxPages)
		log.Printf("Main func: %s", b.Opts.MainFunc)
		log.Printf("Expected runtime: %s", b.Meta.ExpectedRuntime)
		log.Printf("Author: %s", b.Meta.Author)
		log.Printf("Name: %s", b.Meta.Name)
	} else {
		log.Printf("Binary has no custom settings defined. This may result in undesirable performance.")
	}

	hdl := &archway.Handler{
		Topic:  tp.Value.String(),
		Module: data,
		Env:    map[string]string{},
	}

	for k, v := range envValues {
		hdl.Env[k] = v[0]
	}

	hdl, err = cli.CreateHandler(ctx, hdl)
	if err != nil {
		log.Fatalf("can't create handler: %v", err)
	}

	log.Printf("handler ID: %s", hdl.Id)
}
