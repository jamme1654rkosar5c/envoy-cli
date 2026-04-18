package env

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makeScorerEntries(kvs ...string) []parser.EnvEntry {
	entries := make([]parser.EnvEntry, 0, len(kvs)/2)
	for i := 0; i+1 < len(kvs); i += 2 {
		entries = append(entries, parser.EnvEntry{Key: kvs[i], Value: kvs[i+1]})
	}
	return entries
}

func TestScore_PerfectEntry(t *testing.T) {
	entries := makeScorerEntries("APP_NAME", "myapp")
	results := Score(entries, DefaultScoreOptions())
	if results[0].Score != 100 {
		t.Errorf("expected 100, got %d", results[0].Score)
	}
}

func TestScore_EmptyValue_Penalized(t *testing.T) {
	entries := makeScorerEntries("APP_NAME", "")
	results := Score(entries, DefaultScoreOptions())
	if results[0].Score != 70 {
		t.Errorf("expected 70, got %d", results[0].Score)
	}
	if len(results[0].Issues) == 0 {
		t.Error("expected issues")
	}
}

func TestScore_UnquotedSpaces_Penalized(t *testing.T) {
	entries := makeScorerEntries("MSG", "hello world")
	results := Score(entries, DefaultScoreOptions())
	if results[0].Score != 80 {
		t.Errorf("expected 80, got %d", results[0].Score)
	}
}

func TestScore_QuotedSpaces_NoPenalty(t *testing.T) {
	entries := makeScorerEntries("MSG", `"hello world"`)
	results := Score(entries, DefaultScoreOptions())
	if results[0].Score != 100 {
		t.Errorf("expected 100, got %d", results[0].Score)
	}
}

func TestScore_NoComment_Penalized_WhenEnabled(t *testing.T) {
	entries := makeScorerEntries("KEY", "val")
	opts := DefaultScoreOptions()
	opts.PenalizeNoComment = true
	results := Score(entries, opts)
	if results[0].Score != 90 {
		t.Errorf("expected 90, got %d", results[0].Score)
	}
}

func TestScore_ScoreFloorIsZero(t *testing.T) {
	entries := []parser.EnvEntry{{Key: "K", Value: ""}}
	opts := ScoreOptions{PenalizeEmpty: true, PenalizeUnquoted: true, PenalizeNoComment: true}
	results := Score(entries, opts)
	if results[0].Score < 0 {
		t.Error("score should not be negative")
	}
}
