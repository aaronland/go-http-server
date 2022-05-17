package server

import (
	"context"
	"testing"
)

func TestLambdaServer(t *testing.T) {

	ctx := context.Background()

	s, err := NewServer(ctx, "lambda://")

	if err != nil {
		t.Fatalf("Failed to create new server, %v", err)
	}

	if s.Address() != "lambda:" {
		t.Fatalf("Unexpected address: %s", s.Address())
	}
}
