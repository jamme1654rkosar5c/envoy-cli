package env

import (
	"strings"
	"testing"
)

func makeScores(pairs ...interface{}) []EntryScore {
	var out []EntryScore
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, EntryScore{
			Key:   pairs[i].(string),
			Score: pairs[i+1].(int),
		})
	}
	return out
}

func TestFormatScores_ContainsHeaders(t *testing.T) {
	scores := makeScores("APP_KEY", 100)
	out := FormatScores(scores)
	if !strings.Contains(out, "KEY") || !strings.Contains(out, "SCORE") {
		t.Error("expected headers in output")
	}
}

func TestFormatScores_ShowsKey(t *testing.T) {
	scores := makeScores("MY_VAR", 80)
	out := FormatScores(scores)
	if !strings.Contains(out, "MY_VAR") {
		t.Error("expected key in output")
	}
}

func TestFormatScores_ShowsDash_WhenNoIssues(t *testing.T) {
	scores := makeScores("CLEAN", 100)
	out := FormatScores(scores)
	if !strings.Contains(out, "-") {
		t.Error("expected dash for no issues")
	}
}

func TestAverageScore_Empty(t *testing.T) {
	if AverageScore(nil) != 0 {
		t.Error("expected 0 for empty")
	}
}

func TestAverageScore_Computed(t *testing.T) {
	scores := []EntryScore{{Score: 100}, {Score: 80}, {Score: 60}}
	avg := AverageScore(scores)
	if avg != 80.0 {
		t.Errorf("expected 80.0, got %f", avg)
	}
}
