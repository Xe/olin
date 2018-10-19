package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Xe/olin/rpc/archway"
)

func CreateEvent(fl *flag.FlagSet) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tp := fl.Lookup("topic")
	archwayURL := fl.Lookup("url")
	mimeType := fl.Lookup("mime")
	cli := archway.NewInteropProtobufClient(archwayURL.Value.String(), http.DefaultClient)

	fname := fl.Arg(1)
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("can't read event input: %v", err)
	}

	ev := &archway.Event{
		Topic:    tp.Value.String(),
		Data:     data,
		MimeType: mimeType.Value.String(),
	}

	_, err = cli.CreateEvent(ctx, ev)
	if err != nil {
		log.Fatalf("can't create event: %v", err)
	}
}
