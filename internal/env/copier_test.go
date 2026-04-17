package env

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makeCopyEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "DB_HOST", Value: "localhost"},
	}
}

func TestCopy_CopiesValue(t *testing.T) {
	entries := makeCopyEntries()
	opts := DefaultCopyOptions()

	result, err := Copy(entries, "APP_NAME", "APP_ALIAS", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var found bool
	for _, e := range result {
		if e.Key == "APP_ALIAS" && e.Value == "envoy" {
			found = true
		}
	}
	if !found {
		t.Error("expected APP_ALIAS=envoy in result")
	}
}

func TestCopy_SourceKeyMissing_ReturnsError(t *testing.T) {
	entries := makeCopyEntries()
	_, err := Copy(entries, "MISSING", "NEW_KEY", DefaultCopyOptions())
	if err == nil {
		t.Error("expected error for missing source key")
	}
}

func TestCopy_DestExists_NoOverwrite_ReturnsError(t *testing.T) {
	entries := makeCopyEntries()
	opts := DefaultCopyOptions()
	opts.Overwrite = false

	_, err := Copy(entries, "APP_NAME", "APP_ENV", opts)
	if err == nil {
		t.Error("expected error when destination exists and Overwrite=false")
	}
}

func TestCopy_DestExists_WithOverwrite_Succeeds(t *testing.T) {
	entries := makeCopyEntries()
	opts := DefaultCopyOptions()
	opts.Overwrite = true

	result, err := Copy(entries, "APP_NAME", "APP_ENV", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, e := range result {
		if e.Key == "APP_ENV" && e.Value != "envoy" {
			t.Errorf("expected APP_ENV=envoy, got %q", e.Value)
		}
	}
}

func TestCopy_KeepSourceFalse_RemovesSource(t *testing.T) {
	entries := makeCopyEntries()
	opts := DefaultCopyOptions()
	opts.KeepSource = false

	result, err := Copy(entries, "APP_NAME", "APP_ALIAS", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, e := range result {
		if e.Key == "APP_NAME" {
			t.Error("expected APP_NAME to be removed")
		}
	}
}

func TestCopy_DryRun_DoesNotMutate(t *testing.T) {
	entries := makeCopyEntries()
	opts := DefaultCopyOptions()
	opts.DryRun = true

	result, err := Copy(entries, "APP_NAME", "APP_ALIAS", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != len(entries) {
		t.Errorf("expected no changes in dry run, got %d entries", len(result))
	}
}
