package env

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// InspectResult holds detailed metadata about a single env entry.
type InspectResult struct {
	Key          string
	Value        string
	Comment      string
	HasComment   bool
	IsEmpty      bool
	IsQuoted     bool
	QuoteChar    string
	Length       int
	LineNumber   int
}

// InspectOptions controls Inspect behaviour.
type InspectOptions struct {
	// If true, return an error when the key is not found.
	ErrorOnMissing bool
}

// DefaultInspectOptions returns sensible defaults.
func DefaultInspectOptions() InspectOptions {
	return InspectOptions{ErrorOnMissing: true}
}

// Inspect returns metadata about a specific key in the env file.
func Inspect(entries []parser.Entry, key string, opts InspectOptions) (*InspectResult, error) {
	for i, e := range entries {
		if e.Key != key {
			continue
		}
		res := &InspectResult{
			Key:        e.Key,
			Value:      e.Value,
			Comment:    e.Comment,
			HasComment: strings.TrimSpace(e.Comment) != "",
			IsEmpty:    strings.TrimSpace(e.Value) == "",
			Length:     len(e.Value),
			LineNumber: i + 1,
		}
		if strings.HasPrefix(e.Value, `"`) && strings.HasSuffix(e.Value, `"`) {
			res.IsQuoted = true
			res.QuoteChar = `"`
		} else if strings.HasPrefix(e.Value, "'") && strings.HasSuffix(e.Value, "'") {
			res.IsQuoted = true
			res.QuoteChar = "'"
		}
		return res, nil
	}
	if opts.ErrorOnMissing {
		return nil, fmt.Errorf("inspect: key %q not found", key)
	}
	return nil, nil
}

// FormatInspect returns a human-readable summary of an InspectResult.
func FormatInspect(r *InspectResult) string {
	if r == nil {
		return "(not found)"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Key:        %s\n", r.Key)
	fmt.Fprintf(&sb, "Value:      %s\n", r.Value)
	fmt.Fprintf(&sb, "Length:     %d\n", r.Length)
	fmt.Fprintf(&sb, "Empty:      %v\n", r.IsEmpty)
	fmt.Fprintf(&sb, "Quoted:     %v", r.IsQuoted)
	if r.IsQuoted {
		fmt.Fprintf(&sb, " (%s)", r.QuoteChar)
	}
	sb.WriteString("\n")
	fmt.Fprintf(&sb, "HasComment: %v\n", r.HasComment)
	if r.HasComment {
		fmt.Fprintf(&sb, "Comment:    %s\n", r.Comment)
	}
	fmt.Fprintf(&sb, "Line:       %d\n", r.LineNumber)
	return sb.String()
}
