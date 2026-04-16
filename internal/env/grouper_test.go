package env

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makeGroupEntries(keys ...string) []parser.Entry {
	var entries []parser.Entry
	for _, k := range keys {
		entries = append(entries, parser.Entry{Key: k, Value: "val"})
	}
	return entries
}

func TestGroup_ByPrefix(t *testing.T) {
	entries := makeGroupEntries("DB_HOST", "DB_PORT", "APP_NAME", "PORT")
	groups := Group(entries, DefaultGroupOptions())

	if len(groups["DB"]) != 2 {
		t.Errorf("expected 2 DB entries, got %d", len(groups["DB"]))
	}
	if len(groups["APP"]) != 1 {
		t.Errorf("expected 1 APP entry, got %d", len(groups["APP"]))
	}
	if len(groups["_"]) != 1 {
		t.Errorf("expected 1 ungrouped entry, got %d", len(groups["_"]))
	}
}

func TestGroup_EmptyEntries(t *testing.T) {
	groups := Group([]parser.Entry{}, DefaultGroupOptions())
	if len(groups) != 0 {
		t.Errorf("expected empty groups, got %d", len(groups))
	}
}

func TestGroupNames_Sorted(t *testing.T) {
	entries := makeGroupEntries("Z_KEY", "A_KEY", "M_KEY")
	groups := Group(entries, DefaultGroupOptions())
	names := GroupNames(groups, true)

	if names[0] != "A" || names[1] != "M" || names[2] != "Z" {
		t.Errorf("expected sorted names A,M,Z got %v", names)
	}
}

func TestGroupNames_Unsorted(t *testing.T) {
	entries := makeGroupEntries("DB_HOST", "APP_NAME")
	groups := Group(entries, DefaultGroupOptions())
	names := GroupNames(groups, false)
	if len(names) != 2 {
		t.Errorf("expected 2 group names, got %d", len(names))
	}
}

func TestGroup_NoUnderscore(t *testing.T) {
	entries := makeGroupEntries("PORT", "HOST", "DEBUG")
	groups := Group(entries, DefaultGroupOptions())
	if len(groups["_"]) != 3 {
		t.Errorf("expected 3 ungrouped entries, got %d", len(groups["_"]))
	}
}
