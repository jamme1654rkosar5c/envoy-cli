package env

import (
	"fmt"
	"strings"
)

// DigestSummary holds a labelled digest result for display.
type DigestSummary struct {
	Label     string
	Algorithm DigestAlgorithm
	Hash      string
	KeyCount  int
}

// FormatDigest returns a human-readable table of digest summaries.
func FormatDigest(summaries []DigestSummary) string {
	if len(summaries) == 0 {
		return "No digests to display.\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-20s %-8s %-64s %s\n", "LABEL", "ALGO", "HASH", "KEYS"))
	sb.WriteString(strings.Repeat("-", 100) + "\n")

	for _, s := range summaries {
		sb.WriteString(fmt.Sprintf("%-20s %-8s %-64s %d\n",
			truncateDigestLabel(s.Label, 20),
			string(s.Algorithm),
			s.Hash,
			s.KeyCount,
		))
	}
	return sb.String()
}

// BuildDigestSummary creates a DigestSummary from a label, hash, algorithm, and key count.
func BuildDigestSummary(label string, hash string, algo DigestAlgorithm, keyCount int) DigestSummary {
	return DigestSummary{
		Label:     label,
		Algorithm: algo,
		Hash:      hash,
		KeyCount:  keyCount,
	}
}

func truncateDigestLabel(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
