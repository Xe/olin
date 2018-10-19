package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Xe/olin/rpc/archway"
)

func RecentEvent(fl *flag.FlagSet) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tp := fl.Lookup("topic")
	archwayURL := fl.Lookup("url")
	cli := archway.NewInteropProtobufClient(archwayURL.Value.String(), http.DefaultClient)

	ev, err := cli.GetMostRecentEvent(ctx, &archway.Topic{Topic: tp.Value.String()})
	if err != nil {
		log.Fatalf("can't get most recent event: %v", err)
	}

	fmt.Printf("ID: %s\n", ev.Id)
	fmt.Printf("Time: %s (%d)\n", time.Unix(int64(ev.CreatedAtUnixUtc), 0).UTC().Format(time.RFC3339), ev.CreatedAtUnixUtc)
	fmt.Printf("Topic: %s\n", ev.Topic)
	fmt.Printf("Mime type: %s\n", ev.MimeType)

	switch ev.MimeType {
	case "application/octet-stream":
	default:
		fmt.Printf("Body: \n%s\n", string(ev.Data))
	}
}
