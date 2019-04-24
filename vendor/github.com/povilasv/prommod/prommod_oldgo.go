// +build !go1.12

package prommod

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// NewCollector returns a collector which exports metrics about current dependency information.
func NewCollector(program string) *prometheus.GaugeVec {
	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: program,
			Name:      "go_mod_info",
			Help: fmt.Sprintf(
				"A metric with a constant '1' value labeled by dependency name, version, from which %s was built.",
				program,
			),
		},
		[]string{"name", "version"},
	)

	return gauge
}

// Print returns module version information.
func Print(program string) string {
	return program
}

// Info returns dependency versions
func Info() string {
	return "()"
}
