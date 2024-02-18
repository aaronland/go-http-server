package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/aaronland/go-http-server"
)

func TestRouteHandler(t *testing.T) {

	ctx := context.Background()

	pv_func := func(ctx context.Context) (http.Handler, error) {

		fn := func(rsp http.ResponseWriter, req *http.Request) {
			rsp.Write([]byte(req.PathValue("id")))
			return
		}

		return http.HandlerFunc(fn), nil
	}

	pv2_func := func(ctx context.Context) (http.Handler, error) {

		fn := func(rsp http.ResponseWriter, req *http.Request) {
			body := fmt.Sprintf("%s %s", req.PathValue("hello"), req.PathValue("world"))
			rsp.Write([]byte(body))
			return
		}

		return http.HandlerFunc(fn), nil
	}

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

		// Things we expect to succeed
		"/foo":                                foo_func,
		"/foo/bar":                            bar_func,
		"/id/{id}":                            pv_func,
		"/id/{id}/sub":                        pv_func,
		"/{hello}/omg/wtf/{world}":            pv2_func,
		"GET /this/is/a/{hello}/{world}/yeah": pv2_func,

		// Things we expect to fail below
		"POST /foo/post":                   foo_func,
		"example:com/wrong/host/":          bar_func,
		"GET example:com/also/wrong/host/": bar_func,
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

	tests_to_succeed := map[string]string{
		"http://localhost:8080/foo":                        "foo",
		"http://localhost:8080/foo/":                       "foo",
		"http://localhost:8080/foo/bar":                    "bar",
		"http://localhost:8080/id/1234":                    "1234",
		"http://localhost:8080/id/5678/sub":                "5678",
		"http://localhost:8080/horse/omg/wtf/email":        "horse email",
		"http://localhost:8080/this/is/a/GET/handler/yeah": "GET handler",
	}

	for uri, expected := range tests_to_succeed {

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

	tests_to_fail := map[string]string{
		"http://localhost:8080/foo/post":         "",
		"http://localhost:8080/wrong/host/":      "",
		"http://localhost:8080/also/wrong/host/": "",
	}

	for uri, _ := range tests_to_fail {

		rsp, err := http.Get(uri)

		if err != nil {
			t.Fatalf("Failed to query %s", uri)
		}

		defer rsp.Body.Close()

		if rsp.StatusCode == http.StatusOK {
			t.Fatalf("Expected %s to fail (status code)", uri)
		}

	}

}
