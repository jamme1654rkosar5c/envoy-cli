package env

import (
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makePromoteEntries(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestPromote_AddsNewKeys(t *testing.T) {
	src := makePromoteEntries("NEW_KEY", "hello")
	dst := makePromoteEntries("EXISTING", "world")
	out, err := Promote(src, dst, DefaultPromoteOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestPromote_ConflictNoOverwrite_ReturnsError(t *testing.T) {
	src := makePromoteEntries("KEY", "new")
	dst := makePromoteEntries("KEY", "old")
	_, err := Promote(src, dst, DefaultPromoteOptions())
	if err == nil {
		t.Fatal("expected error on conflict without Overwrite")
	}
}

func TestPromote_ConflictWithOverwrite_ReplacesValue(t *testing.T) {
	src := makePromoteEntries("KEY", "new")
	dst := makePromoteEntries("KEY", "old")
	opts := DefaultPromoteOptions()
	opts.Overwrite = true
	out, err := Promote(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "new" {
		t.Errorf("expected 'new', got %q", out[0].Value)
	}
}

func TestPromote_DryRun_DoesNotMutate(t *testing.T) {
	src := makePromoteEntries("EXTRA", "val")
	dst := makePromoteEntries("BASE", "base")
	opts := DefaultPromoteOptions()
	opts.DryRun = true
	out, err := Promote(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("dry run should not add keys, got %d entries", len(out))
	}
}

func TestPromote_FilterByKeys(t *testing.T) {
	src := makePromoteEntries("A", "1", "B", "2", "C", "3")
	dst := makePromoteEntries("X", "x")
	opts := DefaultPromoteOptions()
	opts.Keys = []string{"A", "C"}
	out, err := Promote(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries (X,A,C), got %d", len(out))
	}
}
