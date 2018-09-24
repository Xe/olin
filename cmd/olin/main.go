package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Xe/olin/rpc/archway"
	"github.com/vutran/srgnt"
)

func CreateHandler(fl *flag.FlagSet) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tp := fl.Lookup("topic")
	archwayURL := fl.Lookup("url")
	cli := archway.NewInteropProtobufClient(archwayURL.Value.String(), http.DefaultClient)

	log.Printf("%v", fl.Args())
	if fl.NArg() != 2 || tp.Value.String() == "" {
		log.Fatal("usage: olin <-topic topic.name> handler_create <file.wasm>")
	}

	fname := fl.Arg(1)
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("can't read wasm file: %v", err)
	}

	hdl := &archway.Handler{
		Topic:  tp.Value.String(),
		Module: data,
	}

	hdl, err = cli.CreateHandler(ctx, hdl)
	if err != nil {
		log.Fatalf("can't create handler: %v", err)
	}
}

func CreateEvent(fl *flag.FlagSet) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tp := fl.Lookup("topic")
	archwayURL := fl.Lookup("url")
	cli := archway.NewInteropProtobufClient(archwayURL.Value.String(), http.DefaultClient)

	fname := fl.Arg(1)
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("can't read event input: %v", err)
	}

	ev := &archway.Event{
		Topic: tp.Value.String(),
		Data:  data,
	}

	_, err = cli.CreateEvent(ctx, ev)
	if err != nil {
		log.Fatalf("can't create event: %v", err)
	}
}

func main() {
	cli := srgnt.CreateProgram("olin")

	cli.AddStringFlag("url", "http://127.0.0.1:1324", "archwayd URL")
	cli.AddStringFlag("topic", "", "topic for this command to relate to")

	cli.AddCommand("event_create", CreateEvent, "creates a new event")
	cli.AddCommand("handler_create", CreateHandler, "creates a WebAssembly handler")

	cli.Run()
}
