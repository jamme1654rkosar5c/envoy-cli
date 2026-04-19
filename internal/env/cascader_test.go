package env

import (
	"testing"

	"github.com/your-org/envoy-cli/internal/parser"
)

func makeCascadeEntries(pairs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestCascade_BaseOnly(t *testing.T) {
	base := makeCascadeEntries("A", "1", "B", "2")
	result := Cascade([][]parser.Entry{base}, DefaultCascadeOptions())
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestCascade_OverwriteEnabled(t *testing.T) {
	base := makeCascadeEntries("A", "base")
	overlay := makeCascadeEntries("A", "overlay")
	opts := DefaultCascadeOptions()
	opts.Overwrite = true
	result := Cascade([][]parser.Entry{base, overlay}, opts)
	if result[0].Value != "overlay" {
		t.Errorf("expected 'overlay', got %q", result[0].Value)
	}
}

func TestCascade_OverwriteDisabled_KeepsBase(t *testing.T) {
	base := makeCascadeEntries("A", "base")
	overlay := makeCascadeEntries("A", "overlay")
	opts := DefaultCascadeOptions()
	opts.Overwrite = false
	result := Cascade([][]parser.Entry{base, overlay}, opts)
	if result[0].Value != "base" {
		t.Errorf("expected 'base', got %q", result[0].Value)
	}
}

func TestCascade_AddsNewKeysFromLayer(t *testing.T) {
	base := makeCascadeEntries("A", "1")
	overlay := makeCascadeEntries("B", "2")
	result := Cascade([][]parser.Entry{base, overlay}, DefaultCascadeOptions())
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestCascade_SkipEmpty_IgnoresBlankValues(t *testing.T) {
	base := makeCascadeEntries("A", "original")
	overlay := makeCascadeEntries("A", "")
	opts := DefaultCascadeOptions()
	opts.SkipEmpty = true
	result := Cascade([][]parser.Entry{base, overlay}, opts)
	if result[0].Value != "original" {
		t.Errorf("expected 'original', got %q", result[0].Value)
	}
}

func TestCascade_EmptyLayers_ReturnsNil(t *testing.T) {
	result := Cascade(nil, DefaultCascadeOptions())
	if result != nil {
		t.Errorf("expected nil result for empty layers")
	}
}
