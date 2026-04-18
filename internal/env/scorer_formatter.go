package env

import (
	"fmt"
	"strings"
)

// FormatScores returns a human-readable table of scored entries.
func FormatScores(scores []EntryScore) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %6s  %s\n", "KEY", "SCORE", "ISSUES"))
	sb.WriteString(strings.Repeat("-", 60) + "\n")
	for _, s := range scores {
		issues := "-"
		if len(s.Issues) > 0 {
			issues = strings.Join(s.Issues, "; ")
		}
		sb.WriteString(fmt.Sprintf("%-30s %6d  %s\n", s.Key, s.Score, issues))
	}
	return sb.String()
}

// AverageScore computes the mean score across all entries.
func AverageScore(scores []EntryScore) float64 {
	if len(scores) == 0 {
		return 0
	}
	total := 0
	for _, s := range scores {
		total += s.Score
	}
	return float64(total) / float64(len(scores))
}
