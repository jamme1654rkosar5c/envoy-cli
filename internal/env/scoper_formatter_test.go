package env

import (
	"strings"
	"testing"
)

func TestBuildScopeSummaries_SameKey(t *testing.T) {
	original := []Entry{
		{Key: "APP_HOST", Value: "localhost"},
	}
	scoped := []Entry{
		{Key: "APP_HOST", Value: "localhost"},
	}

	summaries := BuildScopeSummaries(original, scoped)

	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}
	if summaries[0].ScopedKey != "APP_HOST" {
		t.Errorf("expected ScopedKey APP_HOST, got %s", summaries[0].ScopedKey)
	}
}

func TestBuildScopeSummaries_DetectsRenamedKey(t *testing.T) {
	original := []Entry{
		{Key: "APP_HOST", Value: "localhost"},
	}
	scoped := []Entry{
		{Key: "HOST", Value: "localhost"},
	}

	summaries := BuildScopeSummaries(original, scoped)

	if summaries[0].OriginalKey != "APP_HOST" {
		t.Errorf("expected OriginalKey APP_HOST, got %s", summaries[0].OriginalKey)
	}
	if summaries[0].ScopedKey != "HOST" {
		t.Errorf("expected ScopedKey HOST, got %s", summaries[0].ScopedKey)
	}
}

func TestFormatScope_ContainsHeaders(t *testing.T) {
	summaries := []ScopeSummary{
		{OriginalKey: "APP_HOST", ScopedKey: "HOST", Value: "localhost"},
	}

	out := FormatScope(summaries)

	if !strings.Contains(out, "ORIGINAL KEY") {
		t.Error("expected ORIGINAL KEY header")
	}
	if !strings.Contains(out, "SCOPED KEY") {
		t.Error("expected SCOPED KEY header")
	}
}

func TestFormatScope_EmptySummaries(t *testing.T) {
	out := FormatScope([]ScopeSummary{})

	if !strings.Contains(out, "No scoped entries") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatScope_ShowsKeyAndValue(t *testing.T) {
	summaries := []ScopeSummary{
		{OriginalKey: "DB_HOST", ScopedKey: "DB_HOST", Value: "db.local"},
	}

	out := FormatScope(summaries)

	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", out)
	}
	if !strings.Contains(out, "db.local") {
		t.Errorf("expected db.local in output, got: %s", out)
	}
}
