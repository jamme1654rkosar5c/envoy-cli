package env

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeAliasFormatterEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DATABASE_HOST", Value: "localhost", Comment: "alias of DB_HOST"},
		{Key: "APP_ENV", Value: "staging"},
		{Key: "ENVIRONMENT", Value: "staging", Comment: "alias of APP_ENV"},
	}
}

func TestListAliases_ReturnsAliasedEntries(t *testing.T) {
	entries := makeAliasFormatterEntries()
	summaries := ListAliases(entries)
	if len(summaries) != 2 {
		t.Fatalf("expected 2 aliases, got %d", len(summaries))
	}
}

func TestListAliases_NoAliases_ReturnsEmpty(t *testing.T) {
	entries := []parser.Entry{
		{Key: "FOO", Value: "bar"},
	}
	summaries := ListAliases(entries)
	if len(summaries) != 0 {
		t.Errorf("expected 0 aliases, got %d", len(summaries))
	}
}

func TestListAliases_CorrectSourceAndAlias(t *testing.T) {
	entries := makeAliasFormatterEntries()
	summaries := ListAliases(entries)
	if summaries[0].Source != "DB_HOST" {
		t.Errorf("expected source DB_HOST, got %q", summaries[0].Source)
	}
	if summaries[0].Alias != "DATABASE_HOST" {
		t.Errorf("expected alias DATABASE_HOST, got %q", summaries[0].Alias)
	}
}

func TestFormatAliases_ContainsHeaders(t *testing.T) {
	summaries := ListAliases(makeAliasFormatterEntries())
	out := FormatAliases(summaries)
	if !strings.Contains(out, "ALIAS") || !strings.Contains(out, "SOURCE") {
		t.Error("expected table headers ALIAS and SOURCE")
	}
}

func TestFormatAliases_EmptySummaries(t *testing.T) {
	out := FormatAliases(nil)
	if !strings.Contains(out, "No aliases") {
		t.Errorf("expected empty message, got: %q", out)
	}
}
