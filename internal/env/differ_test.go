package env

import (
	"strings"
	"testing"
)

func makeDifferEntries(pairs ...string) []EnvEntry {
	var entries []EnvEntry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, EnvEntry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestDiff_NoChanges(t *testing.T) {
	base := makeDifferEntries("HOST", "localhost", "PORT", "8080")
	next := makeDifferEntries("HOST", "localhost", "PORT", "8080")

	results, err := Diff(base, next, DefaultDiffOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Status != "unchanged" {
			t.Errorf("expected unchanged, got %q for key %q", r.Status, r.Key)
		}
	}
}

func TestDiff_AddedKey(t *testing.T) {
	base := makeDifferEntries("HOST", "localhost")
	next := makeDifferEntries("HOST", "localhost", "PORT", "9000")

	results, err := Diff(base, next, DefaultDiffOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var found bool
	for _, r := range results {
		if r.Key == "PORT" && r.Status == "added" {
			found = true
		}
	}
	if !found {
		t.Error("expected PORT to be marked as added")
	}
}

func TestDiff_RemovedKey(t *testing.T) {
	base := makeDifferEntries("HOST", "localhost", "DEBUG", "true")
	next := makeDifferEntries("HOST", "localhost")

	results, err := Diff(base, next, DefaultDiffOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var found bool
	for _, r := range results {
		if r.Key == "DEBUG" && r.Status == "removed" {
			found = true
		}
	}
	if !found {
		t.Error("expected DEBUG to be marked as removed")
	}
}

func TestDiff_ChangedKey(t *testing.T) {
	base := makeDifferEntries("PORT", "8080")
	next := makeDifferEntries("PORT", "9090")

	results, err := Diff(base, next, DefaultDiffOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Status != "changed" {
		t.Errorf("expected changed status, got %+v", results)
	}
	if results[0].OldValue != "8080" || results[0].NewValue != "9090" {
		t.Errorf("unexpected values: %+v", results[0])
	}
}

func TestDiff_MaskSecrets(t *testing.T) {
	base := makeDifferEntries("API_SECRET", "old-secret")
	next := makeDifferEntries("API_SECRET", "new-secret")

	opts := DefaultDiffOptions()
	opts.MaskSecrets = true

	results, err := Diff(base, next, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].OldValue != "***" || results[0].NewValue != "***" {
		t.Errorf("expected masked values, got old=%q new=%q", results[0].OldValue, results[0].NewValue)
	}
}

func TestDiff_IgnoreKeys(t *testing.T) {
	base := makeDifferEntries("HOST", "a", "SKIP", "x")
	next := makeDifferEntries("HOST", "b", "SKIP", "y")

	opts := DefaultDiffOptions()
	opts.IgnoreKeys = []string{"SKIP"}

	results, err := Diff(base, next, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "SKIP" {
			t.Error("SKIP key should have been ignored")
		}
	}
}

func TestFormatDiff_ContainsHeaders(t *testing.T) {
	entries := []DiffEntry{
		{Key: "FOO", Status: "added", OldValue: "", NewValue: "bar"},
	}
	out := FormatDiff(entries)
	if !strings.Contains(out, "Key") || !strings.Contains(out, "Op") {
		t.Errorf("expected headers in output, got:\n%s", out)
	}
}

func TestFormatDiff_EmptyEntries(t *testing.T) {
	out := FormatDiff(nil)
	if out != "(no differences)" {
		t.Errorf("expected no-differences message, got %q", out)
	}
}
