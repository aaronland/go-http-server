package server

// https://medium.com/@simonfrey/go-as-in-golang-standard-net-http-config-will-break-your-production-environment-1360871cb72b

// https://ieftimov.com/post/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/

// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
// https://blog.cloudflare.com/exposing-go-on-the-internet/

import (
	"context"
	_ "log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func init() {
	ctx := context.Background()
	RegisterServer(ctx, "http", NewHTTPServer)
}

type HTTPServer struct {
	Server
	url         *url.URL
	http_server *http.Server
}

func NewHTTPServer(ctx context.Context, uri string) (Server, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	u.Scheme = "http"

	read_timeout := 1 * time.Second
	write_timeout := 1 * time.Second
	idle_timeout := 30 * time.Second
	header_timeout := 2 * time.Second

	q := u.Query()

	if q.Get("read_timeout") != "" {

		to, err := strconv.Atoi(q.Get("read_timeout"))

		if err != nil {
			return nil, err
		}

		read_timeout = time.Duration(to) * time.Second
	}

	if q.Get("write_timeout") != "" {

		to, err := strconv.Atoi(q.Get("write_timeout"))

		if err != nil {
			return nil, err
		}

		write_timeout = time.Duration(to) * time.Second
	}

	if q.Get("idle_timeout") != "" {

		to, err := strconv.Atoi(q.Get("idle_timeout"))

		if err != nil {
			return nil, err
		}

		idle_timeout = time.Duration(to) * time.Second
	}

	if q.Get("header_timeout") != "" {

		to, err := strconv.Atoi(q.Get("header_timeout"))

		if err != nil {
			return nil, err
		}

		header_timeout = time.Duration(to) * time.Second
	}

	srv := &http.Server{
		Addr:              u.String(),
		ReadTimeout:       read_timeout,
		WriteTimeout:      write_timeout,
		IdleTimeout:       idle_timeout,
		ReadHeaderTimeout: header_timeout,
	}

	server := HTTPServer{
		url:         u,
		http_server: srv,
	}

	return &server, nil
}

func (s *HTTPServer) Address() string {
	return s.url.String()
}

func (s *HTTPServer) ListenAndServe(ctx context.Context, mux *http.ServeMux) error {
	return s.ListenAndServe(ctx, mux)
	// return http.ListenAndServe(s.url.Host, mux)
}
