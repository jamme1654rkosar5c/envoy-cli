package env

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// MaskSummary holds display information for a single masked entry.
type MaskSummary struct {
	Key     string
	OrigLen int
	Masked  string
	WasMasked bool
}

// BuildMaskSummaries compares original and masked entry slices and returns
// per-entry summaries describing what was masked.
func BuildMaskSummaries(original, masked []parser.Entry) []MaskSummary {
	summaries := make([]MaskSummary, 0, len(original))
	for i, orig := range original {
		var m string
		wasMasked := false
		if i < len(masked) {
			m = masked[i].Value
			wasMasked = orig.Value != masked[i].Value
		}
		summaries = append(summaries, MaskSummary{
			Key:       orig.Key,
			OrigLen:   len(orig.Value),
			Masked:    m,
			WasMasked: wasMasked,
		})
	}
	return summaries
}

// FormatMask returns a human-readable table of mask results.
func FormatMask(summaries []MaskSummary) string {
	if len(summaries) == 0 {
		return "No entries to display.\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %-10s %-12s %s\n", "KEY", "MASKED", "ORIG LEN", "DISPLAY"))
	sb.WriteString(strings.Repeat("-", 70) + "\n")
	for _, s := range summaries {
		maskedFlag := "no"
		if s.WasMasked {
			maskedFlag = "yes"
		}
		display := s.Masked
		if display == "" {
			display = "(empty)"
		}
		sb.WriteString(fmt.Sprintf("%-30s %-10s %-12d %s\n",
			truncateMask(s.Key, 29), maskedFlag, s.OrigLen, display))
	}
	return sb.String()
}

func truncateMask(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
