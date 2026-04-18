package env

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makeSanitizeEntries(pairs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestSanitize_TrimKeysAndValues(t *testing.T) {
	entries := makeSanitizeEntries("  KEY  ", "  value  ")
	opts := DefaultSanitizeOptions()
	out := Sanitize(entries, opts)
	if out[0].Key != "KEY" {
		t.Errorf("expected trimmed key, got %q", out[0].Key)
	}
	if out[0].Value != "value" {
		t.Errorf("expected trimmed value, got %q", out[0].Value)
	}
}

func TestSanitize_RemoveEmpty(t *testing.T) {
	entries := makeSanitizeEntries("KEY1", "val", "KEY2", "")
	opts := DefaultSanitizeOptions()
	opts.RemoveEmpty = true
	out := Sanitize(entries, opts)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Key != "KEY1" {
		t.Errorf("unexpected key %q", out[0].Key)
	}
}

func TestSanitize_NormalizeKeys(t *testing.T) {
	entries := makeSanitizeEntries("my_key", "val")
	opts := DefaultSanitizeOptions()
	opts.NormalizeKeys = true
	out := Sanitize(entries, opts)
	if out[0].Key != "MY_KEY" {
		t.Errorf("expected MY_KEY, got %q", out[0].Key)
	}
}

func TestSanitize_StripQuotes(t *testing.T) {
	entries := makeSanitizeEntries("KEY", `"hello world"`)
	opts := DefaultSanitizeOptions()
	opts.StripQuotes = true
	out := Sanitize(entries, opts)
	if out[0].Value != "hello world" {
		t.Errorf("expected unquoted value, got %q", out[0].Value)
	}
}

func TestSanitize_DoesNotMutateOriginal(t *testing.T) {
	entries := makeSanitizeEntries("  KEY  ", `'val'`)
	opts := DefaultSanitizeOptions()
	opts.StripQuotes = true
	Sanitize(entries, opts)
	if entries[0].Key != "  KEY  " {
		t.Error("original entry was mutated")
	}
}
