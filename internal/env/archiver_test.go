package env

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeArchiveEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASS", Value: "secret"},
		{Key: "APP_ENV", Value: "production"},
	}
}

func TestArchive_RenamesKey(t *testing.T) {
	entries := makeArchiveEntries()
	opts := DefaultArchiveOptions()
	out, err := Archive(entries, []string{"DB_HOST"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out {
		if e.Key == "DB_HOST" {
			t.Error("original key should be removed")
		}
	}
	found := false
	for _, e := range out {
		if e.Key == "ARCHIVED_DB_HOST" {
			found = true
			if e.Value != "localhost" {
				t.Errorf("expected value localhost, got %s", e.Value)
			}
		}
	}
	if !found {
		t.Error("archived key not found")
	}
}

func TestArchive_KeepsOriginal_WhenRemoveOriginalFalse(t *testing.T) {
	entries := makeArchiveEntries()
	opts := DefaultArchiveOptions()
	opts.RemoveOriginal = false
	out, err := Archive(entries, []string{"DB_PASS"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	origFound, archFound := false, false
	for _, e := range out {
		if e.Key == "DB_PASS" {
			origFound = true
		}
		if e.Key == "ARCHIVED_DB_PASS" {
			archFound = true
		}
	}
	if !origFound || !archFound {
		t.Error("expected both original and archived key")
	}
}

func TestArchive_MissingKey_ReturnsError(t *testing.T) {
	entries := makeArchiveEntries()
	opts := DefaultArchiveOptions()
	_, err := Archive(entries, []string{"MISSING_KEY"}, opts)
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if !strings.Contains(err.Error(), "MISSING_KEY") {
		t.Errorf("error should mention key, got: %v", err)
	}
}

func TestArchive_SetsComment(t *testing.T) {
	entries := makeArchiveEntries()
	opts := DefaultArchiveOptions()
	out, err := Archive(entries, []string{"APP_ENV"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range out {
		if e.Key == "ARCHIVED_APP_ENV" && !strings.Contains(e.Comment, "APP_ENV") {
			t.Errorf("expected comment to reference original key, got: %s", e.Comment)
		}
	}
}
