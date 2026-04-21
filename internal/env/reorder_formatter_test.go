package env

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeReorderFormatterEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "ALPHA", Value: "a"},
		{Key: "BETA", Value: "b"},
		{Key: "GAMMA", Value: "c"},
	}
}

func TestBuildReorderSummaries_NoChange(t *testing.T) {
	entries := makeReorderFormatterEntries()
	summaries := BuildReorderSummaries(entries, entries)
	for _, s := range summaries {
		if s.Moved {
			t.Errorf("expected no movement for key %q", s.Key)
		}
	}
}

func TestBuildReorderSummaries_DetectsMove(t *testing.T) {
	original := makeReorderFormatterEntries()
	reordered := []parser.Entry{
		{Key: "GAMMA", Value: "c"},
		{Key: "ALPHA", Value: "a"},
		{Key: "BETA", Value: "b"},
	}
	summaries := BuildReorderSummaries(original, reordered)
	if !summaries[0].Moved {
		t.Error("GAMMA should be marked as moved")
	}
	if !summaries[1].Moved {
		t.Error("ALPHA should be marked as moved")
	}
}

func TestFormatReorder_ContainsHeaders(t *testing.T) {
	original := makeReorderFormatterEntries()
	summaries := BuildReorderSummaries(original, original)
	output := FormatReorder(summaries)
	if !strings.Contains(output, "KEY") {
		t.Error("expected header KEY in output")
	}
	if !strings.Contains(output, "MOVED") {
		t.Error("expected header MOVED in output")
	}
}

func TestFormatReorder_ShowsKeyNames(t *testing.T) {
	original := makeReorderFormatterEntries()
	summaries := BuildReorderSummaries(original, original)
	output := FormatReorder(summaries)
	for _, e := range original {
		if !strings.Contains(output, e.Key) {
			t.Errorf("expected key %q in output", e.Key)
		}
	}
}

func TestFormatReorder_EmptySummaries(t *testing.T) {
	output := FormatReorder(nil)
	if !strings.Contains(output, "No entries") {
		t.Error("expected 'No entries' message for empty summaries")
	}
}
