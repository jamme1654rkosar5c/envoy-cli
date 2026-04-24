package env

import (
	"fmt"
	"strings"
)

const (
	symAdded     = "+"
	symRemoved   = "-"
	symChanged   = "~"
	symUnchanged = " "
)

// BuildDiffSummaries converts a slice of DiffEntry into human-readable summary
// lines suitable for display.
func BuildDiffSummaries(entries []DiffEntry) []string {
	var lines []string
	for _, e := range entries {
		switch e.Status {
		case "added":
			lines = append(lines, fmt.Sprintf("%s %s=%s", symAdded, e.Key, truncateDiffVal(e.NewValue)))
		case "removed":
			lines = append(lines, fmt.Sprintf("%s %s=%s", symRemoved, e.Key, truncateDiffVal(e.OldValue)))
		case "changed":
			lines = append(lines, fmt.Sprintf("%s %s: %s → %s", symChanged, e.Key, truncateDiffVal(e.OldValue), truncateDiffVal(e.NewValue)))
		case "unchanged":
			lines = append(lines, fmt.Sprintf("%s %s=%s", symUnchanged, e.Key, truncateDiffVal(e.NewValue)))
		}
	}
	return lines
}

// FormatDiff renders a diff table with a header and per-entry rows.
func FormatDiff(entries []DiffEntry) string {
	if len(entries) == 0 {
		return "(no differences)"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-3s %-30s %-20s %-20s\n", "Op", "Key", "Old", "New"))
	sb.WriteString(strings.Repeat("-", 77) + "\n")

	for _, e := range entries {
		op := statusSymbolDiff(e.Status)
		sb.WriteString(fmt.Sprintf("%-3s %-30s %-20s %-20s\n",
			op,
			truncateDiffKey(e.Key, 30),
			truncateDiffVal(e.OldValue),
			truncateDiffVal(e.NewValue),
		))
	}
	return sb.String()
}

func statusSymbolDiff(status string) string {
	switch status {
	case "added":
		return symAdded
	case "removed":
		return symRemoved
	case "changed":
		return symChanged
	default:
		return symUnchanged
	}
}

func truncateDiffVal(s string) string {
	if len(s) > 20 {
		return s[:17] + "..."
	}
	return s
}

func truncateDiffKey(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}
