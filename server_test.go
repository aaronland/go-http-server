package server

import (
	"context"
	"testing"
)

func TestRegisterServer(t *testing.T) {

	ctx := context.Background()

	err := RegisterServer(ctx, "http", NewHTTPServer)

	if err == nil {
		t.Fatalf("Expected NewEncryptedServer to be registered already")
	}
}

func TestNewServer(t *testing.T) {

	ctx := context.Background()

	uri := "lambda://"

	_, err := NewServer(ctx, uri)

	if err != nil {
		t.Fatalf("Failed to create new server for '%s', %v", uri, err)
	}
}
