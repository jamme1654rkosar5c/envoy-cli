package env

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// DefaultExpandOptions returns a safe default configuration for Expand.
func DefaultExpandOptions() ExpandOptions {
	return ExpandOptions{
		AllowMissing: false,
		MaxDepth:     10,
	}
}

// ExpandOptions controls the behaviour of Expand.
type ExpandOptions struct {
	// AllowMissing silences errors for unresolvable references.
	AllowMissing bool
	// MaxDepth limits recursive expansion to prevent infinite loops.
	MaxDepth int
}

// Expand resolves ${VAR} and $VAR references within entry values using other
// entries in the same slice as the lookup table. Entries are processed in
// order; earlier entries are visible to later ones.
func Expand(entries []parser.Entry, opts ExpandOptions) ([]parser.Entry, error) {
	lookup := make(map[string]string, len(entries))
	result := make([]parser.Entry, 0, len(entries))

	for _, e := range entries {
		expanded, err := expandValue(e.Value, lookup, opts, 0)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", e.Key, err)
		}
		e.Value = expanded
		lookup[e.Key] = expanded
		result = append(result, e)
	}

	return result, nil
}

func expandValue(val string, lookup map[string]string, opts ExpandOptions, depth int) (string, error) {
	if depth > opts.MaxDepth {
		return val, fmt.Errorf("expansion depth exceeded %d (possible cycle)", opts.MaxDepth)
	}

	var sb strings.Builder
	i := 0
	for i < len(val) {
		if val[i] != '$' {
			sb.WriteByte(val[i])
			i++
			continue
		}
		// skip escaped $$
		if i+1 < len(val) && val[i+1] == '$' {
			sb.WriteByte('$')
			i += 2
			continue
		}
		key, advance, err := extractRefKey(val, i)
		if err != nil {
			return "", err
		}
		resolved, ok := lookup[key]
		if !ok {
			if opts.AllowMissing {
				sb.WriteString(val[i : i+advance])
				i += advance
				continue
			}
			return "", fmt.Errorf("undefined variable %q", key)
		}
		inner, err := expandValue(resolved, lookup, opts, depth+1)
		if err != nil {
			return "", err
		}
		sb.WriteString(inner)
		i += advance
	}
	return sb.String(), nil
}

// extractRefKey parses a $VAR or ${VAR} reference starting at pos.
// Returns the key name, the number of bytes consumed, and any error.
func extractRefKey(s string, pos int) (string, int, error) {
	if pos+1 >= len(s) {
		return "", 1, nil
	}
	if s[pos+1] == '{' {
		end := strings.Index(s[pos+2:], "}")
		if end < 0 {
			return "", 0, fmt.Errorf("unclosed '{' in variable reference")
		}
		key := s[pos+2 : pos+2+end]
		return key, end + 3, nil
	}
	// bare $VAR — consume [A-Za-z0-9_]+
	j := pos + 1
	for j < len(s) && isIdentChar(s[j]) {
		j++
	}
	if j == pos+1 {
		return "", 1, nil
	}
	return s[pos+1 : j], j - pos, nil
}

func isIdentChar(b byte) bool {
	return (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') ||
		(b >= '0' && b <= '9') || b == '_'
}
