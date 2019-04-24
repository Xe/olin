package ln

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"within.website/ln/opname"
)

var (
	// DefaultTimeFormat represents the way in which time will be formatted by default
	DefaultTimeFormat = time.RFC3339
)

// Formatter defines the formatting of events
type Formatter interface {
	Format(ctx context.Context, e Event) ([]byte, error)
}

// DefaultFormatter is the default way in which to format events
var DefaultFormatter Formatter

func init() {
	DefaultFormatter = NewTextFormatter()
}

// TextFormatter formats events as key value pairs.
// Any remaining text not wrapped in an instance of `F` will be
// placed at the end.
type TextFormatter struct {
	TimeFormat string
}

// NewTextFormatter returns a Formatter that outputs as text.
func NewTextFormatter() Formatter {
	return &TextFormatter{TimeFormat: DefaultTimeFormat}
}

// Format implements the Formatter interface
func (t *TextFormatter) Format(ctx context.Context, e Event) ([]byte, error) {
	var writer bytes.Buffer

	writer.WriteString("time=\"")
	writer.WriteString(e.Time.Format(t.TimeFormat))
	writer.WriteString("\"")

	if op, ok := opname.Get(ctx); ok {
		e.Data["operation"] = op
	}

	keys := make([]string, len(e.Data))
	i := 0

	for k := range e.Data {
		keys[i] = k
		i++
	}

	for _, k := range keys {
		v := e.Data[k]

		writer.WriteByte(' ')
		if shouldQuote(k) {
			writer.WriteString(fmt.Sprintf("%q", k))
		} else {
			writer.WriteString(k)
		}

		writer.WriteByte('=')

		switch e := v.(type) {
		case string:
			vs := e
			if shouldQuote(vs) {
				fmt.Fprintf(&writer, "%q", vs)
			} else {
				writer.WriteString(vs)
			}
		case error:
			tmperr := e
			es := tmperr.Error()

			if shouldQuote(es) {
				fmt.Fprintf(&writer, "%q", es)
			} else {
				writer.WriteString(es)
			}
		case time.Time:
			tmptime := e
			writer.WriteString(tmptime.Format(time.RFC3339))
		default:
			fmt.Fprint(&writer, v)
		}
	}

	if len(e.Message) > 0 {
		fmt.Fprintf(&writer, " _msg=%q", e.Message)
	}

	writer.WriteByte('\n')
	return writer.Bytes(), nil
}

func shouldQuote(s string) bool {
	for _, b := range s {
		if !((b >= 'A' && b <= 'Z') ||
			(b >= 'a' && b <= 'z') ||
			(b >= '0' && b <= '9') ||
			(b == '-' || b == '.' || b == '#' ||
				b == '/' || b == '_')) {
			return true
		}
	}
	return false
}

type jsonFormatter struct{}

// JSONFormatter outputs json lines for use with tools like https://github.com/koenbollen/jl.
func JSONFormatter() Formatter {
	return jsonFormatter{}
}

func (j jsonFormatter) Format(ctx context.Context, e Event) ([]byte, error) {
	if op, ok := opname.Get(ctx); ok {
		e.Data["operation"] = op
	}

	e.Data["time"] = e.Time.Format(time.RFC3339)

	data, err := json.Marshal(e.Data)
	return append(data, '\n'), err
}
