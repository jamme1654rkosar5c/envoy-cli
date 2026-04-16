package env

import (
	"strings"
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func TestListTags_ReturnsTaggedEntries(t *testing.T) {
	opts := DefaultTagOptions()
	entries := []parser.Entry{
		{Key: "APP_ENV", Value: "production", Comment: "@tag:stable"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "API_KEY", Value: "secret", Comment: "@tag:sensitive"},
	}
	summaries := ListTags(entries, opts)
	if len(summaries) != 2 {
		t.Fatalf("expected 2 summaries, got %d", len(summaries))
	}
	if summaries[0].Tag != "stable" {
		t.Errorf("expected 'stable', got %q", summaries[0].Tag)
	}
	if summaries[1].Tag != "sensitive" {
		t.Errorf("expected 'sensitive', got %q", summaries[1].Tag)
	}
}

func TestListTags_NoTags_ReturnsEmpty(t *testing.T) {
	opts := DefaultTagOptions()
	entries := []parser.Entry{
		{Key: "APP_ENV", Value: "production"},
	}
	summaries := ListTags(entries, opts)
	if len(summaries) != 0 {
		t.Errorf("expected empty, got %d", len(summaries))
	}
}

func TestFormatTags_ContainsHeaders(t *testing.T) {
	summaries := []TagSummary{
		{Key: "APP_ENV", Value: "production", Tag: "stable"},
	}
	out := FormatTags(summaries)
	if !strings.Contains(out, "KEY") || !strings.Contains(out, "TAG") {
		t.Error("expected header row in output")
	}
	if !strings.Contains(out, "stable") {
		t.Error("expected tag value in output")
	}
}

func TestFormatTags_EmptySummaries(t *testing.T) {
	out := FormatTags(nil)
	if !strings.Contains(out, "No tagged") {
		t.Errorf("expected empty message, got %q", out)
	}
}
