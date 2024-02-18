package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aaronland/go-http-server"
	"github.com/sfomuseum/go-flags/flagset"
)

func NewHandler() http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		msg := fmt.Sprintf("Hello, %s", req.Host)
		rsp.Write([]byte(msg))
	}

	h := http.HandlerFunc(fn)
	return h
}

func main() {

	var server_uri string

	fs := flagset.NewFlagSet("server")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "AARONLAND")

	if err != nil {
		log.Fatalf("Failed to set flags from environment variables, %v", err)
	}

	ctx := context.Background()

	s, err := server.NewServer(ctx, server_uri)

	if err != nil {
		log.Fatalf("Unable to create server (%s), %v", server_uri, err)
	}

	mux := http.NewServeMux()
	handler := NewHandler()

	mux.Handle("/", handler)

	log.Printf("Listening on %s", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}
}
