package server

import (
	"context"
	"io"
	"log"
	"net/http"
	"testing"
)

func testHandler() http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {
		rsp.Write([]byte("Hello world"))
	}

	h := http.HandlerFunc(fn)
	return h
}

func TestHTTPServer(t *testing.T) {

	ctx := context.Background()

	s, err := NewServer(ctx, "http://localhost:8080")

	if err != nil {
		t.Fatalf("Failed to create server, %v", err)
	}

	go func() {
		err := s.ListenAndServe(ctx, testHandler())

		if err != nil {
			log.Fatalf("Failed to start server, %v", err)
		}
	}()

	rsp, err := http.Get("http://localhost:8080")

	if err != nil {
		t.Fatalf("Failed to GET request, %v", err)
	}

	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)

	if err != nil {
		t.Fatalf("Failed to read response, %v", err)
	}

	if string(body) != "Hello world" {
		t.Fatalf("Unexpected response")
	}

}
