package server

import (
	"context"
	"testing"
)

func TestLambdaFunctionURLServer(t *testing.T) {

	ctx := context.Background()

	s, err := NewServer(ctx, "functionurl://")

	if err != nil {
		t.Fatalf("Failed to create new server, %v", err)
	}

	if s.Address() != "functionurl://" {
		t.Fatalf("Unexpected address: %s", s.Address())
	}
}
