package env

import (
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// DefaultMaskOptions returns sensible defaults for masking.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		MaskChar:        '*',
		VisibleSuffix:   0,
		SensitiveKeys:   []string{"SECRET", "PASSWORD", "TOKEN", "KEY", "PRIVATE", "PASS", "CREDENTIAL"},
		MaskAllValues:   false,
		Placeholder:     "[REDACTED]",
		UsePlaceholder:  false,
	}
}

// MaskOptions controls how values are masked.
type MaskOptions struct {
	// MaskChar is the character used to replace value characters.
	MaskChar rune
	// VisibleSuffix is the number of trailing characters to leave visible.
	VisibleSuffix int
	// SensitiveKeys is a list of substrings that mark a key as sensitive.
	SensitiveKeys []string
	// MaskAllValues masks every value regardless of key name.
	MaskAllValues bool
	// Placeholder replaces the entire value when UsePlaceholder is true.
	Placeholder string
	// UsePlaceholder replaces the value with Placeholder instead of masking chars.
	UsePlaceholder bool
}

// Mask returns a new slice of entries with sensitive values masked.
// Original entries are not mutated.
func Mask(entries []parser.Entry, opts MaskOptions) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		copy := e
		if e.Value != "" && (opts.MaskAllValues || isSensitiveMaskKey(e.Key, opts.SensitiveKeys)) {
			if opts.UsePlaceholder {
				copy.Value = opts.Placeholder
			} else {
				copy.Value = maskString(e.Value, opts.MaskChar, opts.VisibleSuffix)
			}
		}
		out[i] = copy
	}
	return out
}

func isSensitiveMaskKey(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

func maskString(s string, ch rune, visibleSuffix int) string {
	runes := []rune(s)
	n := len(runes)
	visible := visibleSuffix
	if visible > n {
		visible = n
	}
	maskLen := n - visible
	var b strings.Builder
	for i := 0; i < maskLen; i++ {
		b.WriteRune(ch)
	}
	for i := maskLen; i < n; i++ {
		b.WriteRune(runes[i])
	}
	return b.String()
}
