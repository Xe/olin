package main

import (
	"context"
	"expvar"
	"log"
	"net/http"

	"github.com/Xe/ln"
	"github.com/Xe/ln/opname"
	"github.com/Xe/olin/internal/archwayserver"
	"github.com/Xe/olin/rpc/archway"
	pubsub "github.com/alash3al/go-pubsub"
	"github.com/boltdb/bolt"
	"github.com/go-kit/kit/metrics/provider"
	"github.com/joeshaw/envdecode"
)

type config struct {
	Port   string `env:"PORT,default=1324"`
	APIKey string `env:"API_KEY,default=hunter2"`
	DBPath string `env:"DB_PATH,default=./var/archwayd.db"`
}

func (c config) F() ln.F {
	return ln.F{
		"port":    c.Port,
		"api_key": c.APIKey,
		"db_path": c.DBPath,
	}
}

func main() {
	ctx := context.Background()
	ctx = opname.With(ctx, "main")
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var cfg config
	err := envdecode.StrictDecode(&cfg)
	if err != nil {
		ln.FatalErr(ctx, err)
	}

	b := pubsub.NewBroker()

	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	iop := archwayserver.New(db, b)
	rnr := archwayserver.NewRunner(b)
	go rnr.Manage(ctx)

	mux := http.NewServeMux()
	mux.Handle(archway.InteropPathPrefix, archway.NewInteropServer(
		archway.NewInteropLogging(
			archway.NewInteropMetrics(iop, provider.NewExpvarProvider()),
		), nil,
	))
	mux.Handle("/", expvar.Handler())

	ln.Log(ctx, cfg, ln.Info("Listening on HTTP"))
	ln.FatalErr(ctx, http.ListenAndServe(":"+cfg.Port, mux))
}
