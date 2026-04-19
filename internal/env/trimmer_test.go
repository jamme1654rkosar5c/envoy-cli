package env

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makeTrimEntries(kvs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(kvs)/2)
	for i := 0; i+1 < len(kvs); i += 2 {
		entries = append(entries, parser.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return entries
}

func TestTrim_TrimsKeyAndValue(t *testing.T) {
	entries := makeTrimEntries("  MY_KEY  ", "  hello  ")
	opts := DefaultTrimOptions()
	out := Trim(entries, opts)
	if out[0].Key != "MY_KEY" {
		t.Errorf("expected trimmed key, got %q", out[0].Key)
	}
	if out[0].Value != "hello" {
		t.Errorf("expected trimmed value, got %q", out[0].Value)
	}
}

func TestTrim_TrimPrefix_RemovesFromKey(t *testing.T) {
	entries := makeTrimEntries("APP_HOST", "localhost", "APP_PORT", "8080")
	opts := DefaultTrimOptions()
	opts.TrimPrefixes = []string{"APP_"}
	out := Trim(entries, opts)
	if out[0].Key != "HOST" {
		t.Errorf("expected HOST, got %q", out[0].Key)
	}
	if out[1].Key != "PORT" {
		t.Errorf("expected PORT, got %q", out[1].Key)
	}
}

func TestTrim_TrimSuffix_RemovesFromKey(t *testing.T) {
	entries := makeTrimEntries("DB_URL_DEV", "postgres://")
	opts := DefaultTrimOptions()
	opts.TrimSuffixes = []string{"_DEV"}
	out := Trim(entries, opts)
	if out[0].Key != "DB_URL" {
		t.Errorf("expected DB_URL, got %q", out[0].Key)
	}
}

func TestTrim_SkipEmpty_LeavesEmptyValueUntouched(t *testing.T) {
	entries := makeTrimEntries("  EMPTY_KEY  ", "")
	opts := DefaultTrimOptions()
	opts.SkipEmpty = true
	out := Trim(entries, opts)
	// key should NOT be trimmed because SkipEmpty skips processing
	if out[0].Key != "  EMPTY_KEY  " {
		t.Errorf("expected untouched key, got %q", out[0].Key)
	}
}

func TestTrim_DoesNotMutateOriginal(t *testing.T) {
	entries := makeTrimEntries("  KEY  ", "  val  ")
	opts := DefaultTrimOptions()
	Trim(entries, opts)
	if entries[0].Key != "  KEY  " {
		t.Error("original entries were mutated")
	}
}
