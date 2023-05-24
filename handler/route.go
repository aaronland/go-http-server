package handler

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
)

type RouteHandlerFunc func() (http.Handler, error)

type RouteHandlerOptions struct {
	Handlers map[string]RouteHandlerFunc
	Patterns []string
	Matches  *sync.Map
	Logger   *log.Logger
}

func RouteHandler(handlers map[string]RouteHandlerFunc) (http.Handler, error) {

	matches := new(sync.Map)
	patterns := make([]string, 0)

	for p, _ := range handlers {
		patterns = append(patterns, p)
	}

	// Sort longest to shortest
	// https://cs.opensource.google/go/go/+/refs/tags/go1.20.4:src/net/http/server.go;l=2533

	sort.Slice(patterns, func(i, j int) bool {
		return len(patterns[i]) > len(patterns[j])
	})

	logger := log.Default()

	opts := &RouteHandlerOptions{
		Matches:  matches,
		Patterns: patterns,
		Handlers: handlers,
		Logger:   logger,
	}

	return RouteHandlerWithOptions(opts)
}

func RouteHandlerWithOptions(opts *RouteHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		var handler http.Handler

		path := req.URL.Path

		v, ok := opts.Matches.Load(path)

		if ok {
			handler = v.(http.Handler)
		} else {

			h, err := deriveHandler(opts, path)

			if err != nil {
				opts.Logger.Printf("%v", err)
				http.Error(rsp, "Internal server error", http.StatusInternalServerError)
				return
			}

			// Don't fill up the matches cache with 404 handlers

			if h == nil {
				http.Error(rsp, "Not found", http.StatusNotFound)
				return
			}

			handler = h
			opts.Matches.Store(path, handler)
		}

		handler.ServeHTTP(rsp, req)
		return
	}

	return http.HandlerFunc(fn), nil
}

func deriveHandler(opts *RouteHandlerOptions, path string) (http.Handler, error) {

	// Basically do what the default http.ServeMux does but inflate the
	// handler (func) on demand at run-time. Handler is cached above.
	// https://cs.opensource.google/go/go/+/refs/tags/go1.20.4:src/net/http/server.go;l=2363

	var matching_pattern string

	for _, p := range opts.Patterns {
		if strings.HasPrefix(path, p) {
			matching_pattern = p
			break
		}
	}

	if matching_pattern == "" {
		return nil, nil
	}

	handler_func, ok := opts.Handlers[matching_pattern]

	if !ok {
		return nil, nil
	}

	handler, err := handler_func()

	if err != nil {
		return nil, fmt.Errorf("Failed to instantiate handler func for '%s' matching '%s', %v", path, matching_pattern, err)
	}

	return handler, nil
}
