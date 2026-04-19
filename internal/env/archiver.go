package env

import (
	"fmt"
	"time"

	"github.com/user/envoy-cli/internal/parser"
)

// ArchiveOptions controls archiving behaviour.
type ArchiveOptions struct {
	Prefix    string // prefix added to archived keys, default "ARCHIVED_"
	Timestamp bool   // append unix timestamp to prefix
	RemoveOriginal bool
}

// DefaultArchiveOptions returns sensible defaults.
func DefaultArchiveOptions() ArchiveOptions {
	return ArchiveOptions{
		Prefix:         "ARCHIVED_",
		Timestamp:      false,
		RemoveOriginal: true,
	}
}

// Archive renames matching keys by applying an archive prefix and optionally
// removes the originals from the entry list.
func Archive(entries []parser.Entry, keys []string, opts ArchiveOptions) ([]parser.Entry, error) {
	if opts.Prefix == "" {
		opts.Prefix = "ARCHIVED_"
	}

	prefix := opts.Prefix
	if opts.Timestamp {
		prefix = fmt.Sprintf("%s%d_", prefix, time.Now().Unix())
	}

	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[k] = true
	}

	result := make([]parser.Entry, 0, len(entries))
	archived := make(map[string]bool)

	for _, e := range entries {
		if keySet[e.Key] {
			newKey := prefix + e.Key
			result = append(result, parser.Entry{
				Key:     newKey,
				Value:   e.Value,
				Comment: fmt.Sprintf("archived from %s", e.Key),
			})
			archived[e.Key] = true
			if !opts.RemoveOriginal {
				result = append(result, e)
			}
		} else {
			result = append(result, e)
		}
	}

	for _, k := range keys {
		if !archived[k] {
			return nil, fmt.Errorf("archiver: key %q not found", k)
		}
	}

	return result, nil
}
