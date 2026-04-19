package env

import (
	"fmt"
	"strings"

	"github.com/your-org/envoy-cli/internal/parser"
)

// CascadeSummary describes how a key was resolved across layers.
type CascadeSummary struct {
	Key        string
	FinalValue string
	SourceLayer int
}

// BuildCascadeSummaries returns resolution metadata for each key across layers.
func BuildCascadeSummaries(layers [][]parser.Entry, opts CascadeOptions) []CascadeSummary {
	tracked := map[string]*CascadeSummary{}
	order := []string{}

	for layerIdx, layer := range layers {
		for _, e := range layer {
			if opts.SkipEmpty && e.Value == "" {
				continue
			}
			if s, exists := tracked[e.Key]; exists {
				if opts.Overwrite {
					s.FinalValue = e.Value
					s.SourceLayer = layerIdx
				}
			} else {
				tracked[e.Key] = &CascadeSummary{Key: e.Key, FinalValue: e.Value, SourceLayer: layerIdx}
				order = append(order, e.Key)
			}
		}
	}

	summaries := make([]CascadeSummary, 0, len(order))
	for _, k := range order {
		summaries = append(summaries, *tracked[k])
	}
	return summaries
}

// FormatCascade returns a human-readable table of cascade resolution results.
func FormatCascade(summaries []CascadeSummary) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-30s %-20s %s\n", "KEY", "VALUE", "LAYER"))
	sb.WriteString(strings.Repeat("-", 60) + "\n")
	for _, s := range summaries {
		val := s.FinalValue
		if len(val) > 18 {
			val = val[:15] + "..."
		}
		sb.WriteString(fmt.Sprintf("%-30s %-20s %d\n", s.Key, val, s.SourceLayer))
	}
	return sb.String()
}
