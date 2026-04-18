package env

import (
	"testing"

	"github.com/yourorg/envoy-cli/internal/parser"
)

func makeSplitEntries(pairs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestSplit_EvenDistribution(t *testing.T) {
	entries := makeSplitEntries("A", "1", "B", "2", "C", "3", "D", "4")
	buckets := Split(entries, 2, DefaultSplitOptions())
	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(buckets))
	}
	if len(buckets[0]) != 2 || len(buckets[1]) != 2 {
		t.Errorf("expected 2 entries per bucket")
	}
}

func TestSplit_MoreBucketsThanEntries(t *testing.T) {
	entries := makeSplitEntries("A", "1", "B", "2")
	buckets := Split(entries, 5, DefaultSplitOptions())
	if len(buckets) != 5 {
		t.Fatalf("expected 5 buckets, got %d", len(buckets))
	}
	total := 0
	for _, b := range buckets {
		total += len(b)
	}
	if total != 2 {
		t.Errorf("expected 2 total entries, got %d", total)
	}
}

func TestSplit_SkipEmpty(t *testing.T) {
	entries := makeSplitEntries("A", "1", "B", "", "C", "3")
	opts := DefaultSplitOptions()
	opts.SkipEmpty = true
	buckets := Split(entries, 1, opts)
	if len(buckets[0]) != 2 {
		t.Errorf("expected 2 non-empty entries, got %d", len(buckets[0]))
	}
}

func TestSplit_MaxBucketsCaps(t *testing.T) {
	entries := makeSplitEntries("A", "1", "B", "2", "C", "3")
	opts := DefaultSplitOptions()
	opts.MaxBuckets = 2
	buckets := Split(entries, 10, opts)
	if len(buckets) != 2 {
		t.Errorf("expected max 2 buckets, got %d", len(buckets))
	}
}

func TestSplit_EmptyEntries(t *testing.T) {
	buckets := Split([]parser.Entry{}, 3, DefaultSplitOptions())
	if len(buckets) != 0 {
		t.Errorf("expected empty result for empty input")
	}
}
