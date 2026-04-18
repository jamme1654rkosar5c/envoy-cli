package env

import (
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeAnnotatorEntries(kvs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		entries = append(entries, parser.Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return entries
}

func TestAnnotate_SetsComment(t *testing.T) {
	entries := makeAnnotatorEntries("DB_HOST", "localhost")
	out, err := Annotate(entries, "DB_HOST", "primary database host", DefaultAnnotateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := GetAnnotation(out, "DB_HOST")
	if got != "primary database host" {
		t.Errorf("expected annotation text, got %q", got)
	}
}

func TestAnnotate_WithPrefix(t *testing.T) {
	entries := makeAnnotatorEntries("API_KEY", "secret")
	opts := DefaultAnnotateOptions()
	opts.Prefix = "TODO"
	out, err := Annotate(entries, "API_KEY", "rotate this key", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := GetAnnotation(out, "API_KEY")
	if got != "[TODO] rotate this key" {
		t.Errorf("unexpected annotation: %q", got)
	}
}

func TestAnnotate_NoOverwrite_KeepsOriginal(t *testing.T) {
	entries := makeAnnotatorEntries("PORT", "8080")
	opts := DefaultAnnotateOptions()
	out, _ := Annotate(entries, "PORT", "first annotation", opts)
	out, _ = Annotate(out, "PORT", "second annotation", opts)
	got := GetAnnotation(out, "PORT")
	if got != "first annotation" {
		t.Errorf("expected original annotation preserved, got %q", got)
	}
}

func TestAnnotate_Overwrite_ReplacesAnnotation(t *testing.T) {
	entries := makeAnnotatorEntries("PORT", "8080")
	opts := DefaultAnnotateOptions()
	out, _ := Annotate(entries, "PORT", "first annotation", opts)
	opts.Overwrite = true
	out, _ = Annotate(out, "PORT", "updated annotation", opts)
	got := GetAnnotation(out, "PORT")
	if got != "updated annotation" {
		t.Errorf("expected updated annotation, got %q", got)
	}
}

func TestAnnotate_KeyNotFound_ReturnsError(t *testing.T) {
	entries := makeAnnotatorEntries("FOO", "bar")
	_, err := Annotate(entries, "MISSING", "note", DefaultAnnotateOptions())
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRemoveAnnotation_ClearsComment(t *testing.T) {
	entries := makeAnnotatorEntries("HOST", "localhost")
	out, _ := Annotate(entries, "HOST", "some note", DefaultAnnotateOptions())
	out = RemoveAnnotation(out, "HOST")
	got := GetAnnotation(out, "HOST")
	if got != "" {
		t.Errorf("expected empty annotation after removal, got %q", got)
	}
}

func TestGetAnnotation_NoAnnotation_ReturnsEmpty(t *testing.T) {
	entries := makeAnnotatorEntries("KEY", "val")
	got := GetAnnotation(entries, "KEY")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}
