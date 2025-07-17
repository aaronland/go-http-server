package main

import (
	"context"
	_ "fmt"
	"log"
	_ "io/fs"
	"net/http"
	"time"
	"os"
	
	"github.com/aaronland/go-http-server/v2"
	"github.com/aaronland/go-http-server/v2/handler"
	"github.com/sfomuseum/go-flags/flagset"
)

func main() {

	var server_uri string
	var cgibin_root string

	fs := flagset.NewFlagSet("server")

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")
	fs.StringVar(&cgibin_root, "cgi-bin-root", "", "...")

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

	cgi_root, err := os.OpenRoot(cgibin_root)

	if err != nil {
		log.Fatalf("Failed to open cgi root, %v", err)
	}

	cgi_fs := cgi_root.FS()
	
	cgi_opts := &handler.CgiBinHandlerOptions{
		Root: cgi_fs,
		Timeout: 30 * time.Second,
	}

	cgi_handler, err := handler.CgiBinHandler(cgi_opts)

	if err != nil {
		log.Fatalf("Failed to create cgibin handler, %v", err)
	}
	
	mux := http.NewServeMux()

	mux.Handle("/cgi-bin/", cgi_handler)
	
	log.Printf("Listening on %s", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}
}
