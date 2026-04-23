package env

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makeTraceEntries(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestTrace_NoReferences(t *testing.T) {
	entries := makeTraceEntries("FOO", "bar", "BAZ", "qux")
	results, err := Trace(entries, DefaultTraceOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r := results["FOO"]; len(r.Chain) != 1 || r.Chain[0] != "FOO" {
		t.Errorf("expected chain [FOO], got %v", r.Chain)
	}
	if results["FOO"].Depth != 0 {
		t.Errorf("expected depth 0, got %d", results["FOO"].Depth)
	}
}

func TestTrace_SingleHop(t *testing.T) {
	entries := makeTraceEntries("A", "hello", "B", "${A}")
	results, err := Trace(entries, DefaultTraceOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := results["B"]
	if len(r.Chain) != 2 {
		t.Fatalf("expected chain length 2, got %d: %v", len(r.Chain), r.Chain)
	}
	if r.Chain[0] != "B" || r.Chain[1] != "A" {
		t.Errorf("unexpected chain: %v", r.Chain)
	}
	if r.Depth != 1 {
		t.Errorf("expected depth 1, got %d", r.Depth)
	}
}

func TestTrace_MultiHop(t *testing.T) {
	entries := makeTraceEntries("X", "base", "Y", "$X", "Z", "${Y}")
	results, err := Trace(entries, DefaultTraceOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := results["Z"]
	if r.Depth != 2 {
		t.Errorf("expected depth 2, got %d", r.Depth)
	}
}

func TestTrace_CycleDetected_ReturnsError(t *testing.T) {
	entries := makeTraceEntries("A", "${B}", "B", "${A}")
	_, err := Trace(entries, DefaultTraceOptions())
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
}

func TestTrace_CycleAllowed_SetsCycleFlag(t *testing.T) {
	entries := makeTraceEntries("A", "${B}", "B", "${A}")
	opts := DefaultTraceOptions()
	opts.AllowCycles = true
	results, err := Trace(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results["A"].Cycles && !results["B"].Cycles {
		t.Error("expected at least one cycle flag to be true")
	}
}

func TestTrace_MaxDepth_ReturnsError(t *testing.T) {
	entries := makeTraceEntries("A", "${B}", "B", "${C}", "C", "${D}", "D", "${E}", "E", "val")
	opts := DefaultTraceOptions()
	opts.MaxDepth = 2
	_, err := Trace(entries, opts)
	if err == nil {
		t.Fatal("expected max depth error, got nil")
	}
}
