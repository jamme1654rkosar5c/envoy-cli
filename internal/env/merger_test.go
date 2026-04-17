package env

import (
	"testing"

	"github.com/your-org/envoy-cli/internal/parser"
)

func makeOverlapEntries(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestOverlap_AddsNewKeys(t *testing.T) {
	dst := makeOverlapEntries("A", "1")
	src := makeOverlapEntries("B", "2")
	out := Overlap(dst, src, DefaultOverlapOptions())
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestOverlap_NoOverwrite_KeepsDstValue(t *testing.T) {
	dst := makeOverlapEntries("A", "original")
	src := makeOverlapEntries("A", "new")
	opts := DefaultOverlapOptions()
	opts.Overwrite = false
	out := Overlap(dst, src, opts)
	if out[0].Value != "original" {
		t.Errorf("expected 'original', got %q", out[0].Value)
	}
}

func TestOverlap_Overwrite_ReplacesDstValue(t *testing.T) {
	dst := makeOverlapEntries("A", "original")
	src := makeOverlapEntries("A", "new")
	opts := DefaultOverlapOptions()
	opts.Overwrite = true
	out := Overlap(dst, src, opts)
	if out[0].Value != "new" {
		t.Errorf("expected 'new', got %q", out[0].Value)
	}
}

func TestOverlap_SkipEmpty_IgnoresEmptySrc(t *testing.T) {
	dst := makeOverlapEntries("A", "keep")
	src := makeOverlapEntries("A", "", "B", "")
	opts := DefaultOverlapOptions()
	opts.Overwrite = true
	opts.SkipEmpty = true
	out := Overlap(dst, src, opts)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Value != "keep" {
		t.Errorf("expected 'keep', got %q", out[0].Value)
	}
}

func TestOverlap_DoesNotMutateDst(t *testing.T) {
	dst := makeOverlapEntries("A", "1")
	src := makeOverlapEntries("A", "2")
	opts := DefaultOverlapOptions()
	opts.Overwrite = true
	Overlap(dst, src, opts)
	if dst[0].Value != "1" {
		t.Errorf("dst was mutated, expected '1', got %q", dst[0].Value)
	}
}
