package env

import (
	"testing"

	"github.com/your-org/envoy-cli/internal/parser"
)

func makePruneEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DEBUG", Value: ""},
		{Key: "SECRET_KEY", Value: "abc123"},
		{Key: "LEGACY_FLAG", Value: "true"},
		{Key: "EMPTY_AGAIN", Value: ""},
	}
}

func TestPrune_RemoveEmpty(t *testing.T) {
	opts := DefaultPruneOptions()
	kept, pruned := Prune(makePruneEntries(), opts)

	if len(kept) != 3 {
		t.Fatalf("expected 3 kept entries, got %d", len(kept))
	}
	if len(pruned) != 2 {
		t.Fatalf("expected 2 pruned keys, got %d", len(pruned))
	}
	for _, e := range kept {
		if e.Value == "" {
			t.Errorf("kept entry %q has empty value", e.Key)
		}
	}
}

func TestPrune_ExplicitKeys(t *testing.T) {
	opts := DefaultPruneOptions()
	opts.RemoveEmpty = false
	opts.Keys = []string{"LEGACY_FLAG", "DEBUG"}

	kept, pruned := Prune(makePruneEntries(), opts)

	if len(pruned) != 2 {
		t.Fatalf("expected 2 pruned, got %d", len(pruned))
	}
	for _, e := range kept {
		if e.Key == "LEGACY_FLAG" || e.Key == "DEBUG" {
			t.Errorf("expected %q to be pruned", e.Key)
		}
	}
}

func TestPrune_DryRun_DoesNotMutateOriginal(t *testing.T) {
	src := makePruneEntries()
	opts := DefaultPruneOptions()
	opts.DryRun = true

	kept, _ := Prune(src, opts)

	// Mutate the returned slice.
	if len(kept) > 0 {
		kept[0].Value = "MUTATED"
	}

	for _, e := range src {
		if e.Value == "MUTATED" {
			t.Error("DryRun mutated the original entries")
		}
	}
}

func TestPrune_RemoveCommented(t *testing.T) {
	entries := []parser.Entry{
		{Key: "A", Value: "hello"},
		{Key: "B", Value: "# deprecated"},
		{Key: "C", Value: "world"},
	}
	opts := DefaultPruneOptions()
	opts.RemoveEmpty = false
	opts.RemoveCommented = true

	kept, pruned := Prune(entries, opts)

	if len(kept) != 2 {
		t.Fatalf("expected 2 kept, got %d", len(kept))
	}
	if len(pruned) != 1 || pruned[0] != "B" {
		t.Errorf("expected B to be pruned, got %v", pruned)
	}
}

func TestPrune_NoOptions_KeepsAll(t *testing.T) {
	opts := PruneOptions{}
	kept, pruned := Prune(makePruneEntries(), opts)

	if len(kept) != len(makePruneEntries()) {
		t.Errorf("expected all entries kept, got %d", len(kept))
	}
	if len(pruned) != 0 {
		t.Errorf("expected no pruned keys, got %v", pruned)
	}
}
