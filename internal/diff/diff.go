package diff

import (
	"fmt"
	"sort"
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// ChangeType represents the kind of difference between two env files.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
)

// Change describes a single difference between two env files.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Result holds the full diff between two env files.
type Result struct {
	Changes []Change
}

// HasChanges returns true if there are any differences.
func (r *Result) HasChanges() bool {
	return len(r.Changes) > 0
}

// Summary returns a human-readable summary of all changes.
func (r *Result) Summary() string {
	if !r.HasChanges() {
		return "No differences found."
	}
	var sb strings.Builder
	for _, c := range r.Changes {
		switch c.Type {
		case Added:
			fmt.Fprintf(&sb, "+ %s=%s\n", c.Key, c.NewValue)
		case Removed:
			fmt.Fprintf(&sb, "- %s=%s\n", c.Key, c.OldValue)
		case Modified:
			fmt.Fprintf(&sb, "~ %s: %q -> %q\n", c.Key, c.OldValue, c.NewValue)
		}
	}
	return sb.String()
}

// Compare computes the diff between two parsed env files.
// base is the reference (e.g. .env.example), target is the file being compared.
func Compare(base, target *parser.EnvFile) *Result {
	baseMap := toMap(base)
	targetMap := toMap(target)

	keys := unionKeys(baseMap, targetMap)
	sort.Strings(keys)

	var changes []Change
	for _, key := range keys {
		baseVal, inBase := baseMap[key]
		targetVal, inTarget := targetMap[key]

		switch {
		case inBase && !inTarget:
			changes = append(changes, Change{Key: key, Type: Removed, OldValue: baseVal})
		case !inBase && inTarget:
			changes = append(changes, Change{Key: key, Type: Added, NewValue: targetVal})
		case baseVal != targetVal:
			changes = append(changes, Change{Key: key, Type: Modified, OldValue: baseVal, NewValue: targetVal})
		}
	}
	return &Result{Changes: changes}
}

func toMap(f *parser.EnvFile) map[string]string {
	m := make(map[string]string, len(f.Entries))
	for _, e := range f.Entries {
		m[e.Key] = e.Value
	}
	return m
}

func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}
