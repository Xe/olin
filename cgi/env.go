package cgi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/iancoleman/strcase"
)

// DeriveEnv creates the base environment variables for
func DeriveEnv(r *http.Request, runID, workerID string) map[string]string {
	result := map[string]string{
		"REQUEST_METHOD":    r.Method,
		"PATH_INFO":         r.RequestURI,
		"QUERY_STRING":      r.URL.Query().Encode(),
		"RUN_ID":            runID,
		"WORKER_ID":         workerID,
		"CONTENT_LENGTH":    fmt.Sprint(r.ContentLength),
		"CONTENT_TYPE":      r.Header.Get("Content-Type"),
		"GATEWAY_INTERFACE": "CGI/1.1",
		"SERVER_SOFTWARE":   "Olin",
		"SERVER_PROTOCOL":   "HTTP/1.1",
		"HTTPS":             "on",
		"HTTP_HOST":         r.Host,
		"SERVER_NAME":       r.Host,
	}

	for name, headers := range r.Header {
		if strings.ToUpper(name) == "PROXY" {
			continue
		}
		result["HTTP_"+strcase.ToScreamingSnake(name)] = headers[0]
	}

	return result
}
