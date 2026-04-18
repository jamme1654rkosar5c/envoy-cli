package env

import (
	"strings"

	"github.com/envoy-cli/internal/parser"
)

// ScoreOptions controls how entries are scored.
type ScoreOptions struct {
	// PenalizeEmpty reduces score for entries with empty values.
	PenalizeEmpty bool
	// PenalizeUnquoted reduces score for values containing spaces but no quotes.
	PenalizeUnquoted bool
	// PenalizeNoComment reduces score for entries without comments.
	PenalizeNoComment bool
}

// DefaultScoreOptions returns sensible defaults.
func DefaultScoreOptions() ScoreOptions {
	return ScoreOptions{
		PenalizeEmpty:     true,
		PenalizeUnquoted:  true,
		PenalizeNoComment: false,
	}
}

// EntryScore holds the score result for a single entry.
type EntryScore struct {
	Key    string
	Score  int
	Issues []string
}

// Score evaluates the quality of each entry and returns scored results.
// A perfect entry scores 100; penalties are subtracted.
func Score(entries []parser.EnvEntry, opts ScoreOptions) []EntryScore {
	results := make([]EntryScore, 0, len(entries))
	for _, e := range entries {
		es := EntryScore{Key: e.Key, Score: 100}
		if opts.PenalizeEmpty && strings.TrimSpace(e.Value) == "" {
			es.Score -= 30
			es.Issues = append(es.Issues, "empty value")
		}
		if opts.PenalizeUnquoted && strings.Contains(e.Value, " ") &&
			!strings.HasPrefix(e.Value, "\"") && !strings.HasPrefix(e.Value, "'") {
			es.Score -= 20
			es.Issues = append(es.Issues, "unquoted value with spaces")
		}
		if opts.PenalizeNoComment && strings.TrimSpace(e.Comment) == "" {
			es.Score -= 10
			es.Issues = append(es.Issues, "missing comment")
		}
		if es.Score < 0 {
			es.Score = 0
		}
		results = append(results, es)
	}
	return results
}
