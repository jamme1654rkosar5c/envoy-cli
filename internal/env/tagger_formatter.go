package env

import (
	"fmt"
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// TagSummary holds tag information for a single entry.
type TagSummary struct {
	Key   string
	Value string
	Tag   string
}

// ListTags returns a TagSummary for every entry that carries a tag.
func ListTags(entries []parser.Entry, opts TagOptions) []TagSummary {
	var out []TagSummary
	for _, e := range entries {
		v := extractTagValue(e.Comment, opts.TagPrefix)
		if v != "" {
			out = append(out, TagSummary{Key: e.Key, Value: e.Value, Tag: v})
		}
	}
	return out
}

// FormatTags returns a human-readable table of tagged entries.
func FormatTags(summaries []TagSummary) string {
	if len(summaries) == 0 {
		return "No tagged entries found.\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %-20s %s\n", "KEY", "TAG", "VALUE"))
	sb.WriteString(strings.Repeat("-", 70) + "\n")
	for _, s := range summaries {
		sb.WriteString(fmt.Sprintf("%-30s %-20s %s\n", s.Key, s.Tag, s.Value))
	}
	return sb.String()
}
