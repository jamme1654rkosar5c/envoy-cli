package env

import (
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeFlattenEntries(pairs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestFlatten_NoSeparatorInKey(t *testing.T) {
	entries := makeFlattenEntries("APP_HOST", "localhost", "APP_PORT", "5432")
	opts := DefaultFlattenOptions()
	out, err := Flatten(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Key != "APP_HOST" {
		t.Errorf("expected APP_HOST, got %s", out[0].Key)
	}
}

func TestFlatten_CollapsesDoubleSeparator(t *testing.T) {
	entries := makeFlattenEntries("APP__DB__HOST", "db.local")
	opts := DefaultFlattenOptions()
	out, err := Flatten(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Key != "APP_DB_HOST" {
		t.Errorf("expected APP_DB_HOST, got %s", out[0].Key)
	}
	if out[0].Value != "db.local" {
		t.Errorf("expected db.local, got %s", out[0].Value)
	}
}

func TestFlatten_DeduplicatesResultKeys_KeepsLast(t *testing.T) {
	entries := makeFlattenEntries("APP__HOST", "first", "APP__HOST", "second")
	opts := DefaultFlattenOptions()
	out, err := Flatten(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 deduplicated entry, got %d", len(out))
	}
	if out[0].Value != "second" {
		t.Errorf("expected 'second', got %s", out[0].Value)
	}
}

func TestFlatten_PrefixFilter_SkipsNonMatching(t *testing.T) {
	entries := makeFlattenEntries("APP__KEY", "yes", "OTHER__KEY", "no")
	opts := DefaultFlattenOptions()
	opts.Prefix = "APP"
	out, err := Flatten(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Key != "APP_KEY" {
		t.Errorf("expected APP_KEY, got %s", out[0].Key)
	}
	if out[1].Key != "OTHER__KEY" {
		t.Errorf("expected OTHER__KEY unchanged, got %s", out[1].Key)
	}
}

func TestFlatten_EmptySeparator_ReturnsError(t *testing.T) {
	entries := makeFlattenEntries("A", "1")
	opts := DefaultFlattenOptions()
	opts.Separator = ""
	_, err := Flatten(entries, opts)
	if err == nil {
		t.Fatal("expected error for empty separator")
	}
}
