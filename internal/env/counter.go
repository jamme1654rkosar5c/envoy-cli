package env

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// CountOptions controls what the Count function tallies.
type CountOptions struct {
	// IncludeEmpty counts entries with empty values.
	IncludeEmpty bool
	// IncludeCommented counts entries that have inline comments.
	IncludeCommented bool
	// PrefixBreakdown returns per-prefix counts when true.
	PrefixBreakdown bool
	// PrefixSep is the separator used to detect prefixes (default "_").
	PrefixSep string
}

// DefaultCountOptions returns sensible defaults.
func DefaultCountOptions() CountOptions {
	return CountOptions{
		IncludeEmpty:     true,
		IncludeCommented: true,
		PrefixBreakdown:  false,
		PrefixSep:        "_",
	}
}

// CountResult holds the tallied statistics for a set of env entries.
type CountResult struct {
	Total      int
	Empty      int
	Commented  int
	Unique     int
	Duplicates int
	// Prefixes maps each detected prefix to its entry count.
	// Populated only when CountOptions.PrefixBreakdown is true.
	Prefixes map[string]int
}

// Count tallies statistics about the provided entries.
func Count(entries []parser.Entry, opts CountOptions) CountResult {
	if opts.PrefixSep == "" {
		opts.PrefixSep = "_"
	}

	seen := make(map[string]int)
	result := CountResult{
		Prefixes: make(map[string]int),
	}

	for _, e := range entries {
		result.Total++
		seen[e.Key]++

		if e.Value == "" {
			result.Empty++
		}
		if e.Comment != "" {
			result.Commented++
		}

		if opts.PrefixBreakdown {
			prefix := extractPrefix(e.Key, opts.PrefixSep)
			result.Prefixes[prefix]++
		}
	}

	for _, count := range seen {
		if count == 1 {
			result.Unique++
		} else {
			result.Duplicates += count
		}
	}

	return result
}

// FormatCount returns a human-readable summary of a CountResult.
func FormatCount(r CountResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Total:      %d\n", r.Total))
	sb.WriteString(fmt.Sprintf("Unique:     %d\n", r.Unique))
	sb.WriteString(fmt.Sprintf("Duplicates: %d\n", r.Duplicates))
	sb.WriteString(fmt.Sprintf("Empty:      %d\n", r.Empty))
	sb.WriteString(fmt.Sprintf("Commented:  %d\n", r.Commented))
	if len(r.Prefixes) > 0 {
		sb.WriteString("Prefixes:\n")
		for p, c := range r.Prefixes {
			sb.WriteString(fmt.Sprintf("  %-20s %d\n", p, c))
		}
	}
	return sb.String()
}

func extractPrefix(key, sep string) string {
	idx := strings.Index(key, sep)
	if idx <= 0 {
		return key
	}
	return key[:idx]
}
