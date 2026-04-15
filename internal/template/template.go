// Package template provides functionality for generating .env template files
// from existing env files, replacing values with placeholder descriptors.
package template

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// Options controls how the template is generated.
type Options struct {
	// IncludeDefaults keeps non-sensitive values as their actual defaults.
	IncludeDefaults bool
	// PlaceholderFormat is the format string for placeholders, e.g. "<KEY>".
	// Use %s to embed the key name.
	PlaceholderFormat string
}

var defaultOptions = Options{
	IncludeDefaults:   false,
	PlaceholderFormat: "<%s>",
}

// Generate produces a .env template string from the given EnvFile.
// Sensitive keys always get a placeholder; non-sensitive keys respect opts.IncludeDefaults.
func Generate(file parser.EnvFile, opts *Options) string {
	if opts == nil {
		opts = &defaultOptions
	}

	fmt = opts.PlaceholderFormat
	if fmt == "" {
		fmt = defaultOptions.PlaceholderFormat
	}

	var sb strings.Builder

	for _, entry := range file.Entries {
		if entry.Comment != "" {
			sb.WriteString("# ")
			sb.WriteString(entry.Comment)
			sb.WriteString("\n")
		}

		placeholder := buildPlaceholder(entry.Key, fmt)

		if isSensitiveKey(entry.Key) {
			sb.WriteString(entry.Key)
			sb.WriteString("=")
			sb.WriteString(placeholder)
		} else if opts.IncludeDefaults && entry.Value != "" {
			sb.WriteString(entry.Key)
			sb.WriteString("=")
			sb.WriteString(entry.Value)
		} else {
			sb.WriteString(entry.Key)
			sb.WriteString("=")
			sb.WriteString(placeholder)
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

func buildPlaceholder(key, format string) string {
	return fmt.Sprintf(format, key)
}

// isSensitiveKey mirrors the logic in secrets.IsSensitive without importing
// to keep the package dependency-light.
func isSensitiveKey(key string) bool {
	upper := strings.ToUpper(key)
	sensitiveTerms := []string{"SECRET", "PASSWORD", "PASSWD", "TOKEN", "APIKEY", "API_KEY", "PRIVATE", "CREDENTIAL"}
	for _, term := range sensitiveTerms {
		if strings.Contains(upper, term) {
			return true
		}
	}
	return false
}
