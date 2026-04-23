package env

import (
	"fmt"
	"strings"
)

// ClassifySummary holds display data for a single classified entry.
type ClassifySummary struct {
	Key      string
	Category Category
	Value    string
}

// BuildClassifySummaries converts ClassifiedEntry slice to summaries.
func BuildClassifySummaries(entries []ClassifiedEntry) []ClassifySummary {
	out := make([]ClassifySummary, 0, len(entries))
	for _, ce := range entries {
		out = append(out, ClassifySummary{
			Key:      ce.Entry.Key,
			Category: ce.Category,
			Value:    truncateClassVal(ce.Entry.Value, 32),
		})
	}
	return out
}

// FormatClassify returns a human-readable table of classified entries.
func FormatClassify(summaries []ClassifySummary) string {
	if len(summaries) == 0 {
		return "no entries to classify\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-40s %-12s %s\n", "KEY", "CATEGORY", "VALUE"))
	sb.WriteString(strings.Repeat("-", 70) + "\n")
	for _, s := range summaries {
		sb.WriteString(fmt.Sprintf("%-40s %-12s %s\n", s.Key, s.Category, s.Value))
	}
	return sb.String()
}

// GroupByCategory returns a map of category → summaries.
func GroupByCategory(summaries []ClassifySummary) map[Category][]ClassifySummary {
	result := make(map[Category][]ClassifySummary)
	for _, s := range summaries {
		result[s.Category] = append(result[s.Category], s)
	}
	return result
}

func truncateClassVal(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
