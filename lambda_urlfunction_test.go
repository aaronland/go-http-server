package server

import (
	"context"
	"testing"
)

func TestLambdaURLFunctionServer(t *testing.T) {

	ctx := context.Background()

	s, err := NewServer(ctx, "urlfunction://")

	if err != nil {
		t.Fatalf("Failed to create new server, %v", err)
	}

	if s.Address() != "urlfunction://" {
		t.Fatalf("Unexpected address: %s", s.Address())
	}
}
