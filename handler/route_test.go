package handler

import (
	"context"
	_ "fmt"
	"io"
	"net/http"
	"testing"

	"github.com/aaronland/go-http-server"
)

func TestRouteHandler(t *testing.T) {

	ctx := context.Background()

	foo_func := func(ctx context.Context) (http.Handler, error) {

		fn := func(rsp http.ResponseWriter, req *http.Request) {
			rsp.Write([]byte(`foo`))
			return
		}

		return http.HandlerFunc(fn), nil
	}

	bar_func := func(ctx context.Context) (http.Handler, error) {

		fn := func(rsp http.ResponseWriter, req *http.Request) {
			rsp.Write([]byte(`bar`))
			return
		}

		return http.HandlerFunc(fn), nil
	}

	handlers := map[string]RouteHandlerFunc{
		"/foo":     foo_func,
		"/foo/bar": bar_func,
	}

	route_handler, err := RouteHandler(handlers)

	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", route_handler)

	s, err := server.NewServer(ctx, "http://localhost:8080")

	if err != nil {
		t.Fatal(err)
	}

	go func() {
		s.ListenAndServe(ctx, mux)
	}()

	tests := map[string]string{
		"http://localhost:8080/foo":     "foo",
		"http://localhost:8080/foo/":    "foo",
		"http://localhost:8080/foo/bar": "bar",
	}

	for uri, expected := range tests {

		rsp, err := http.Get(uri)

		if err != nil {
			t.Fatalf("Failed to get %s, %v", uri, err)
		}

		defer rsp.Body.Close()

		body, err := io.ReadAll(rsp.Body)

		if err != nil {
			t.Fatalf("Failed to read body for %s, %v", uri, err)
		}

		str_body := string(body)

		if str_body != expected {
			t.Fatalf("Unexpected value for %s. Expected '%s' but got '%s'", uri, expected, str_body)
		}
	}
}
