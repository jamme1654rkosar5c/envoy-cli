package compare

import (
	"fmt"
	"strings"
)

// FormatTable renders the comparison report as a plain-text table.
func FormatTable(r Report) string {
	var sb strings.Builder

	// Header
	header := fmt.Sprintf("%-30s", "KEY")
	for _, env := range r.Environments {
		header += fmt.Sprintf(" %-15s", env)
	}
	sb.WriteString(header + "\n")
	sb.WriteString(strings.Repeat("-", 30+16*len(r.Environments)) + "\n")

	for _, ks := range r.Keys {
		row := fmt.Sprintf("%-30s", ks.Key)
		for _, env := range r.Environments {
			if ks.Present[env] {
				row += fmt.Sprintf(" %-15s", truncate(ks.Values[env], 13))
			} else {
				row += fmt.Sprintf(" %-15s", "<missing>")
			}
		}
		sb.WriteString(row + "\n")
	}
	return sb.String()
}

// FormatMissing returns a focused report of missing keys per environment.
func FormatMissing(r Report) string {
	if len(r.MissingIn) == 0 {
		return "No missing keys across environments.\n"
	}
	var sb strings.Builder
	for _, env := range r.Environments {
		keys, ok := r.MissingIn[env]
		if !ok || len(keys) == 0 {
			continue
		}
		sb.WriteString(fmt.Sprintf("[%s] missing keys:\n", env))
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("  - %s\n", k))
		}
	}
	return sb.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
