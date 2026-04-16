package env

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makePinEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
		{Key: "DEBUG", Value: "true"},
	}
}

func TestPin_SetsComment(t *testing.T) {
	entries := makePinEntries()
	result, err := Pin(entries, []string{"HOST", "PORT"}, DefaultPinOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Comment != "pinned" {
		t.Errorf("expected HOST to be pinned")
	}
	if result[1].Comment != "pinned" {
		t.Errorf("expected PORT to be pinned")
	}
	if result[2].Comment == "pinned" {
		t.Errorf("expected DEBUG to not be pinned")
	}
}

func TestPin_MissingKey_Error(t *testing.T) {
	entries := makePinEntries()
	_, err := Pin(entries, []string{"MISSING"}, DefaultPinOptions())
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestPin_MissingKey_NoError(t *testing.T) {
	entries := makePinEntries()
	opts := DefaultPinOptions()
	opts.FailIfMissing = false
	_, err := Pin(entries, []string{"MISSING"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnpin_ClearsComment(t *testing.T) {
	entries := makePinEntries()
	entries[0].Comment = "pinned"
	entries[1].Comment = "pinned"

	result, err := Unpin(entries, []string{"HOST"}, DefaultPinOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Comment == "pinned" {
		t.Errorf("expected HOST to be unpinned")
	}
	if result[1].Comment != "pinned" {
		t.Errorf("expected PORT to remain pinned")
	}
}

func TestIsPinned_True(t *testing.T) {
	e := parser.Entry{Key: "X", Value: "1", Comment: "pinned"}
	if !IsPinned(e) {
		t.Error("expected IsPinned to return true")
	}
}

func TestIsPinned_False(t *testing.T) {
	e := parser.Entry{Key: "X", Value: "1", Comment: "other"}
	if IsPinned(e) {
		t.Error("expected IsPinned to return false")
	}
}
