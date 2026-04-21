package env

import (
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeReorderEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "FOO", Value: "1"},
		{Key: "BAR", Value: "2"},
		{Key: "BAZ", Value: "3"},
		{Key: "QUX", Value: "4"},
	}
}

func TestReorder_ExplicitOrder(t *testing.T) {
	entries := makeReorderEntries()
	opts := DefaultReorderOptions()
	opts.Keys = []string{"BAZ", "FOO", "QUX", "BAR"}

	result, err := Reorder(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"BAZ", "FOO", "QUX", "BAR"}
	for i, e := range result {
		if e.Key != expected[i] {
			t.Errorf("position %d: want %q, got %q", i, expected[i], e.Key)
		}
	}
}

func TestReorder_UnknownKeysPushedToEnd(t *testing.T) {
	entries := makeReorderEntries()
	opts := DefaultReorderOptions()
	opts.Keys = []string{"BAR", "FOO"}

	result, err := Reorder(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Key != "BAR" || result[1].Key != "FOO" {
		t.Errorf("first two keys should be BAR, FOO; got %q %q", result[0].Key, result[1].Key)
	}
	// BAZ and QUX should follow in original order
	if result[2].Key != "BAZ" || result[3].Key != "QUX" {
		t.Errorf("trailing keys should be BAZ, QUX; got %q %q", result[2].Key, result[3].Key)
	}
}

func TestReorder_UnknownKeysPrependedWhenFlagFalse(t *testing.T) {
	entries := makeReorderEntries()
	opts := DefaultReorderOptions()
	opts.Keys = []string{"QUX"}
	opts.PushUnknownToEnd = false

	result, err := Reorder(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// QUX should be last since unlisted keys come first
	if result[len(result)-1].Key != "QUX" {
		t.Errorf("expected QUX at end, got %q", result[len(result)-1].Key)
	}
}

func TestReorder_ErrorOnMissing(t *testing.T) {
	entries := makeReorderEntries()
	opts := DefaultReorderOptions()
	opts.Keys = []string{"MISSING_KEY"}
	opts.ErrorOnMissing = true

	_, err := Reorder(entries, opts)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestReorder_NoErrorOnMissingByDefault(t *testing.T) {
	entries := makeReorderEntries()
	opts := DefaultReorderOptions()
	opts.Keys = []string{"MISSING_KEY", "FOO"}

	result, err := Reorder(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(result))
	}
}
