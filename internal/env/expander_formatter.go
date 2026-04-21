package env

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// ExpandSummary describes a single expansion result for display purposes.
type ExpandSummary struct {
	Key      string
	Original string
	Expanded string
	Changed  bool
}

// BuildExpandSummaries compares original and expanded entry slices and returns
// a summary for each entry that changed during expansion.
func BuildExpandSummaries(original, expanded []parser.Entry) []ExpandSummary {
	origMap := make(map[string]string, len(original))
	for _, e := range original {
		origMap[e.Key] = e.Value
	}

	var summaries []ExpandSummary
	for _, e := range expanded {
		orig := origMap[e.Key]
		summaries = append(summaries, ExpandSummary{
			Key:      e.Key,
			Original: orig,
			Expanded: e.Value,
			Changed:  orig != e.Value,
		})
	}
	return summaries
}

// FormatExpand renders the expansion summaries as a human-readable table.
// Only changed entries are shown unless showAll is true.
func FormatExpand(summaries []ExpandSummary, showAll bool) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %-40s %s\n", "KEY", "ORIGINAL", "EXPANDED"))
	sb.WriteString(strings.Repeat("-", 90) + "\n")

	printed := 0
	for _, s := range summaries {
		if !s.Changed && !showAll {
			continue
		}
		sb.WriteString(fmt.Sprintf("%-30s %-40s %s\n",
			truncateExp(s.Key, 29),
			truncateExp(s.Original, 39),
			s.Expanded,
		))
		printed++
	}

	if printed == 0 {
		sb.WriteString("(no expansions performed)\n")
	}
	return sb.String()
}

func truncateExp(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
