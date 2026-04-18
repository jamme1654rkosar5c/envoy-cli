package env

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// AnnotateOptions controls annotation behaviour.
type AnnotateOptions struct {
	// Overwrite replaces an existing annotation if present.
	Overwrite bool
	// Prefix is prepended to the annotation text, e.g. "TODO", "NOTE".
	Prefix string
}

// DefaultAnnotateOptions returns sensible defaults.
func DefaultAnnotateOptions() AnnotateOptions {
	return AnnotateOptions{
		Overwrite: false,
		Prefix:    "",
	}
}

// Annotate sets a free-text annotation on the named key's comment field.
// The annotation is stored as "# @annotation <text>" appended to any existing comment.
func Annotate(entries []parser.Entry, key, text string, opts AnnotateOptions) ([]parser.Entry, error) {
	idx := -1
	for i, e := range entries {
		if e.Key == key {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil, fmt.Errorf("annotate: key %q not found", key)
	}

	annotation := buildAnnotation(opts.Prefix, text)

	existing := entries[idx].Comment
	if existing != "" && hasAnnotation(existing) {
		if !opts.Overwrite {
			return entries, nil
		}
		// Replace existing annotation.
		entries[idx].Comment = replaceAnnotation(existing, annotation)
		return entries, nil
	}

	if existing == "" {
		entries[idx].Comment = annotation
	} else {
		entries[idx].Comment = existing + " " + annotation
	}
	return entries, nil
}

// GetAnnotation returns the annotation text for the given key, or empty string.
func GetAnnotation(entries []parser.Entry, key string) string {
	for _, e := range entries {
		if e.Key == key && hasAnnotation(e.Comment) {
			return extractAnnotation(e.Comment)
		}
	}
	return ""
}

// RemoveAnnotation strips the annotation from the named key's comment.
func RemoveAnnotation(entries []parser.Entry, key string) []parser.Entry {
	for i, e := range entries {
		if e.Key == key {
			entries[i].Comment = stripAnnotation(e.Comment)
		}
	}
	return entries
}

func buildAnnotation(prefix, text string) string {
	if prefix != "" {
		return fmt.Sprintf("@annotation [%s] %s", prefix, text)
	}
	return "@annotation " + text
}

func hasAnnotation(comment string) bool {
	return strings.Contains(comment, "@annotation")
}

func extractAnnotation(comment string) string {
	idx := strings.Index(comment, "@annotation")
	if idx == -1 {
		return ""
	}
	return strings.TrimSpace(comment[idx+len("@annotation"):])
}

func replaceAnnotation(comment, newAnnotation string) string {
	idx := strings.Index(comment, "@annotation")
	if idx == -1 {
		return comment
	}
	return strings.TrimSpace(comment[:idx]) + newAnnotation
}

func stripAnnotation(comment string) string {
	idx := strings.Index(comment, "@annotation")
	if idx == -1 {
		return comment
	}
	return strings.TrimSpace(comment[:idx])
}
