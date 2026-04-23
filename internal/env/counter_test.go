package env

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeCountEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "DB_HOST", Value: "localhost", Comment: "primary db"},
		{Key: "DB_PORT", Value: "5432"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "APP_DEBUG", Value: ""},
		{Key: "SECRET_KEY", Value: "abc123", Comment: "rotate often"},
		{Key: "DB_HOST", Value: "replica"},
	}
}

func TestCount_Total(t *testing.T) {
	entries := makeCountEntries()
	opts := DefaultCountOptions()
	r := Count(entries, opts)
	if r.Total != 6 {
		t.Errorf("expected Total=6, got %d", r.Total)
	}
}

func TestCount_EmptyValues(t *testing.T) {
	entries := makeCountEntries()
	opts := DefaultCountOptions()
	r := Count(entries, opts)
	if r.Empty != 1 {
		t.Errorf("expected Empty=1, got %d", r.Empty)
	}
}

func TestCount_CommentedEntries(t *testing.T) {
	entries := makeCountEntries()
	opts := DefaultCountOptions()
	r := Count(entries, opts)
	if r.Commented != 2 {
		t.Errorf("expected Commented=2, got %d", r.Commented)
	}
}

func TestCount_DuplicateKeys(t *testing.T) {
	entries := makeCountEntries()
	opts := DefaultCountOptions()
	r := Count(entries, opts)
	// DB_HOST appears twice → 2 duplicates; rest are unique
	if r.Duplicates != 2 {
		t.Errorf("expected Duplicates=2, got %d", r.Duplicates)
	}
	if r.Unique != 4 {
		t.Errorf("expected Unique=4, got %d", r.Unique)
	}
}

func TestCount_PrefixBreakdown(t *testing.T) {
	entries := makeCountEntries()
	opts := DefaultCountOptions()
	opts.PrefixBreakdown = true
	r := Count(entries, opts)

	if r.Prefixes["DB"] != 3 {
		t.Errorf("expected DB prefix count=3, got %d", r.Prefixes["DB"])
	}
	if r.Prefixes["APP"] != 2 {
		t.Errorf("expected APP prefix count=2, got %d", r.Prefixes["APP"])
	}
	if r.Prefixes["SECRET"] != 1 {
		t.Errorf("expected SECRET prefix count=1, got %d", r.Prefixes["SECRET"])
	}
}

func TestCount_NoPrefixBreakdown_EmptyMap(t *testing.T) {
	entries := makeCountEntries()
	opts := DefaultCountOptions()
	opts.PrefixBreakdown = false
	r := Count(entries, opts)
	if len(r.Prefixes) != 0 {
		t.Errorf("expected empty Prefixes map, got %v", r.Prefixes)
	}
}

func TestFormatCount_ContainsLabels(t *testing.T) {
	r := CountResult{Total: 5, Unique: 4, Duplicates: 1, Empty: 2, Commented: 3}
	out := FormatCount(r)
	for _, label := range []string{"Total", "Unique", "Duplicates", "Empty", "Commented"} {
		if !strings.Contains(out, label) {
			t.Errorf("expected output to contain %q", label)
		}
	}
}
