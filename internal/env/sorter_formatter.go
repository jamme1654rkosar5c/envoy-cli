package env

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// SortSummary holds metadata about a single entry after sorting.
type SortSummary struct {
	Key      string
	Value    string
	Position int
}

// BuildSortSummaries returns a slice of SortSummary for the given sorted entries.
func BuildSortSummaries(entries []parser.Entry) []SortSummary {
	summaries := make([]SortSummary, 0, len(entries))
	for i, e := range entries {
		summaries = append(summaries, SortSummary{
			Key:      e.Key,
			Value:    e.Value,
			Position: i + 1,
		})
	}
	return summaries
}

// FormatSort returns a human-readable table of sorted entries.
func FormatSort(summaries []SortSummary) string {
	if len(summaries) == 0 {
		return "No entries to display.\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-4s  %-30s  %s\n", "#", "KEY", "VALUE"))
	sb.WriteString(strings.Repeat("-", 60) + "\n")

	for _, s := range summaries {
		val := truncateSortVal(s.Value, 24)
		sb.WriteString(fmt.Sprintf("%-4d  %-30s  %s\n", s.Position, s.Key, val))
	}

	return sb.String()
}

func truncateSortVal(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
