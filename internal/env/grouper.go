package env

import (
	"sort"
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// GroupOptions controls how entries are grouped.
type GroupOptions struct {
	SortGroupNames bool
}

// DefaultGroupOptions returns sensible defaults.
func DefaultGroupOptions() GroupOptions {
	return GroupOptions{
		SortGroupNames: true,
	}
}

// Group partitions env entries by their prefix (the part before the first '_').
// Entries with no underscore are placed under the "_" group.
func Group(entries []parser.Entry, opts GroupOptions) map[string][]parser.Entry {
	result := make(map[string][]parser.Entry)
	for _, e := range entries {
		key := groupKey(e.Key)
		result[key] = append(result[key], e)
	}
	return result
}

// GroupNames returns the sorted (or unsorted) group names from a grouped map.
func GroupNames(groups map[string][]parser.Entry, sorted bool) []string {
	names := make([]string, 0, len(groups))
	for k := range groups {
		names = append(names, k)
	}
	if sorted {
		sort.Strings(names)
	}
	return names
}

func groupKey(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return strings.ToUpper(key[:idx])
	}
	return "_"
}
