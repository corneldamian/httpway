# httpway

Simple middleware for [httprouter](https://github.com/julienschmidt/httprouter/)

- simple middleware without overhead, httprouter will have the same performance
- context available from first middleware until handler
- http server with gracefully shutdown

You can get some middlewares from here [httpwaymid](https://github.com/corneldamian/httpwaymid)
Integrates very well with [golog](https://github.com/corneldamian/golog.git)

[![GoDoc](https://godoc.org/github.com/corneldamian/httpway?status.svg)](https://godoc.org/github.com/corneldamian/httpway)
[![Build Status](https://travis-ci.org/corneldamian/httpway.svg?branch=master)](https://travis-ci.org/corneldamian/httpway)

```
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/corneldamian/httpway"
	"github.com/julienschmidt/httprouter"
)

var server *httpway.Server

func main() {
	router := httpway.New()

	public := router.Middleware(AccessLogger)
	private := public.Middleware(AuthCheck)

	public.GET("/public", testHandler("public"))

	private.GET("/private", testHandler("private"))
	private.GET("/stop", stopServer)

	server = httpway.NewServer(nil)
	server.Addr = ":8080"
	server.Handler = router

	if err := server.Start(); err != nil {
		fmt.Println("Error", err)
		return
	}

	if err := server.WaitStop(10 * time.Second); err != nil {
		fmt.Println("Error", err)
	}
}

func testHandler(str string) httpway.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s", str)
	}
}

func stopServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Stopping")
	server.Stop()
}

func AccessLogger(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	httpway.GetContext(r).Next(w, r)

	fmt.Printf("Request: %s duration: %s\n", r.URL.EscapedPath(), time.Since(startTime))
}

func AuthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := httpway.GetContext(r)

	if r.URL.EscapedPath() == "/public" {
		http.Error(w, "Auth required", http.StatusForbidden)
		return
	}

	ctx.Next(w, r)
}

```


