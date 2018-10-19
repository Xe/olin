package main

import (
	"github.com/vutran/srgnt"
)

func main() {
	cli := srgnt.CreateProgram("olin")

	cli.AddStringFlag("url", "http://127.0.0.1:1324", "archwayd URL")
	cli.AddStringFlag("topic", "", "topic for this command to relate to")
	cli.AddStringFlag("mime", "", "mime-type for the event payload")
	cli.AddStringFlag("env", "", "url-encoded environment for the process (sorry)")

	cli.AddCommand("event_create", CreateEvent, "creates a new event")
	cli.AddCommand("event_recent", RecentEvent, "gets the most recent event for a topic")
	cli.AddCommand("handler_create", CreateHandler, "creates a WebAssembly handler")

	cli.Run()
}
