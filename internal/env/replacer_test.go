package env

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makeReplaceEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_URL", Value: "postgres://localhost:5432/mydb"},
		{Key: "API_URL", Value: "https://localhost/api"},
		{Key: "APP_ENV", Value: "development"},
	}
}

func TestReplace_SubstringMatch(t *testing.T) {
	entries := makeReplaceEntries()
	opts := DefaultReplaceOptions()
	opts.OldValue = "localhost"
	opts.NewValue = "prod.example.com"

	result, count, err := Replace(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 replacements, got %d", count)
	}
	if result[0].Value != "prod.example.com" {
		t.Errorf("DB_HOST: expected 'prod.example.com', got %q", result[0].Value)
	}
}

func TestReplace_ExactMatch_OnlyReplacesFull(t *testing.T) {
	entries := makeReplaceEntries()
	opts := DefaultReplaceOptions()
	opts.OldValue = "localhost"
	opts.NewValue = "prod.example.com"
	opts.ExactMatch = true

	_, count, err := Replace(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Only DB_HOST has value exactly "localhost"
	if count != 1 {
		t.Errorf("expected 1 exact replacement, got %d", count)
	}
}

func TestReplace_KeyFilter_RestrictsScope(t *testing.T) {
	entries := makeReplaceEntries()
	opts := DefaultReplaceOptions()
	opts.OldValue = "localhost"
	opts.NewValue = "remote"
	opts.KeyFilter = "DB_"

	_, count, err := Replace(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 replacements (DB_HOST + DB_URL), got %d", count)
	}
}

func TestReplace_DryRun_DoesNotMutateOriginal(t *testing.T) {
	entries := makeReplaceEntries()
	opts := DefaultReplaceOptions()
	opts.OldValue = "localhost"
	opts.NewValue = "remote"
	opts.DryRun = true

	_, _, err := Replace(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Value != "localhost" {
		t.Errorf("DryRun mutated original: got %q", entries[0].Value)
	}
}

func TestReplace_EmptyOldValue_ReturnsError(t *testing.T) {
	entries := makeReplaceEntries()
	opts := DefaultReplaceOptions()
	opts.OldValue = ""

	_, _, err := Replace(entries, opts)
	if err == nil {
		t.Fatal("expected error for empty OldValue, got nil")
	}
}
