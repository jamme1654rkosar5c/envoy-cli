package env

import (
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeAliasEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_ENV", Value: "production"},
	}
}

func TestAlias_CreatesAliasKey(t *testing.T) {
	entries := makeAliasEntries()
	opts := DefaultAliasOptions()
	out, err := Alias(entries, "DB_HOST", "DATABASE_HOST", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var found *parser.Entry
	for i := range out {
		if out[i].Key == "DATABASE_HOST" {
			found = &out[i]
		}
	}
	if found == nil {
		t.Fatal("expected alias key DATABASE_HOST to exist")
	}
	if found.Value != "localhost" {
		t.Errorf("expected value %q, got %q", "localhost", found.Value)
	}
	if found.Comment != "alias of DB_HOST" {
		t.Errorf("unexpected comment: %q", found.Comment)
	}
}

func TestAlias_KeepsOriginalByDefault(t *testing.T) {
	entries := makeAliasEntries()
	out, err := Alias(entries, "DB_HOST", "DATABASE_HOST", DefaultAliasOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var hasOrig bool
	for _, e := range out {
		if e.Key == "DB_HOST" {
			hasOrig = true
		}
	}
	if !hasOrig {
		t.Error("expected original key DB_HOST to be retained")
	}
}

func TestAlias_SourceNotFound_ReturnsError(t *testing.T) {
	entries := makeAliasEntries()
	_, err := Alias(entries, "MISSING", "NEW_KEY", DefaultAliasOptions())
	if err == nil {
		t.Fatal("expected error for missing source key")
	}
}

func TestAlias_ExistingKey_NoOverwrite_ReturnsError(t *testing.T) {
	entries := makeAliasEntries()
	_, err := Alias(entries, "DB_HOST", "DB_PORT", DefaultAliasOptions())
	if err == nil {
		t.Fatal("expected error when alias key already exists")
	}
}

func TestAlias_ExistingKey_WithOverwrite_Succeeds(t *testing.T) {
	entries := makeAliasEntries()
	opts := DefaultAliasOptions()
	opts.Overwrite = true
	out, err := Alias(entries, "DB_HOST", "DB_PORT", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out {
		if e.Key == "DB_PORT" && e.Value == "localhost" {
			return
		}
	}
	t.Error("expected DB_PORT to be overwritten with DB_HOST value")
}

func TestAlias_DryRun_DoesNotMutate(t *testing.T) {
	entries := makeAliasEntries()
	opts := DefaultAliasOptions()
	opts.DryRun = true
	out, err := Alias(entries, "DB_HOST", "DATABASE_HOST", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out {
		if e.Key == "DATABASE_HOST" {
			t.Error("dry run should not create alias key")
		}
	}
}
