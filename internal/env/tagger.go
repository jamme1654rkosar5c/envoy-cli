package env

import (
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// TagOptions controls tagging behaviour.
type TagOptions struct {
	// Overwrite replaces existing tags if true.
	Overwrite bool
	// TagPrefix is the prefix used to identify tag comments, e.g. "@tag".
	TagPrefix string
}

// DefaultTagOptions returns sensible defaults.
func DefaultTagOptions() TagOptions {
	return TagOptions{
		Overwrite: false,
		TagPrefix: "@tag",
	}
}

// Tag sets a tag on the entry matching key.
// Tags are stored as a structured comment suffix: # @tag:<value>
func Tag(entries []parser.Entry, key, tag string, opts TagOptions) ([]parser.Entry, error) {
	out := make([]parser.Entry, len(entries))
	copy(out, entries)

	for i, e := range out {
		if e.Key != key {
			continue
		}
		marker := opts.TagPrefix + ":" + tag
		if !opts.Overwrite && strings.Contains(e.Comment, opts.TagPrefix) {
			return out, nil
		}
		// Strip any existing tag.
		base := stripTag(e.Comment, opts.TagPrefix)
		if base != "" {
			out[i].Comment = base + " " + marker
		} else {
			out[i].Comment = marker
		}
		return out, nil
	}
	return out, nil
}

// Untag removes the tag comment from the entry matching key.
func Untag(entries []parser.Entry, key string, opts TagOptions) []parser.Entry {
	out := make([]parser.Entry, len(entries))
	copy(out, entries)
	for i, e := range out {
		if e.Key == key {
			out[i].Comment = stripTag(e.Comment, opts.TagPrefix)
		}
	}
	return out
}

// GetTag returns the tag value for the given key, or empty string if none.
func GetTag(entries []parser.Entry, key string, opts TagOptions) string {
	for _, e := range entries {
		if e.Key != key {
			continue
		}
		return extractTagValue(e.Comment, opts.TagPrefix)
	}
	return ""
}

func stripTag(comment, prefix string) string {
	parts := strings.Fields(comment)
	var kept []string
	for _, p := range parts {
		if !strings.HasPrefix(p, prefix) {
			kept = append(kept, p)
		}
	}
	return strings.Join(kept, " ")
}

func extractTagValue(comment, prefix string) string {
	for _, p := range strings.Fields(comment) {
		if strings.HasPrefix(p, prefix+":") {
			return strings.TrimPrefix(p, prefix+":")
		}
	}
	return ""
}
