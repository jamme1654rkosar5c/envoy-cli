package env

import "github.com/your-org/envoy-cli/internal/parser"

// CascadeOptions controls how cascading merges are applied.
type CascadeOptions struct {
	// Overwrite allows later layers to overwrite earlier ones.
	Overwrite bool
	// SkipEmpty ignores entries with empty values in source layers.
	SkipEmpty bool
}

// DefaultCascadeOptions returns sensible defaults.
func DefaultCascadeOptions() CascadeOptions {
	return CascadeOptions{
		Overwrite: true,
		SkipEmpty: false,
	}
}

// Cascade applies multiple env layers in order, with each layer overriding
// the previous according to the given options. The first layer is the base.
func Cascade(layers [][]parser.Entry, opts CascadeOptions) []parser.Entry {
	if len(layers) == 0 {
		return nil
	}

	result := make([]parser.Entry, 0, len(layers[0]))
	result = append(result, layers[0]...)

	for _, layer := range layers[1:] {
		for _, src := range layer {
			if opts.SkipEmpty && src.Value == "" {
				continue
			}
			idx := indexInEntries(result, src.Key)
			if idx >= 0 {
				if opts.Overwrite {
					result[idx].Value = src.Value
				}
			} else {
				result = append(result, src)
			}
		}
	}
	return result
}

func indexInEntries(entries []parser.Entry, key string) int {
	for i, e := range entries {
		if e.Key == key {
			return i
		}
	}
	return -1
}
