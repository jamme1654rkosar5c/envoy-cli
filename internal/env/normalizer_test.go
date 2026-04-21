package env

import (
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeNormEntries(pairs ...string) []parser.EnvEntry {
	entries := make([]parser.EnvEntry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.EnvEntry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	entries := makeNormEntries("db_host", "localhost", "api_key", "secret")
	opts := DefaultNormalizeOptions()
	opts.UppercaseKeys = true

	result := Normalize(entries, opts)
	if result[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", result[0].Key)
	}
	if result[1].Key != "API_KEY" {
		t.Errorf("expected API_KEY, got %s", result[1].Key)
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	entries := makeNormEntries("KEY", "  value  ")
	opts := DefaultNormalizeOptions()
	opts.TrimValues = true

	result := Normalize(entries, opts)
	if result[0].Value != "value" {
		t.Errorf("expected 'value', got %q", result[0].Value)
	}
}

func TestNormalize_StripQuotes(t *testing.T) {
	entries := makeNormEntries("KEY", `"hello world"`, "OTHER", "'single'")
	opts := DefaultNormalizeOptions()
	opts.StripQuotes = true

	result := Normalize(entries, opts)
	if result[0].Value != "hello world" {
		t.Errorf("expected 'hello world', got %q", result[0].Value)
	}
	if result[1].Value != "single" {
		t.Errorf("expected 'single', got %q", result[1].Value)
	}
}

func TestNormalize_RemoveEmpty(t *testing.T) {
	entries := makeNormEntries("KEY", "value", "EMPTY", "", "ANOTHER", "x")
	opts := DefaultNormalizeOptions()
	opts.RemoveEmpty = true

	result := Normalize(entries, opts)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestNormalize_CollapseWhitespace(t *testing.T) {
	entries := makeNormEntries("KEY", "hello   world  foo")
	opts := DefaultNormalizeOptions()
	opts.CollapseWhitespace = true

	result := Normalize(entries, opts)
	if result[0].Value != "hello world foo" {
		t.Errorf("expected 'hello world foo', got %q", result[0].Value)
	}
}

func TestNormalize_DoesNotMutateOriginal(t *testing.T) {
	entries := makeNormEntries("db_host", "  localhost  ")
	opts := DefaultNormalizeOptions()

	_ = Normalize(entries, opts)

	if entries[0].Key != "db_host" {
		t.Errorf("original key was mutated")
	}
	if entries[0].Value != "  localhost  " {
		t.Errorf("original value was mutated")
	}
}
