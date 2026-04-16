package env

import (
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// FilterOptions controls which entries are included.
type FilterOptions struct {
	Prefix   string // keep only entries whose key starts with Prefix
	Contains string // keep only entries whose key contains Contains
	NoEmpty  bool   // drop entries with empty values
}

// Filter returns a subset of entries matching all active criteria in opts.
func Filter(entries []parser.Entry, opts FilterOptions) []parser.Entry {
	var out []parser.Entry
	for _, e := range entries {
		if opts.Prefix != "" && !strings.HasPrefix(e.Key, opts.Prefix) {
			continue
		}
		if opts.Contains != "" && !strings.Contains(e.Key, opts.Contains) {
			continue
		}
		if opts.NoEmpty && strings.TrimSpace(e.Value) == "" {
			continue
		}
		out = append(out, e)
	}
	return out
}
