package env

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makeDedupEntries(pairs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestDedupe_NoDuplicates(t *testing.T) {
	entries := makeDedupEntries("A", "1", "B", "2", "C", "3")
	out, removed := Dedupe(entries, DefaultDedupeOptions())
	if len(out) != 3 || len(removed) != 0 {
		t.Fatalf("expected 3 entries and 0 removed, got %d and %d", len(out), len(removed))
	}
}

func TestDedupe_KeepFirst(t *testing.T) {
	entries := makeDedupEntries("A", "first", "B", "2", "A", "second")
	out, removed := Dedupe(entries, DedupeOptions{Strategy: KeepFirst})
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Value != "first" {
		t.Errorf("expected 'first', got %q", out[0].Value)
	}
	if len(removed) != 1 || removed[0] != "A" {
		t.Errorf("expected removed [A], got %v", removed)
	}
}

func TestDedupe_KeepLast(t *testing.T) {
	entries := makeDedupEntries("A", "first", "B", "2", "A", "second")
	out, removed := Dedupe(entries, DedupeOptions{Strategy: KeepLast})
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Value != "second" {
		t.Errorf("expected 'second', got %q", out[0].Value)
	}
	if len(removed) != 1 || removed[0] != "A" {
		t.Errorf("expected removed [A], got %v", removed)
	}
}

func TestDedupe_MultipleOccurrences(t *testing.T) {
	entries := makeDedupEntries("X", "1", "X", "2", "X", "3")
	out, removed := Dedupe(entries, DefaultDedupeOptions())
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if len(removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(removed))
	}
}
