package env

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makeTagEntries(kvs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		entries = append(entries, parser.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return entries
}

func TestTag_SetsComment(t *testing.T) {
	opts := DefaultTagOptions()
	entries := makeTagEntries("APP_ENV", "production")
	out, err := Tag(entries, "APP_ENV", "stable", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := GetTag(out, "APP_ENV", opts)
	if got != "stable" {
		t.Errorf("expected tag 'stable', got %q", got)
	}
}

func TestTag_DoesNotOverwriteByDefault(t *testing.T) {
	opts := DefaultTagOptions()
	entries := makeTagEntries("APP_ENV", "production")
	out, _ := Tag(entries, "APP_ENV", "stable", opts)
	out, _ = Tag(out, "APP_ENV", "beta", opts)
	got := GetTag(out, "APP_ENV", opts)
	if got != "stable" {
		t.Errorf("expected original tag 'stable', got %q", got)
	}
}

func TestTag_OverwriteReplacesTag(t *testing.T) {
	opts := DefaultTagOptions()
	opts.Overwrite = true
	entries := makeTagEntries("APP_ENV", "production")
	out, _ := Tag(entries, "APP_ENV", "stable", opts)
	out, _ = Tag(out, "APP_ENV", "beta", opts)
	got := GetTag(out, "APP_ENV", opts)
	if got != "beta" {
		t.Errorf("expected tag 'beta', got %q", got)
	}
}

func TestUntag_ClearsTag(t *testing.T) {
	opts := DefaultTagOptions()
	entries := makeTagEntries("DB_HOST", "localhost")
	out, _ := Tag(entries, "DB_HOST", "infra", opts)
	out = Untag(out, "DB_HOST", opts)
	got := GetTag(out, "DB_HOST", opts)
	if got != "" {
		t.Errorf("expected empty tag, got %q", got)
	}
}

func TestGetTag_MissingKey_ReturnsEmpty(t *testing.T) {
	opts := DefaultTagOptions()
	entries := makeTagEntries("APP_ENV", "production")
	got := GetTag(entries, "NONEXISTENT", opts)
	if got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestTag_DoesNotMutateOriginal(t *testing.T) {
	opts := DefaultTagOptions()
	entries := makeTagEntries("APP_ENV", "production")
	_, _ = Tag(entries, "APP_ENV", "stable", opts)
	if entries[0].Comment != "" {
		t.Error("original entries were mutated")
	}
}
