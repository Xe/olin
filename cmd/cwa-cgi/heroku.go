// +build heroku

package main

import (
	"context"
	"log"

	"github.com/heroku/x/hmetrics"
)

func init() {
	go func() {
		if err := hmetrics.Report(context.Background(), hmetrics.DefaultEndpoint, func(err error) error {
			log.Println("Error reporting metrics to heroku:", err)
			return nil
		}); err != nil {
			log.Fatal("Error starting hmetrics reporting:", err)
		}
	}()
}
