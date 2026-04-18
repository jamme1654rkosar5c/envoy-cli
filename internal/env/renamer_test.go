package env

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makeRenameEntries() []parser.EnvEntry {
	return []parser.EnvEntry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost", Comment: "primary db"},
	}
}

func TestRename_Success(t *testing.T) {
	entries := makeRenameEntries()
	result, err := Rename(entries, "APP_NAME", "SERVICE_NAME", DefaultRenameOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Key != "SERVICE_NAME" {
		t.Errorf("expected SERVICE_NAME, got %s", result[0].Key)
	}
	if result[0].Value != "envoy" {
		t.Errorf("expected value to be preserved, got %s", result[0].Value)
	}
}

func TestRename_PreservesComment(t *testing.T) {
	entries := makeRenameEntries()
	result, err := Rename(entries, "DB_HOST", "DATABASE_HOST", DefaultRenameOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[2].Comment != "primary db" {
		t.Errorf("expected comment to be preserved, got %q", result[2].Comment)
	}
}

func TestRename_NotFound_Error(t *testing.T) {
	entries := makeRenameEntries()
	_, err := Rename(entries, "MISSING_KEY", "NEW_KEY", DefaultRenameOptions())
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRename_NotFound_NoError(t *testing.T) {
	entries := makeRenameEntries()
	opts := DefaultRenameOptions()
	opts.FailIfNotFound = false
	result, err := Rename(entries, "MISSING_KEY", "NEW_KEY", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != len(entries) {
		t.Errorf("expected same length entries")
	}
}

func TestRename_DestExists_Error(t *testing.T) {
	entries := makeRenameEntries()
	_, err := Rename(entries, "APP_NAME", "APP_ENV", DefaultRenameOptions())
	if err == nil {
		t.Fatal("expected error when dest key exists")
	}
}

func TestRename_DryRun_DoesNotMutate(t *testing.T) {
	entries := makeRenameEntries()
	opts := DefaultRenameOptions()
	opts.DryRun = true
	_, err := Rename(entries, "APP_NAME", "SERVICE_NAME", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[0].Key != "APP_NAME" {
		t.Errorf("dry run should not mutate original, got %s", entries[0].Key)
	}
}
