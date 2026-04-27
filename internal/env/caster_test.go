package env

import (
	"testing"

	"github.com/envoy-cli/envoy/internal/parser"
)

func makeCastEntries(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestCast_IntValues(t *testing.T) {
	entries := makeCastEntries("PORT", "8080", "TIMEOUT", "30")
	opts := DefaultCastOptions()
	opts.TargetType = CastInt

	out, results, err := Cast(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.OK {
			t.Errorf("expected OK for key %q", r.Key)
		}
	}
	if out[0].Value != "8080" || out[1].Value != "30" {
		t.Errorf("unexpected values: %v", out)
	}
}

func TestCast_BoolValues(t *testing.T) {
	entries := makeCastEntries("DEBUG", "1", "VERBOSE", "false")
	opts := DefaultCastOptions()
	opts.TargetType = CastBool

	out, results, err := Cast(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "true" {
		t.Errorf("expected 'true' for DEBUG, got %q", out[0].Value)
	}
	if out[1].Value != "false" {
		t.Errorf("expected 'false' for VERBOSE, got %q", out[1].Value)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestCast_InvalidInt_ReturnsError(t *testing.T) {
	entries := makeCastEntries("PORT", "not-a-number")
	opts := DefaultCastOptions()
	opts.TargetType = CastInt

	_, _, err := Cast(entries, opts)
	if err == nil {
		t.Fatal("expected error for invalid int, got nil")
	}
}

func TestCast_SkipInvalid_DoesNotError(t *testing.T) {
	entries := makeCastEntries("PORT", "abc", "WORKERS", "4")
	opts := DefaultCastOptions()
	opts.TargetType = CastInt
	opts.SkipInvalid = true

	out, results, err := Cast(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// PORT should remain unchanged
	if out[0].Value != "abc" {
		t.Errorf("expected unchanged value for PORT, got %q", out[0].Value)
	}
	// WORKERS should be cast
	if out[1].Value != "4" {
		t.Errorf("expected '4' for WORKERS, got %q", out[1].Value)
	}
	var skipped int
	for _, r := range results {
		if !r.OK {
			skipped++
		}
	}
	if skipped != 1 {
		t.Errorf("expected 1 skipped result, got %d", skipped)
	}
}

func TestCast_KeyFilter_OnlyCastsMatchingKeys(t *testing.T) {
	entries := makeCastEntries("PORT", "8080", "NAME", "hello")
	opts := DefaultCastOptions()
	opts.TargetType = CastInt
	opts.Keys = []string{"PORT"}

	out, results, err := Cast(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "hello" {
		t.Errorf("NAME should be untouched, got %q", out[1].Value)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result (only PORT), got %d", len(results))
	}
}

func TestCast_FloatNormalisation(t *testing.T) {
	entries := makeCastEntries("RATE", "3.14000")
	opts := DefaultCastOptions()
	opts.TargetType = CastFloat

	out, _, err := Cast(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "3.14" {
		t.Errorf("expected '3.14', got %q", out[0].Value)
	}
}
