package env

import (
	"fmt"
	"strings"
)

// ScopeSummary describes a single scoped entry result.
type ScopeSummary struct {
	OriginalKey string
	ScopedKey   string
	Value       string
}

// BuildScopeSummaries produces a summary of original vs scoped keys.
func BuildScopeSummaries(original, scoped []Entry) []ScopeSummary {
	origMap := make(map[string]string, len(original))
	for _, e := range original {
		origMap[e.Key] = e.Value
	}

	summaries := make([]ScopeSummary, 0, len(scoped))
	for _, e := range scoped {
		s := ScopeSummary{
			OriginalKey: e.Key,
			ScopedKey:   e.Key,
			Value:       e.Value,
		}
		// try to find the original key by value match
		for origKey, origVal := range origMap {
			if origVal == e.Value && origKey != e.Key {
				s.OriginalKey = origKey
				break
			}
		}
		summaries = append(summaries, s)
	}
	return summaries
}

// FormatScope returns a human-readable table of scoped entries.
func FormatScope(summaries []ScopeSummary) string {
	if len(summaries) == 0 {
		return "No scoped entries."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %-30s %s\n", "ORIGINAL KEY", "SCOPED KEY", "VALUE"))
	sb.WriteString(strings.Repeat("-", 80) + "\n")
	for _, s := range summaries {
		sb.WriteString(fmt.Sprintf("%-30s %-30s %s\n",
			truncateScopeVal(s.OriginalKey, 28),
			truncateScopeVal(s.ScopedKey, 28),
			truncateScopeVal(s.Value, 20),
		))
	}
	return sb.String()
}

func truncateScopeVal(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
