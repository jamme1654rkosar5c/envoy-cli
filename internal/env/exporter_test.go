package env

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeExporterEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "DEBUG", Value: "true"},
		{Key: "DB_PASSWORD", Value: "s3cr3t"},
	}
}

func TestExportEntries_DotEnvFormat(t *testing.T) {
	entries := makeExporterEntries()
	out, err := ExportEntries(entries, DefaultExportOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_NAME=envoy") {
		t.Errorf("expected APP_NAME=envoy in output, got:\n%s", out)
	}
	if !strings.Contains(out, "DEBUG=true") {
		t.Errorf("expected DEBUG=true in output, got:\n%s", out)
	}
}

func TestExportEntries_ShellFormat(t *testing.T) {
	entries := makeExporterEntries()
	opts := DefaultExportOptions()
	opts.Format = ExportFormatShell
	out, err := ExportEntries(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export APP_NAME=envoy") {
		t.Errorf("expected export prefix, got:\n%s", out)
	}
}

func TestExportEntries_InlineFormat(t *testing.T) {
	entries := makeExporterEntries()
	opts := DefaultExportOptions()
	opts.Format = ExportFormatInline
	out, err := ExportEntries(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_NAME=envoy") {
		t.Errorf("expected inline entry, got:\n%s", out)
	}
	// inline entries are space-separated, not newline-separated
	if strings.Contains(out, "\n") {
		t.Errorf("inline format should not contain newlines, got:\n%s", out)
	}
}

func TestExportEntries_UnsupportedFormat(t *testing.T) {
	opts := DefaultExportOptions()
	opts.Format = "xml"
	_, err := ExportEntries(makeExporterEntries(), opts)
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExportEntries_IncludeKeys(t *testing.T) {
	opts := DefaultExportOptions()
	opts.IncludeKeys = []string{"APP_NAME"}
	out, err := ExportEntries(makeExporterEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "DEBUG") {
		t.Errorf("expected DEBUG to be excluded, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_NAME") {
		t.Errorf("expected APP_NAME to be present, got:\n%s", out)
	}
}

func TestExportEntries_ExcludeKeys(t *testing.T) {
	opts := DefaultExportOptions()
	opts.ExcludeKeys = []string{"DB_PASSWORD"}
	out, err := ExportEntries(makeExporterEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "DB_PASSWORD") {
		t.Errorf("expected DB_PASSWORD to be excluded, got:\n%s", out)
	}
}

func TestExportEntries_QuoteValues(t *testing.T) {
	opts := DefaultExportOptions()
	opts.QuoteValues = true
	out, err := ExportEntries(makeExporterEntries(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `APP_NAME="envoy"`) {
		t.Errorf("expected quoted value, got:\n%s", out)
	}
}
