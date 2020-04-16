package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-http-server"
	"log"
	"net/http"
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

	server_uri := flag.String("server-uri", "http://localhost:8080", "...")

	flag.Parse()

	ctx := context.Background()

	s, err := server.NewServer(ctx, *server_uri)

	if err != nil {
		log.Fatalf("Unable to create server (%s), %v", *server_uri, err)
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
