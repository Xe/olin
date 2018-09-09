/* Copyright (c) 2018 Salesforce
 * All rights reserved.
 * Licensed under the BSD 3-Clause license.
 * For full license text, see LICENSE.txt file in the repo root  or https://opensource.org/licenses/BSD-3-Clause
 */

/*
Package onload automatically starts hmetrics reporting, ignoring errors and
retrying reporting, backing off in 10 second increments.

Use this package when you don't care about stopping reporting, specifying the
endpoint, or being notified of any reporting errors.

This package will log all reporting errors via the stdlib log package if the
environment variable HMETRICS_VERBOSE is set to "1", or any other true-like
value as defined by strconv#ParseBool (see $ godoc strconv for more
information).

usage:

  import (
	_ "github.com/heroku/x/hmetrics/onload"
  )

See the package documentation at https://godoc.org/github.com/heroku/x/hmetrics
for more info about Heroku Go metrics and advanced usage.

*/
package onload

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/heroku/x/hmetrics"
)

const (
	interval = 10
)

func init() {
	var logger hmetrics.ErrHandler = func(_ error) error { return nil }

	val := os.Getenv("HMETRICS_VERBOSE")
	should, err := strconv.ParseBool(val)
	if err == nil && should {
		logger = func(err error) error {
			log.Printf("[hmetrics] error: %v", err)
			return nil
		}
	}

	go func() {
		var backoff int64
		for backoff = 1; ; backoff++ {
			start := time.Now()

			err = hmetrics.Report(context.Background(), hmetrics.DefaultEndpoint, logger)
			if time.Since(start) > 5*time.Minute {
				backoff = 1
			}
			logger(err)

			time.Sleep(time.Duration(backoff*interval) * time.Second)
		}
	}()
}
