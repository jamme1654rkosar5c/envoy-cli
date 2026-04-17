package env

import (
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeClonerEntries(kvs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		entries = append(entries, parser.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return entries
}

func TestClone_AppendsNewKeys(t *testing.T) {
	dst := makeClonerEntries("HOST", "localhost")
	src := makeClonerEntries("PORT", "8080", "DEBUG", "true")
	out, err := Clone(dst, src, DefaultCloneOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}

func TestClone_WithPrefix(t *testing.T) {
	dst := makeClonerEntries()
	src := makeClonerEntries("KEY", "val")
	opts := DefaultCloneOptions()
	opts.Prefix = "PROD_"
	out, err := Clone(dst, src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Key != "PROD_KEY" {
		t.Errorf("expected PROD_KEY, got %s", out[0].Key)
	}
}

func TestClone_SkipKeys(t *testing.T) {
	dst := makeClonerEntries()
	src := makeClonerEntries("A", "1", "B", "2")
	opts := DefaultCloneOptions()
	opts.SkipKeys = map[string]bool{"A": true}
	out, err := Clone(dst, src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0].Key != "B" {
		t.Errorf("expected only B, got %+v", out)
	}
}

func TestClone_ConflictError(t *testing.T) {
	dst := makeClonerEntries("KEY", "old")
	src := makeClonerEntries("KEY", "new")
	_, err := Clone(dst, src, DefaultCloneOptions())
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
}

func TestClone_OverwriteExisting(t *testing.T) {
	dst := makeClonerEntries("KEY", "old")
	src := makeClonerEntries("KEY", "new")
	opts := DefaultCloneOptions()
	opts.OverwriteExisting = true
	out, err := Clone(dst, src, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "new" {
		t.Errorf("expected value 'new', got %s", out[0].Value)
	}
}
