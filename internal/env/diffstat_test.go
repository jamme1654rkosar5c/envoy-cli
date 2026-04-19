package env

import (
	"strings"
	"testing"
)

func makeDiffEntries(pairs ...string) []EnvEntry {
	var entries []EnvEntry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, EnvEntry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestDiffStat_NoChanges(t *testing.T) {
	base := makeDiffEntries("A", "1", "B", "2")
	target := makeDiffEntries("A", "1", "B", "2")
	results := DiffStat(base, target, DefaultDiffStatOptions())
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestDiffStat_AddedKey(t *testing.T) {
	base := makeDiffEntries("A", "1")
	target := makeDiffEntries("A", "1", "B", "2")
	results := DiffStat(base, target, DefaultDiffStatOptions())
	if len(results) != 1 || results[0].Status != "added" || results[0].Key != "B" {
		t.Fatalf("expected added B, got %+v", results)
	}
}

func TestDiffStat_RemovedKey(t *testing.T) {
	base := makeDiffEntries("A", "1", "B", "2")
	target := makeDiffEntries("A", "1")
	results := DiffStat(base, target, DefaultDiffStatOptions())
	if len(results) != 1 || results[0].Status != "removed" || results[0].Key != "B" {
		t.Fatalf("expected removed B, got %+v", results)
	}
}

func TestDiffStat_ChangedKey(t *testing.T) {
	base := makeDiffEntries("A", "old")
	target := makeDiffEntries("A", "new")
	results := DiffStat(base, target, DefaultDiffStatOptions())
	if len(results) != 1 || results[0].Status != "changed" {
		t.Fatalf("expected changed A, got %+v", results)
	}
	if results[0].OldVal != "old" || results[0].NewVal != "new" {
		t.Errorf("unexpected values: %+v", results[0])
	}
}

func TestDiffStat_IncludeUnchanged(t *testing.T) {
	base := makeDiffEntries("A", "1")
	target := makeDiffEntries("A", "1")
	opts := DefaultDiffStatOptions()
	opts.IncludeUnchanged = true
	results := DiffStat(base, target, opts)
	if len(results) != 1 || results[0].Status != "unchanged" {
		t.Fatalf("expected unchanged entry, got %+v", results)
	}
}

func TestDiffStat_RedactValues(t *testing.T) {
	base := makeDiffEntries("SECRET", "hunter2")
	target := makeDiffEntries("SECRET", "newpass")
	opts := DefaultDiffStatOptions()
	opts.RedactValues = true
	results := DiffStat(base, target, opts)
	if len(results) != 1 {
		t.Fatalf("expected 1 result")
	}
	if results[0].OldVal != "***" || results[0].NewVal != "***" {
		t.Errorf("expected redacted values, got %+v", results[0])
	}
}

func TestFormatDiffStat_ContainsHeaders(t *testing.T) {
	out := FormatDiffStat([]DiffStatEntry{{Key: "X", Status: "added", NewVal: "y"}})
	if !strings.Contains(out, "KEY") {
		t.Errorf("expected header in output: %s", out)
	}
}

func TestFormatDiffStat_EmptyReturnsMessage(t *testing.T) {
	out := FormatDiffStat(nil)
	if !strings.Contains(out, "No differences") {
		t.Errorf("expected no-diff message, got: %s", out)
	}
}
