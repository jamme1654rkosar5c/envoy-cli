// Package compare provides environment comparison across multiple named environments.
package compare

import (
	"fmt"
	"sort"

	"github.com/yourorg/envoy-cli/internal/parser"
)

// EnvMap maps environment names to their parsed entries.
type EnvMap map[string][]parser.Entry

// KeyStatus describes the presence of a key across environments.
type KeyStatus struct {
	Key     string
	Present map[string]bool
	Values  map[string]string
}

// Report holds the full cross-environment comparison result.
type Report struct {
	Environments []string
	Keys         []KeyStatus
	MissingIn    map[string][]string // env -> keys missing in that env
}

// CrossCompare compares entries across multiple named environments.
func CrossCompare(envs EnvMap) Report {
	envNames := make([]string, 0, len(envs))
	for name := range envs {
		envNames = append(envNames, name)
	}
	sort.Strings(envNames)

	allKeys := collectAllKeys(envs)
	statuses := make([]KeyStatus, 0, len(allKeys))
	missingIn := make(map[string][]string)

	for _, key := range allKeys {
		ks := KeyStatus{
			Key:     key,
			Present: make(map[string]bool),
			Values:  make(map[string]string),
		}
		for _, env := range envNames {
			val, found := findKey(envs[env], key)
			ks.Present[env] = found
			if found {
				ks.Values[env] = val
			} else {
				missingIn[env] = append(missingIn[env], key)
			}
		}
		statuses = append(statuses, ks)
	}

	return Report{
		Environments: envNames,
		Keys:         statuses,
		MissingIn:    missingIn,
	}
}

// Summary returns a human-readable summary string for the report.
func Summary(r Report) string {
	total := len(r.Keys)
	missingCount := 0
	for _, keys := range r.MissingIn {
		missingCount += len(keys)
	}
	return fmt.Sprintf("Environments: %v | Total keys: %d | Missing entries: %d",
		r.Environments, total, missingCount)
}

func collectAllKeys(envs EnvMap) []string {
	seen := make(map[string]bool)
	for _, entries := range envs {
		for _, e := range entries {
			seen[e.Key] = true
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func findKey(entries []parser.Entry, key string) (string, bool) {
	for _, e := range entries {
		if e.Key == key {
			return e.Value, true
		}
	}
	return "", false
}
