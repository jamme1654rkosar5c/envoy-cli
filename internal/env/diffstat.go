package env

import (
	"fmt"
	"strings"
)

// DiffStatEntry holds a single key's change summary.
type DiffStatEntry struct {
	Key    string
	Status string // "added", "removed", "changed", "unchanged"
	OldVal string
	NewVal string
}

// DiffStatOptions controls DiffStat behaviour.
type DiffStatOptions struct {
	IncludeUnchanged bool
	RedactValues     bool
}

// DefaultDiffStatOptions returns sensible defaults.
func DefaultDiffStatOptions() DiffStatOptions {
	return DiffStatOptions{
		IncludeUnchanged: false,
		RedactValues:     false,
	}
}

// DiffStat compares two slices of EnvEntry and returns a change summary.
func DiffStat(base, target []EnvEntry, opts DiffStatOptions) []DiffStatEntry {
	baseMap := toEntryMap(base)
	targetMap := toEntryMap(target)

	var results []DiffStatEntry

	for k, bv := range baseMap {
		if tv, ok := targetMap[k]; ok {
			if bv == tv {
				if opts.IncludeUnchanged {
					results = append(results, DiffStatEntry{Key: k, Status: "unchanged", OldVal: mask(bv, opts.RedactValues), NewVal: mask(tv, opts.RedactValues)})
				}
			} else {
				results = append(results, DiffStatEntry{Key: k, Status: "changed", OldVal: mask(bv, opts.RedactValues), NewVal: mask(tv, opts.RedactValues)})
			}
		} else {
			results = append(results, DiffStatEntry{Key: k, Status: "removed", OldVal: mask(bv, opts.RedactValues)})
		}
	}

	for k, tv := range targetMap {
		if _, ok := baseMap[k]; !ok {
			results = append(results, DiffStatEntry{Key: k, Status: "added", NewVal: mask(tv, opts.RedactValues)})
		}
	}

	return results
}

// FormatDiffStat returns a human-readable summary table.
func FormatDiffStat(entries []DiffStatEntry) string {
	if len(entries) == 0 {
		return "No differences found.\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-4s %-30s %-20s %-20s\n", "ST", "KEY", "OLD", "NEW"))
	sb.WriteString(strings.Repeat("-", 78) + "\n")
	for _, e := range entries {
		symbol := statusSymbol(e.Status)
		sb.WriteString(fmt.Sprintf("%-4s %-30s %-20s %-20s\n", symbol, e.Key, truncateStat(e.OldVal, 18), truncateStat(e.NewVal, 18)))
	}
	return sb.String()
}

func statusSymbol(s string) string {
	switch s {
	case "added":
		return "+"
	case "removed":
		return "-"
	case "changed":
		return "~"
	default:
		return "="
	}
}

func mask(v string, redact bool) string {
	if redact && v != "" {
		return "***"
	}
	return v
}

func truncateStat(s string, max int) string {
	if len(s) > max {
		return s[:max-1] + "…"
	}
	return s
}

func toEntryMap(entries []EnvEntry) map[string]string {
	m := make(map[string]string, len(entries))
	for _, e := range entries {
		m[e.Key] = e.Value
	}
	return m
}
