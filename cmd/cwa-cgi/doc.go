/*
Command cwa-cgi is a simple wrapper for applications that handle HTTP requests
mediated by package internal/cwagi.

Usage is simple:

    cwa-cgi [options] <file.wasm>

      -addr string
	    TCP host:port to listen on (default ":8400")
      -main-func string
	    main function to call (because rust takes over the name main) (default "main")
      -max-pool-size int
	    maximum worker pool size (default 32)
      -pool-size int
	    initial worker pool size (default 1)

Some notes about the protocol for requests getting in:

- The overall flow for requests is documented here: https://github.com/Xe/olin/blob/master/doc/diagrams/cwagi-request-flow.dot
- This part actually serves the requests: https://github.com/Xe/olin/blob/master/internal/cwagi/cwagi.go#L60
- The request method is in the environment variable REQUEST_METHOD: https://github.com/Xe/olin/blob/master/internal/cwagi/cwagi.go#L82
- The request URI is in the environment variable REQUEST_URI: https://github.com/Xe/olin/blob/master/internal/cwagi/cwagi.go#L83
- The query string (if applicable) is in the environment variable QUERY_STRING: https://github.com/Xe/olin/blob/master/internal/cwagi/cwagi.go#L84
- The unique run ID is in the environment variable RUN_ID: https://github.com/Xe/olin/blob/master/internal/cwagi/cwagi.go#L85
- The unique worker ID is in the environment variable WORKER_ID: https://github.com/Xe/olin/blob/master/internal/cwagi/cwagi.go#L86
- The process must write a valid HTTP/1.1 request to its standard output: https://github.com/Xe/olin/blob/master/internal/cwagi/cwagi.go#L110
- The pooling logic will automatically spin up more VM's if there is need: https://github.com/Xe/olin/blob/master/internal/cwagi/vmpool.go
*/
package main

//go:generate ./generate_readme.sh
