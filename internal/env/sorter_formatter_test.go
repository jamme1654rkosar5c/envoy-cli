package env

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeSortFormatterEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "LOG_LEVEL", Value: "info"},
	}
}

func TestBuildSortSummaries_CorrectPositions(t *testing.T) {
	entries := makeSortFormatterEntries()
	summaries := BuildSortSummaries(entries)

	if len(summaries) != 3 {
		t.Fatalf("expected 3 summaries, got %d", len(summaries))
	}
	for i, s := range summaries {
		if s.Position != i+1 {
			t.Errorf("entry %d: expected position %d, got %d", i, i+1, s.Position)
		}
	}
}

func TestBuildSortSummaries_PreservesKeyValue(t *testing.T) {
	entries := makeSortFormatterEntries()
	summaries := BuildSortSummaries(entries)

	if summaries[0].Key != "APP_NAME" {
		t.Errorf("expected key APP_NAME, got %s", summaries[0].Key)
	}
	if summaries[1].Value != "localhost" {
		t.Errorf("expected value localhost, got %s", summaries[1].Value)
	}
}

func TestFormatSort_ContainsHeaders(t *testing.T) {
	summaries := BuildSortSummaries(makeSortFormatterEntries())
	out := FormatSort(summaries)

	if !strings.Contains(out, "KEY") {
		t.Error("expected output to contain KEY header")
	}
	if !strings.Contains(out, "VALUE") {
		t.Error("expected output to contain VALUE header")
	}
}

func TestFormatSort_ContainsKeys(t *testing.T) {
	summaries := BuildSortSummaries(makeSortFormatterEntries())
	out := FormatSort(summaries)

	for _, key := range []string{"APP_NAME", "DB_HOST", "LOG_LEVEL"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected output to contain key %s", key)
		}
	}
}

func TestFormatSort_EmptySummaries(t *testing.T) {
	out := FormatSort([]SortSummary{})
	if !strings.Contains(out, "No entries") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatSort_TruncatesLongValues(t *testing.T) {
	entries := []parser.Entry{
		{Key: "LONG_VAL", Value: strings.Repeat("x", 50)},
	}
	summaries := BuildSortSummaries(entries)
	out := FormatSort(summaries)

	if !strings.Contains(out, "...") {
		t.Error("expected truncated value to contain ellipsis")
	}
}
