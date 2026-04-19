package env

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeInspectEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "APP_NAME", Value: "myapp", Comment: ""},
		{Key: "DB_PASS", Value: `"secret123"`, Comment: "# sensitive"},
		{Key: "EMPTY_KEY", Value: "", Comment: ""},
		{Key: "SINGLE_QUOTED", Value: "'hello'", Comment: ""},
	}
}

func TestInspect_BasicKey(t *testing.T) {
	entries := makeInspectEntries()
	res, err := Inspect(entries, "APP_NAME", DefaultInspectOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Key != "APP_NAME" || res.Value != "myapp" {
		t.Errorf("unexpected result: %+v", res)
	}
	if res.IsEmpty || res.IsQuoted {
		t.Errorf("expected non-empty, non-quoted")
	}
	if res.LineNumber != 1 {
		t.Errorf("expected line 1, got %d", res.LineNumber)
	}
}

func TestInspect_QuotedValue(t *testing.T) {
	entries := makeInspectEntries()
	res, err := Inspect(entries, "DB_PASS", DefaultInspectOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsQuoted || res.QuoteChar != `"` {
		t.Errorf("expected double-quoted, got %+v", res)
	}
	if !res.HasComment {
		t.Errorf("expected comment")
	}
}

func TestInspect_EmptyValue(t *testing.T) {
	entries := makeInspectEntries()
	res, err := Inspect(entries, "EMPTY_KEY", DefaultInspectOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsEmpty {
		t.Errorf("expected empty")
	}
}

func TestInspect_SingleQuoted(t *testing.T) {
	entries := makeInspectEntries()
	res, err := Inspect(entries, "SINGLE_QUOTED", DefaultInspectOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsQuoted || res.QuoteChar != "'" {
		t.Errorf("expected single-quoted")
	}
}

func TestInspect_MissingKey_Error(t *testing.T) {
	entries := makeInspectEntries()
	_, err := Inspect(entries, "MISSING", DefaultInspectOptions())
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestInspect_MissingKey_NoError(t *testing.T) {
	entries := makeInspectEntries()
	opts := DefaultInspectOptions()
	opts.ErrorOnMissing = false
	res, err := Inspect(entries, "MISSING", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != nil {
		t.Errorf("expected nil result")
	}
}

func TestFormatInspect_ContainsKey(t *testing.T) {
	entries := makeInspectEntries()
	res, _ := Inspect(entries, "APP_NAME", DefaultInspectOptions())
	out := FormatInspect(res)
	if !strings.Contains(out, "APP_NAME") {
		t.Errorf("expected key in output")
	}
	if !strings.Contains(out, "myapp") {
		t.Errorf("expected value in output")
	}
}
