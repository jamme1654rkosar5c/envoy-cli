package env

import (
	"strings"

	"github.com/user/envoy-cli/internal/parser"
)

// ClassifyOptions controls classification behaviour.
type ClassifyOptions struct {
	// CustomRules maps a category name to a list of key prefixes.
	CustomRules map[string][]string
}

// DefaultClassifyOptions returns sensible defaults.
func DefaultClassifyOptions() ClassifyOptions {
	return ClassifyOptions{
		CustomRules: map[string][]string{},
	}
}

// Category represents the inferred category of an env entry.
type Category string

const (
	CategorySecret   Category = "secret"
	CategoryDatabase Category = "database"
	CategoryNetwork  Category = "network"
	CategoryFeature  Category = "feature"
	CategoryGeneral  Category = "general"
)

// ClassifiedEntry pairs an env entry with its resolved category.
type ClassifiedEntry struct {
	Entry    parser.Entry
	Category Category
}

// Classify assigns a Category to each entry based on key patterns.
func Classify(entries []parser.Entry, opts ClassifyOptions) []ClassifiedEntry {
	result := make([]ClassifiedEntry, 0, len(entries))
	for _, e := range entries {
		result = append(result, ClassifiedEntry{
			Entry:    e,
			Category: classify(e.Key, opts),
		})
	}
	return result
}

func classify(key string, opts ClassifyOptions) Category {
	upper := strings.ToUpper(key)

	for cat, prefixes := range opts.CustomRules {
		for _, p := range prefixes {
			if strings.HasPrefix(upper, strings.ToUpper(p)) {
				return Category(cat)
			}
		}
	}

	switch {
	case matchesAny(upper, []string{"SECRET", "PASSWORD", "PASSWD", "TOKEN", "API_KEY", "PRIVATE_KEY", "AUTH"}):
		return CategorySecret
	case matchesAny(upper, []string{"DB_", "DATABASE_", "POSTGRES", "MYSQL", "MONGO", "REDIS"}):
		return CategoryDatabase
	case matchesAny(upper, []string{"HOST", "PORT", "URL", "ADDR", "ENDPOINT", "DNS"}):
		return CategoryNetwork
	case matchesAny(upper, []string{"FEATURE_", "FLAG_", "ENABLE_", "DISABLE_"}):
		return CategoryFeature
	default:
		return CategoryGeneral
	}
}

func matchesAny(key string, patterns []string) bool {
	for _, p := range patterns {
		if strings.Contains(key, p) {
			return true
		}
	}
	return false
}
