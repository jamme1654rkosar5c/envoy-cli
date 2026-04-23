package env

import (
	"strings"
	"testing"
)

func makeTraceResults() map[string]TraceResult {
	return map[string]TraceResult{
		"FOO": {Key: "FOO", Chain: []string{"FOO"}, Cycles: false, Depth: 0},
		"BAR": {Key: "BAR", Chain: []string{"BAR", "FOO"}, Cycles: false, Depth: 1},
		"BAZ": {Key: "BAZ", Chain: []string{"BAZ", "BAR", "FOO"}, Cycles: false, Depth: 2},
	}
}

func TestFormatTrace_ContainsHeaders(t *testing.T) {
	results := makeTraceResults()
	out := FormatTrace(results, []string{"FOO", "BAR", "BAZ"})
	if !strings.Contains(out, "KEY") {
		t.Error("expected header KEY")
	}
	if !strings.Contains(out, "DEPTH") {
		t.Error("expected header DEPTH")
	}
	if !strings.Contains(out, "CHAIN") {
		t.Error("expected header CHAIN")
	}
}

func TestFormatTrace_ShowsKeyAndDepth(t *testing.T) {
	results := makeTraceResults()
	out := FormatTrace(results, []string{"BAR"})
	if !strings.Contains(out, "BAR") {
		t.Error("expected BAR in output")
	}
	if !strings.Contains(out, "1") {
		t.Error("expected depth 1 in output")
	}
}

func TestFormatTrace_ShowsChain(t *testing.T) {
	results := makeTraceResults()
	out := FormatTrace(results, []string{"BAZ"})
	if !strings.Contains(out, "->") {
		t.Error("expected chain arrow in output")
	}
}

func TestFormatTrace_EmptyResults(t *testing.T) {
	out := FormatTrace(map[string]TraceResult{}, []string{})
	if !strings.Contains(out, "KEY") {
		t.Error("expected header even for empty results")
	}
}

func TestBuildTraceSummaries_CorrectCount(t *testing.T) {
	results := makeTraceResults()
	summaries := BuildTraceSummaries(results)
	if len(summaries) != 3 {
		t.Errorf("expected 3 summaries, got %d", len(summaries))
	}
}

func TestBuildTraceSummaries_ContainsKey(t *testing.T) {
	results := map[string]TraceResult{
		"FOO": {Key: "FOO", Chain: []string{"FOO"}, Depth: 0},
	}
	summaries := BuildTraceSummaries(results)
	if len(summaries) != 1 || !strings.Contains(summaries[0], "FOO") {
		t.Errorf("expected summary containing FOO, got %v", summaries)
	}
}
