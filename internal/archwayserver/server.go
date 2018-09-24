package archwayserver

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/Xe/ln"
	"github.com/Xe/olin/rpc/archway"
	pubsub "github.com/alash3al/go-pubsub"
	"github.com/boltdb/bolt"
	"github.com/djherbis/stow"
	"github.com/pborman/uuid"
)

type Interop struct {
	hd *stow.Store
	ev *stow.Store
	b  *pubsub.Broker
}

func New(db *bolt.DB, b *pubsub.Broker) Interop {
	hd := stow.NewJSONStore(db, []byte("handlers"))
	ev := stow.NewJSONStore(db, []byte("events"))

	return Interop{
		hd: hd,
		ev: ev,
		b:  b,
	}
}

func (i Interop) CreateHandler(ctx context.Context, hdl *archway.Handler) (*archway.Handler, error) {
	if hdl.GetId() != "" {
		return nil, errors.New("can't create a handler with an ID")
	}

	id := uuid.New()
	hdl.Id = id
	hdl.CreatedAtUnixUtc = time.Now().UTC().Unix()
	u := url.URL{
		Scheme: "handler",
		Host:   id,
	}

	var err error

	err = i.hd.Put(u.String(), hdl)
	if err != nil {
		return nil, err
	}

	i.b.Broadcast(hdl, "archway.internal.handler_create")

	ln.Log(ctx, ln.F{"topic": hdl.Topic}, ln.Info("handler created"))

	return hdl, nil
}

func (i Interop) DeleteHandler(ctx context.Context, id *archway.Id) (*archway.Handler, error) {
	u := url.URL{
		Scheme: "handler",
		Host:   id.GetId(),
	}

	var hdl archway.Handler
	err := i.hd.Get(u.String(), &hdl)
	if err != nil {
		return nil, err
	}

	err = i.hd.Delete(u.String())
	if err != nil {
		return nil, err
	}

	i.b.Broadcast(hdl, "archway.internal.handler_destroy")

	return &hdl, nil
}

func (i Interop) GetHandler(ctx context.Context, id *archway.Id) (*archway.Handler, error) {
	u := url.URL{
		Scheme: "handler",
		Host:   id.GetId(),
	}

	var hdl archway.Handler
	err := i.hd.Get(u.String(), &hdl)
	if err != nil {
		return nil, err
	}

	return &hdl, nil
}

func (i Interop) ListHandlers(ctx context.Context, _ *archway.Nil) (*archway.Handlers, error) {
	result := &archway.Handlers{}

	i.hd.ForEach(func(h *archway.Handler) {
		result.Handlers = append(result.Handlers, h)
	})

	return result, nil
}

func (i Interop) CreateEvent(ctx context.Context, e *archway.Event) (*archway.Nil, error) {
	id := uuid.New()

	u := url.URL{
		Scheme: "event",
		Host:   id,
	}

	e.Id = id
	e.CreatedAtUnixUtc = time.Now().UTC().Unix()

	err := i.ev.Put(u.String(), e)
	if err != nil {
		return nil, err
	}

	i.b.Broadcast(e, e.GetTopic())

	ln.Log(ctx, ln.F{"topic": e.GetTopic()}, ln.Action("EventCreated"))

	return &archway.Nil{}, nil
}

func (i Interop) GetEvent(ctx context.Context, id *archway.Id) (*archway.Event, error) {
	u := url.URL{
		Scheme: "event",
		Host:   id.GetId(),
	}

	var e archway.Event
	err := i.ev.Get(u.String(), &e)
	if err != nil {
		return nil, err
	}

	return &e, nil
}
