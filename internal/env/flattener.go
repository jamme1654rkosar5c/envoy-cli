package env

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// FlattenOptions controls how nested key structures are flattened.
type FlattenOptions struct {
	// Separator is the delimiter used to join nested key segments (default: "__").
	Separator string
	// Prefix filters entries to only flatten keys with this prefix.
	Prefix string
	// StripPrefix removes the matched prefix from output keys.
	StripPrefix bool
}

// DefaultFlattenOptions returns sensible defaults for FlattenOptions.
func DefaultFlattenOptions() FlattenOptions {
	return FlattenOptions{
		Separator:   "__",
		StripPrefix: false,
	}
}

// Flatten normalises entries whose keys contain the separator into a canonical
// lowercase-segment representation. Duplicate resulting keys are deduplicated
// by keeping the last occurrence.
//
// Example: APP__DB__HOST=localhost → APP_DB_HOST=localhost (separator "__", join "_")
func Flatten(entries []parser.Entry, opts FlattenOptions) ([]parser.Entry, error) {
	if opts.Separator == "" {
		return nil, fmt.Errorf("flattener: separator must not be empty")
	}

	seen := make(map[string]int)
	result := make([]parser.Entry, 0, len(entries))

	for _, e := range entries {
		key := e.Key

		if opts.Prefix != "" && !strings.HasPrefix(key, opts.Prefix) {
			result = append(result, e)
			continue
		}

		if opts.StripPrefix && opts.Prefix != "" {
			key = strings.TrimPrefix(key, opts.Prefix)
		}

		segments := strings.Split(key, opts.Separator)
		flat := strings.Join(segments, "_")

		e.Key = flat

		if idx, exists := seen[flat]; exists {
			result[idx] = e
		} else {
			seen[flat] = len(result)
			result = append(result, e)
		}
	}

	return result, nil
}
