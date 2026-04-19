package env

import (
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeLabelEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "APP_NAME", Value: "envoy", Comment: ""},
		{Key: "APP_ENV", Value: "production", Comment: "some note"},
		{Key: "DB_URL", Value: "postgres://localhost/db", Comment: ""},
	}
}

func TestLabel_SetsLabel(t *testing.T) {
	entries := makeLabelEntries()
	out, err := Label(entries, "APP_NAME", "core", DefaultLabelOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := GetLabel(out, "APP_NAME")
	if got != "core" {
		t.Errorf("expected label 'core', got %q", got)
	}
}

func TestLabel_PreservesExistingComment(t *testing.T) {
	entries := makeLabelEntries()
	out, err := Label(entries, "APP_ENV", "infra", DefaultLabelOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry := out[1]
	if entry.Comment == "" || entry.Comment == "label:infra" {
		t.Errorf("expected original comment preserved, got %q", entry.Comment)
	}
	if GetLabel(out, "APP_ENV") != "infra" {
		t.Error("label not set")
	}
}

func TestLabel_NoOverwrite_ReturnsError(t *testing.T) {
	entries := makeLabelEntries()
	out, _ := Label(entries, "APP_NAME", "core", DefaultLabelOptions())
	_, err := Label(out, "APP_NAME", "other", DefaultLabelOptions())
	if err == nil {
		t.Error("expected error when overwrite disabled")
	}
}

func TestLabel_Overwrite_ReplacesLabel(t *testing.T) {
	entries := makeLabelEntries()
	out, _ := Label(entries, "APP_NAME", "core", DefaultLabelOptions())
	opts := DefaultLabelOptions()
	opts.Overwrite = true
	out, err := Label(out, "APP_NAME", "updated", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if GetLabel(out, "APP_NAME") != "updated" {
		t.Error("label not updated")
	}
}

func TestLabel_MissingKey_ReturnsError(t *testing.T) {
	entries := makeLabelEntries()
	_, err := Label(entries, "MISSING", "x", DefaultLabelOptions())
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestLabel_MissingKey_AllowMissing_NoError(t *testing.T) {
	entries := makeLabelEntries()
	opts := DefaultLabelOptions()
	opts.AllowMissing = true
	_, err := Label(entries, "MISSING", "x", opts)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUnlabel_ClearsLabel(t *testing.T) {
	entries := makeLabelEntries()
	out, _ := Label(entries, "APP_NAME", "core", DefaultLabelOptions())
	out = Unlabel(out, "APP_NAME")
	if GetLabel(out, "APP_NAME") != "" {
		t.Error("expected label to be cleared")
	}
}

func TestGetLabel_NoLabel_ReturnsEmpty(t *testing.T) {
	entries := makeLabelEntries()
	if got := GetLabel(entries, "APP_NAME"); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}
