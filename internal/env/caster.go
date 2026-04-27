package env

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/envoy-cli/envoy/internal/parser"
)

// CastType represents the target type for casting an env value.
type CastType string

const (
	CastString  CastType = "string"
	CastInt     CastType = "int"
	CastFloat   CastType = "float"
	CastBool    CastType = "bool"
)

// CastOptions controls how values are cast.
type CastOptions struct {
	// Keys is the set of keys to cast. If empty, all keys are attempted.
	Keys []string
	// TargetType is the type to cast values to.
	TargetType CastType
	// SkipInvalid skips entries that cannot be cast instead of returning an error.
	SkipInvalid bool
}

// DefaultCastOptions returns sensible defaults for Cast.
func DefaultCastOptions() CastOptions {
	return CastOptions{
		TargetType:  CastString,
		SkipInvalid: false,
	}
}

// CastResult holds the result of casting a single entry.
type CastResult struct {
	Key      string
	Original string
	Casted   string
	Type     CastType
	OK       bool
}

// Cast attempts to cast the values of matching entries to the given type,
// normalising their string representation (e.g. "true"/"false" for bools).
// It returns a new slice of entries with updated values and a summary of results.
func Cast(entries []parser.Entry, opts CastOptions) ([]parser.Entry, []CastResult, error) {
	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = true
	}

	out := make([]parser.Entry, len(entries))
	copy(out, entries)

	var results []CastResult

	for i, e := range out {
		if len(keySet) > 0 && !keySet[e.Key] {
			continue
		}
		if e.Key == "" {
			continue
		}

		casted, ok, err := castValue(e.Value, opts.TargetType)
		if err != nil {
			if opts.SkipInvalid {
				results = append(results, CastResult{Key: e.Key, Original: e.Value, Casted: e.Value, Type: opts.TargetType, OK: false})
				continue
			}
			return nil, nil, fmt.Errorf("cast: key %q value %q cannot be cast to %s: %w", e.Key, e.Value, opts.TargetType, err)
		}

		out[i].Value = casted
		results = append(results, CastResult{Key: e.Key, Original: e.Value, Casted: casted, Type: opts.TargetType, OK: ok})
	}

	return out, results, nil
}

func castValue(v string, t CastType) (string, bool, error) {
	v = strings.TrimSpace(v)
	switch t {
	case CastString:
		return v, true, nil
	case CastInt:
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return v, false, err
		}
		return strconv.FormatInt(n, 10), true, nil
	case CastFloat:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return v, false, err
		}
		return strconv.FormatFloat(f, 'f', -1, 64), true, nil
	case CastBool:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return v, false, err
		}
		return strconv.FormatBool(b), true, nil
	default:
		return v, false, fmt.Errorf("unknown cast type: %s", t)
	}
}
