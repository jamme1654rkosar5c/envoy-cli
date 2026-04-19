package env

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// AliasSummary describes a single alias relationship.
type AliasSummary struct {
	Source string
	Alias  string
	Value  string
}

// ListAliases scans entries for alias comments and returns summaries.
func ListAliases(entries []parser.Entry) []AliasSummary {
	var out []AliasSummary
	for _, e := range entries {
		if strings.HasPrefix(e.Comment, "alias of ") {
			src := strings.TrimPrefix(e.Comment, "alias of ")
			out = append(out, AliasSummary{Source: src, Alias: e.Key, Value: e.Value})
		}
	}
	return out
}

// FormatAliases returns a human-readable table of alias summaries.
func FormatAliases(summaries []AliasSummary) string {
	if len(summaries) == 0 {
		return "No aliases defined.\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %-30s %s\n", "ALIAS", "SOURCE", "VALUE"))
	sb.WriteString(strings.Repeat("-", 80) + "\n")
	for _, s := range summaries {
		v := s.Value
		if len(v) > 16 {
			v = v[:13] + "..."
		}
		sb.WriteString(fmt.Sprintf("%-30s %-30s %s\n", s.Alias, s.Source, v))
	}
	return sb.String()
}
