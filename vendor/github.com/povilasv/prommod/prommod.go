// +build go1.12

package prommod

import (
	"bytes"
	"fmt"
	"html/template"
	"runtime/debug"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// versionInfoTmpl contains the template used by Print call.
var versionInfoTmpl = `{{.Program}}
{{range $k,$v := .Deps}}  {{$k}}:  {{$v}}
{{end}}`

// build module dependency information. Populated at build-time.
var (
	buildInfo, ok = debug.ReadBuildInfo()
	info          string
	version       map[string]string

	tmpl = template.Must(template.New("version").Parse(versionInfoTmpl))
)

func init() {
	var versions []string
	if ok {
		for _, dep := range buildInfo.Deps {
			d := dep
			if dep.Replace != nil {
				d = dep.Replace
			}
			versions = append(versions, d.Path+": "+d.Version)
		}
	}

	info = fmt.Sprintf("(%s)", strings.Join(versions, ", "))

	version = make(map[string]string)
	if ok {
		for _, dep := range buildInfo.Deps {
			d := dep
			if dep.Replace != nil {
				d = dep.Replace
			}
			version[d.Path] = d.Version
		}
	}
}

// NewCollector returns a collector which exports metrics about current dependency information.
func NewCollector(program string) *prometheus.GaugeVec {
	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "go_mod_info",
			Help: fmt.Sprintf(
				"A metric with a constant '1' value labeled by dependency name, version, from which %s was built.",
				program,
			),
		},
		[]string{"name", "version", "program"},
	)
	if !ok {
		return gauge
	}

	for _, dep := range buildInfo.Deps {
		d := dep
		if dep.Replace != nil {
			d = dep.Replace
		}
		gauge.WithLabelValues(d.Path, d.Version, program).Set(1)
	}
	return gauge
}

type versionPrint struct {
	Program string
	Deps    map[string]string
}

// Print returns module version information.
func Print(program string) string {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "version", versionPrint{
		Program: program,
		Deps:    version,
	}); err != nil {
		panic(err)
	}
	return strings.TrimSpace(buf.String())
}

// Info returns dependency versions
func Info() string {
	return info
}
