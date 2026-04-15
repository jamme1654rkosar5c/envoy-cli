package merger

import (
	"fmt"

	"github.com/envoy-cli/internal/parser"
)

// Strategy defines how conflicts are resolved during a merge.
type Strategy int

const (
	// StrategyBase keeps the base value on conflict.
	StrategyBase Strategy = iota
	// StrategyOverride replaces base values with override values on conflict.
	StrategyOverride
	// StrategyError returns an error on conflict.
	StrategyError
)

// Result holds the merged env file and metadata about the merge.
type Result struct {
	File      parser.EnvFile
	Conflicts []string
	Added     []string
}

// Merge combines base and override EnvFiles according to the given strategy.
// Keys present only in override are always added to the result.
func Merge(base, override parser.EnvFile, strategy Strategy) (Result, error) {
	baseMap := toMap(base)
	result := Result{}

	// Start with a copy of base entries.
	merged := make([]parser.Entry, len(base.Entries))
	copy(merged, base.Entries)

	for _, entry := range override.Entries {
		if _, exists := baseMap[entry.Key]; exists {
			result.Conflicts = append(result.Conflicts, entry.Key)
			switch strategy {
			case StrategyError:
				return Result{}, fmt.Errorf("merge conflict on key: %s", entry.Key)
			case StrategyOverride:
				for i, e := range merged {
					if e.Key == entry.Key {
						merged[i].Value = entry.Value
						break
					}
				}
			case StrategyBase:
				// keep base value, do nothing
			}
		} else {
			merged = append(merged, entry)
			result.Added = append(result.Added, entry.Key)
		}
	}

	result.File = parser.EnvFile{Entries: merged}
	return result, nil
}

func toMap(f parser.EnvFile) map[string]string {
	m := make(map[string]string, len(f.Entries))
	for _, e := range f.Entries {
		m[e.Key] = e.Value
	}
	return m
}
