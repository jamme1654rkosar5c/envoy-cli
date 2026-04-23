package env

import (
	"fmt"
	"strings"
)

// FormatTrace renders a human-readable table of trace results.
func FormatTrace(results map[string]TraceResult, keys []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-24s %-6s %-8s %s\n", "KEY", "DEPTH", "CYCLES", "CHAIN"))
	sb.WriteString(strings.Repeat("-", 72) + "\n")
	for _, k := range keys {
		r, ok := results[k]
		if !ok {
			continue
		}
		cycles := "no"
		if r.Cycles {
			cycles = "YES"
		}
		chain := strings.Join(r.Chain, " -> ")
		if len(chain) > 40 {
			chain = chain[:37] + "..."
		}
		sb.WriteString(fmt.Sprintf("%-24s %-6d %-8s %s\n", truncateTraceKey(k, 24), r.Depth, cycles, chain))
	}
	return sb.String()
}

// BuildTraceSummaries returns a slice of formatted chain strings keyed by variable name.
func BuildTraceSummaries(results map[string]TraceResult) []string {
	out := make([]string, 0, len(results))
	for k, r := range results {
		out = append(out, fmt.Sprintf("%s: %s", k, strings.Join(r.Chain, " -> ")))
	}
	return out
}

func truncateTraceKey(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
