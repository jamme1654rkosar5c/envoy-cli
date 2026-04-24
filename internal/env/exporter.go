package env

import (
	"fmt"
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// ExportFormat defines the output format for exported entries.
type ExportFormat string

const (
	ExportFormatDotEnv ExportFormat = "dotenv"
	ExportFormatShell  ExportFormat = "shell"
	ExportFormatInline ExportFormat = "inline"
)

// DefaultExportOptions returns sensible defaults for ExportEntries.
func DefaultExportOptions() ExportOptions {
	return ExportOptions{
		Format:      ExportFormatDotEnv,
		IncludeKeys: nil,
		ExcludeKeys: nil,
		QuoteValues: false,
	}
}

// ExportOptions controls how entries are serialised.
type ExportOptions struct {
	Format      ExportFormat
	IncludeKeys []string // if non-empty, only these keys are exported
	ExcludeKeys []string // keys to skip
	QuoteValues bool     // wrap values in double-quotes
}

// ExportEntries serialises env entries into the requested text format.
// It returns the formatted string and any error encountered.
func ExportEntries(entries []parser.Entry, opts ExportOptions) (string, error) {
	excludeSet := toExcludeSet(opts.ExcludeKeys)
	includeSet := toIncludeSet(opts.IncludeKeys)

	var sb strings.Builder
	for _, e := range entries {
		if len(includeSet) > 0 && !includeSet[e.Key] {
			continue
		}
		if excludeSet[e.Key] {
			continue
		}

		val := e.Value
		if opts.QuoteValues {
			val = `"` + val + `"`
		}

		switch opts.Format {
		case ExportFormatDotEnv:
			sb.WriteString(fmt.Sprintf("%s=%s\n", e.Key, val))
		case ExportFormatShell:
			sb.WriteString(fmt.Sprintf("export %s=%s\n", e.Key, val))
		case ExportFormatInline:
			sb.WriteString(fmt.Sprintf("%s=%s ", e.Key, val))
		default:
			return "", fmt.Errorf("unsupported export format: %q", opts.Format)
		}
	}

	return strings.TrimRight(sb.String(), " "), nil
}

func toExcludeSet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}

func toIncludeSet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
