package env

import (
	"sort"

	"github.com/envoy-cli/internal/parser"
)

// SortOrder defines how entries should be sorted.
type SortOrder int

const (
	SortAlpha SortOrder = iota
	SortAlphaReverse
	SortByLength
)

// SortOptions configures sorting behaviour.
type SortOptions struct {
	Order     SortOrder
	StablePos bool // preserve relative order of equal keys
}

// DefaultSortOptions returns sensible defaults.
func DefaultSortOptions() SortOptions {
	return SortOptions{Order: SortAlpha}
}

// Sort returns a new slice of entries sorted according to opts.
func Sort(entries []parser.Entry, opts SortOptions) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	copy(out, entries)

	var less func(i, j int) bool
	switch opts.Order {
	case SortAlphaReverse:
		less = func(i, j int) bool { return out[i].Key > out[j].Key }
	case SortByLength:
		less = func(i, j int) bool {
			if len(out[i].Key) == len(out[j].Key) {
				return out[i].Key < out[j].Key
			}
			return len(out[i].Key) < len(out[j].Key)
		}
	default: // SortAlpha
		less = func(i, j int) bool { return out[i].Key < out[j].Key }
	}

	if opts.StablePos {
		sort.SliceStable(out, less)
	} else {
		sort.Slice(out, less)
	}
	return out
}
