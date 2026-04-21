package env

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// ReorderSummary describes the positional change of a single entry.
type ReorderSummary struct {
	Key      string
	OldIndex int
	NewIndex int
	Moved    bool
}

// BuildReorderSummaries compares original and reordered entry slices and
// returns a summary of positional changes.
func BuildReorderSummaries(original, reordered []parser.Entry) []ReorderSummary {
	oldPos := make(map[string]int, len(original))
	for i, e := range original {
		oldPos[e.Key] = i
	}

	summaries := make([]ReorderSummary, 0, len(reordered))
	for newIdx, e := range reordered {
		oldIdx := oldPos[e.Key]
		summaries = append(summaries, ReorderSummary{
			Key:      e.Key,
			OldIndex: oldIdx,
			NewIndex: newIdx,
			Moved:    oldIdx != newIdx,
		})
	}
	return summaries
}

// FormatReorder returns a human-readable table of reorder changes.
func FormatReorder(summaries []ReorderSummary) string {
	if len(summaries) == 0 {
		return "No entries to display.\n"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%-30s %8s %8s %8s\n", "KEY", "OLD POS", "NEW POS", "MOVED")
	fmt.Fprintf(&sb, "%s\n", strings.Repeat("-", 58))
	for _, s := range summaries {
		moved := "no"
		if s.Moved {
			moved = "yes"
		}
		fmt.Fprintf(&sb, "%-30s %8d %8d %8s\n", truncateReorderKey(s.Key, 29), s.OldIndex, s.NewIndex, moved)
	}
	return sb.String()
}

func truncateReorderKey(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
