package repository

import (
	"context"
	"errors"
	"testing"
)

func TestCreateWithCancelledContext(t *testing.T) {
	store := NewMemory()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := store.Create(ctx, "Should not be created")

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}

	tasks, err := store.GetAll(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != 0 {
		t.Fatalf("expected zero tasks, got %d", len(tasks))
	}
}