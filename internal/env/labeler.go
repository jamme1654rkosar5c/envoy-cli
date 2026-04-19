package env

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

const labelPrefix = "label:"

// DefaultLabelOptions returns sensible defaults for Label operations.
func DefaultLabelOptions() LabelOptions {
	return LabelOptions{
		Overwrite: false,
		AllowMissing: false,
	}
}

// LabelOptions controls behaviour of Label and Unlabel.
type LabelOptions struct {
	Overwrite    bool
	AllowMissing bool
}

// Label attaches a named label to the comment field of the matching entry.
func Label(entries []parser.Entry, key, label string, opts LabelOptions) ([]parser.Entry, error) {
	out := make([]parser.Entry, len(entries))
	copy(out, entries)

	for i, e := range out {
		if e.Key != key {
			continue
		}
		existing := GetLabel(out, key)
		if existing != "" && !opts.Overwrite {
			return out, fmt.Errorf("key %q already has label %q; use Overwrite to replace", key, existing)
		}
		base := stripLabel(e.Comment)
		if base != "" {
			out[i].Comment = base + " " + labelPrefix + label
		} else {
			out[i].Comment = labelPrefix + label
		}
		return out, nil
	}

	if !opts.AllowMissing {
		return out, fmt.Errorf("key %q not found", key)
	}
	return out, nil
}

// Unlabel removes the label annotation from the matching entry.
func Unlabel(entries []parser.Entry, key string) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	copy(out, entries)
	for i, e := range out {
		if e.Key == key {
			out[i].Comment = stripLabel(e.Comment)
		}
	}
	return out
}

// GetLabel returns the label value attached to the entry, or empty string.
func GetLabel(entries []parser.Entry, key string) string {
	for _, e := range entries {
		if e.Key != key {
			continue
		}
		for _, part := range strings.Fields(e.Comment) {
			if strings.HasPrefix(part, labelPrefix) {
				return strings.TrimPrefix(part, labelPrefix)
			}
		}
	}
	return ""
}

func stripLabel(comment string) string {
	parts := strings.Fields(comment)
	var kept []string
	for _, p := range parts {
		if !strings.HasPrefix(p, labelPrefix) {
			kept = append(kept, p)
		}
	}
	return strings.Join(kept, " ")
}
