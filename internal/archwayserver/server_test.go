package archwayserver

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Xe/olin/rpc/archway"
	pubsub "github.com/alash3al/go-pubsub"
	"github.com/boltdb/bolt"
)

func newInterop(t *testing.T) (*Interop, func()) {
	t.Helper()

	file, err := ioutil.TempFile("", "archwayserver_test")
	if err != nil {
		t.Fatal(err)
	}

	db, err := bolt.Open(file.Name(), 0600, nil)
	if err != nil {
		t.Fatal(err)
	}
	b := pubsub.NewBroker()

	i := New(db, b)

	cleanup := func() {
		db.Close()
		os.Remove(file.Name())
	}

	return &i, cleanup
}

func TestSendEvent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	i, cup := newInterop(t)
	defer cup()

	data, err := ioutil.ReadFile("./log-env-message.wasm")
	if err != nil {
		t.Fatal(err)
	}

	hdl := &archway.Handler{
		Topic:  "test",
		Module: data,
		Env:    map[string]string{"MESSAGE": "Hi"},
	}

	hdl, err = i.CreateHandler(ctx, hdl)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("handler ID: %s", hdl.Id)

	ev := &archway.Event{
		Topic:    "test",
		Data:     []byte("{}"),
		MimeType: "application/json",
	}

	_, err = i.CreateEvent(ctx, ev)
	if err != nil {
		t.Fatal(err)
	}
}
