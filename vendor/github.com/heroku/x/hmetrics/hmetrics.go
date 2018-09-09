/* Copyright (c) 2018 Salesforce
 * All rights reserved.
 * Licensed under the BSD 3-Clause license.
 * For full license text, see LICENSE.txt file in the repo root  or https://opensource.org/licenses/BSD-3-Clause
 */

/*
Package hmetrics is a self-contained client for Heroku Go runtime metrics.

Typical usage is through the `github.com/heroku/x/hmetrics/onload` package
imported like so:

  import _ "github.com/heroku/x/hmetrics/onload"

You can find more information about Heroku Go runtime metrics here:
https://devcenter.heroku.com/articles/language-runtime-metrics-go
*/
package hmetrics

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	interval = 20 * time.Second
)

var (
	// DefaultEndpoint to report metrics to Heroku. The "HEROKU_METRICS_URL" env
	// var is set by the Heroku runtime when the `runtime-heroku-metrics` labs
	// flag is enabled. For more info see:
	// https://devcenter.heroku.com/articles/language-runtime-metrics
	//
	// DefaultEndpoint must be changed before Report is called.
	DefaultEndpoint = os.Getenv("HEROKU_METRICS_URL")
)

// AlreadyStarted represents an Error condition of already being started.
type AlreadyStarted struct{}

func (as AlreadyStarted) Error() string {
	return "already started"
}

func (as AlreadyStarted) Fatal() bool {
	return false
}

// HerokuMetricsURLUnset represents the Error condition when the
// HEROKU_METRICS_URL environment variables is unset or an empty string.
type HerokuMetricsURLUnset struct{}

func (e HerokuMetricsURLUnset) Error() string {
	return "cannot report metrics because HEROKU_METRICS_URL is unset"
}

func (e HerokuMetricsURLUnset) Fatal() bool {
	return true
}

// ErrHandler funcations are used to provide custom processing/handling of
// errors encountered during the collection or reporting of metrics to Heroku.
type ErrHandler func(err error) error

var (
	mu      sync.Mutex
	started bool
)

// Report go metrics to the endpoint until the context is canceled.
//
// Only one call to the ErrHandler will happen at a time. Metrics can be dropped
// or delayed if a call to ErrHandler takes longer than the reprting interval.
// Processing of metrics continues if the ErrHandler returns nil, but aborts if
// the ErrHandler itself returns an error. It is safe to pass a nil ErrHandler.
//
// Report is safe for concurrent usage, but calling it again w/o canceling the
// context passed previously returns an AlreadyStarted error. This is to ensure
// that metrics aren't duplicated. Report can be called again to restart
// reporting after the context is canceled.
func Report(ctx context.Context, endpoint string, ef ErrHandler) error {
	if err := startable(endpoint); err != nil {
		return err
	}

	report(ctx, &http.Client{Timeout: 20 * time.Second}, endpoint, errorHandler(ef))

	mu.Lock()
	defer mu.Unlock()
	started = false
	return nil
}

// err if not startable, otherwise start
func startable(endpoint string) error {
	mu.Lock()
	defer mu.Unlock()
	if started {
		return AlreadyStarted{}
	}
	if endpoint == "" {
		return &url.Error{
			Op:  "Empty",
			URL: endpoint,
			Err: errors.New("Empty string"),
		}
	}
	if _, err := url.Parse(endpoint); err != nil {
		return err
	}
	started = true
	return nil
}

func noopErrorHandler(_ error) error {
	return nil
}

func errorHandler(ef ErrHandler) ErrHandler {
	if ef == nil {
		return noopErrorHandler
	}
	return ef
}

func report(ctx context.Context, client *http.Client, endpoint string, ef ErrHandler) {
	t := time.NewTicker(interval)
	defer t.Stop()

	var buf bytes.Buffer
	var pauseTotalNS uint64
	var numGC uint32
	for {
		select {
		case <-t.C:
		case <-ctx.Done():
			return
		}

		buf.Reset()

		var err error
		pauseTotalNS, numGC, err = gatherMetrics(&buf, pauseTotalNS, numGC)
		if err != nil {
			if err := ef(err); err != nil {
				return
			}
			continue
		}
		if err := submitMetrics(ctx, client, &buf, endpoint); err != nil {
			if err := ef(err); err != nil {
				return
			}
			continue
		}
	}
}

// gatherMetrics and write the JSON encoded representation to w.
// returns the sampled PauseTotalNs+NumGC & any encoding errors.
// TODO: If we ever have high frequency charts HeapIdle minus HeapReleased could be interesting.
func gatherMetrics(w io.Writer, prevPauseTotalNS uint64, prevNumGC uint32) (uint64, uint32, error) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	// cribbed from https://github.com/codahale/metrics/blob/master/runtime/memstats.go
	result := struct {
		Counters map[string]float64 `json:"counters"`
		Gauges   map[string]float64 `json:"gauges"`
	}{
		Counters: map[string]float64{
			"go.gc.collections": float64(stats.NumGC - prevNumGC),
			"go.gc.pause.ns":    float64(stats.PauseTotalNs - prevPauseTotalNS),
		},
		Gauges: map[string]float64{
			"go.memory.heap.bytes":   float64(stats.Alloc),
			"go.memory.stack.bytes":  float64(stats.StackInuse),
			"go.memory.heap.objects": float64(stats.Mallocs - stats.Frees), // Number of "live" objects.
			"go.gc.goal":             float64(stats.NextGC),                // Goal heap size for next GC.
			"go.routines":            float64(runtime.NumGoroutine()),      // Current number of goroutines.
		},
	}

	return stats.PauseTotalNs, stats.NumGC, json.NewEncoder(w).Encode(result)
}

// submitMetrics read from r to the endpoint using the provided client
func submitMetrics(ctx context.Context, client *http.Client, r io.Reader, endpoint string) error {
	req, err := http.NewRequest("POST", endpoint, r)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected %v (http.StatusOK) but got %s", http.StatusOK, resp.Status)
	}

	return nil
}
