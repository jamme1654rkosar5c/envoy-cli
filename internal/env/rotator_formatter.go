package env

import (
	"fmt"
	"strings"
)

// RotateSummary describes the outcome of a single rotation.
type RotateSummary struct {
	OldKey  string
	NewKey  string
	Changed bool // true when the value was also replaced
}

// FormatRotations returns a human-readable table of rotation summaries.
func FormatRotations(summaries []RotateSummary) string {
	if len(summaries) == 0 {
		return "No rotations applied."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %-30s %s\n", "OLD KEY", "NEW KEY", "VALUE CHANGED"))
	sb.WriteString(strings.Repeat("-", 72) + "\n")
	for _, s := range summaries {
		changed := "no"
		if s.Changed {
			changed = "yes"
		}
		sb.WriteString(fmt.Sprintf("%-30s %-30s %s\n", s.OldKey, s.NewKey, changed))
	}
	return sb.String()
}

// BuildRotateSummaries compares rotation requests against old keys to produce
// a slice of RotateSummary values suitable for display.
func BuildRotateSummaries(rotations []RotateEntry) []RotateSummary {
	out := make([]RotateSummary, 0, len(rotations))
	for _, r := range rotations {
		out = append(out, RotateSummary{
			OldKey:  r.OldKey,
			NewKey:  r.NewKey,
			Changed: r.NewValue != "",
		})
	}
	return out
}
