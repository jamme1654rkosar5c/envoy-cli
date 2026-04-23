package env

import (
	"fmt"
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// TraceOptions controls how variable reference tracing behaves.
type TraceOptions struct {
	MaxDepth    int
	AllowCycles bool
}

// DefaultTraceOptions returns sensible defaults for tracing.
func DefaultTraceOptions() TraceOptions {
	return TraceOptions{
		MaxDepth:    16,
		AllowCycles: false,
	}
}

// TraceResult holds the resolved chain for a single key.
type TraceResult struct {
	Key    string
	Chain  []string
	Cycles bool
	Depth  int
}

// Trace resolves the full reference chain for each key in entries.
// It returns a map of key -> TraceResult describing how each value was derived.
func Trace(entries []parser.Entry, opts TraceOptions) (map[string]TraceResult, error) {
	lookup := make(map[string]string, len(entries))
	for _, e := range entries {
		lookup[e.Key] = e.Value
	}

	results := make(map[string]TraceResult, len(entries))
	for _, e := range entries {
		chain, cycles, err := resolveChain(e.Key, lookup, opts, nil)
		if err != nil {
			return nil, fmt.Errorf("trace %q: %w", e.Key, err)
		}
		results[e.Key] = TraceResult{
			Key:    e.Key,
			Chain:  chain,
			Cycles: cycles,
			Depth:  len(chain) - 1,
		}
	}
	return results, nil
}

func resolveChain(key string, lookup map[string]string, opts TraceOptions, visited []string) ([]string, bool, error) {
	for _, v := range visited {
		if v == key {
			if opts.AllowCycles {
				return append(visited, key+"(cycle)"), true, nil
			}
			return nil, true, fmt.Errorf("cycle detected at %q", key)
		}
	}
	if len(visited) >= opts.MaxDepth {
		return nil, false, fmt.Errorf("max depth %d exceeded", opts.MaxDepth)
	}

	visited = append(visited, key)
	val, ok := lookup[key]
	if !ok {
		return visited, false, nil
	}

	ref := extractRef(val)
	if ref == "" || ref == key {
		return visited, false, nil
	}
	return resolveChain(ref, lookup, opts, visited)
}

func extractRef(val string) string {
	val = strings.TrimSpace(val)
	if strings.HasPrefix(val, "${") && strings.HasSuffix(val, "}") {
		return val[2 : len(val)-1]
	}
	if strings.HasPrefix(val, "$") {
		rest := val[1:]
		for i, c := range rest {
			if !isIdentChar(c) {
				return rest[:i]
			}
		}
		return rest
	}
	return ""
}
