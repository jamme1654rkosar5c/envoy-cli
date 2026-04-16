package env

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makeEntries(keys ...string) []parser.Entry {
	out := make([]parser.Entry, len(keys))
	for i, k := range keys {
		out[i] = parser.Entry{Key: k, Value: "v"}
	}
	return out
}

func keys(entries []parser.Entry) []string {
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.Key
	}
	return out
}

func TestSort_Alpha(t *testing.T) {
	entries := makeEntries("ZEBRA", "APPLE", "MANGO")
	sorted := Sort(entries, DefaultSortOptions())
	got := keys(sorted)
	want := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range want {
		if got[i] != k {
			t.Errorf("pos %d: got %s want %s", i, got[i], k)
		}
	}
}

func TestSort_AlphaReverse(t *testing.T) {
	entries := makeEntries("ZEBRA", "APPLE", "MANGO")
	sorted := Sort(entries, SortOptions{Order: SortAlphaReverse})
	got := keys(sorted)
	want := []string{"ZEBRA", "MANGO", "APPLE"}
	for i, k := range want {
		if got[i] != k {
			t.Errorf("pos %d: got %s want %s", i, got[i], k)
		}
	}
}

func TestSort_ByLength(t *testing.T) {
	entries := makeEntries("AB", "ABCDE", "ABC")
	sorted := Sort(entries, SortOptions{Order: SortByLength})
	got := keys(sorted)
	want := []string{"AB", "ABC", "ABCDE"}
	for i, k := range want {
		if got[i] != k {
			t.Errorf("pos %d: got %s want %s", i, got[i], k)
		}
	}
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	entries := makeEntries("ZEBRA", "APPLE")
	origFirst := entries[0].Key
	Sort(entries, DefaultSortOptions())
	if entries[0].Key != origFirst {
		t.Errorf("original slice was mutated")
	}
}
