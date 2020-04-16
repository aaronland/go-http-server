package main

import (
	"context"
	"flag"
	"github.com/aaronland/go-http-server"
	"log"
	"net/http"
	"path/filepath"
	"os"
)

func main() {

	server_uri := flag.String("server-uri", "http://localhost:8080", "...")
	root := flag.String("root", "", "...")

	flag.Parse()

	ctx := context.Background()

	s, err := server.NewServer(ctx, *server_uri)

	if err != nil {
		log.Fatalf("Unable to create server (%s), %v", *server_uri, err)
	}

	if *root == "" {
		log.Fatalf("Missing -root")
	}

	abs_root, err := filepath.Abs(*root)

	if err != nil {
		log.Fatalf("Failed to parse root (%s), %v", *root, err)
	}

	info, err := os.Stat(abs_root)

	if err != nil {
		log.Fatalf("Failed to inspect root (%s), %v", abs_root, err)
	}

	if !info.IsDir(){
		log.Fatalf("Root (%s) is not a directory", abs_root)
	}

	http_root := http.Dir(abs_root)
	fs_handler := http.FileServer(http_root)
	
	mux := http.NewServeMux()
	mux.Handle("/", fs_handler)

	log.Printf("Listening on %s", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}
}
