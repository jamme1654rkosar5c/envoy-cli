package env

import (
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeClassifySummaries(pairs ...string) []ClassifySummary {
	out := make([]ClassifySummary, 0)
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, ClassifySummary{
			Key:      pairs[i],
			Category: CategoryGeneral,
			Value:    pairs[i+1],
		})
	}
	return out
}

func TestFormatClassify_EmptySummaries(t *testing.T) {
	out := FormatClassify(nil)
	if out != "no entries to classify\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatClassify_ShowsKeyAndCategory(t *testing.T) {
	entries := []parser.Entry{{Key: "DB_URL", Value: "postgres://localhost"}}
	classified := Classify(entries, DefaultClassifyOptions())
	summaries := BuildClassifySummaries(classified)
	out := FormatClassify(summaries)
	if !containsSubstr(out, "DB_URL") {
		t.Error("expected key DB_URL in output")
	}
	if !containsSubstr(out, string(CategoryDatabase)) {
		t.Errorf("expected category %s in output", CategoryDatabase)
	}
}

func TestTruncateClassVal_ShortValue(t *testing.T) {
	result := truncateClassVal("hello", 32)
	if result != "hello" {
		t.Errorf("expected 'hello', got %q", result)
	}
}

func TestTruncateClassVal_LongValue(t *testing.T) {
	long := "abcdefghijklmnopqrstuvwxyz0123456789extra"
	result := truncateClassVal(long, 10)
	if len(result) > 13 {
		t.Errorf("truncated value too long: %q", result)
	}
	if result[len(result)-3:] != "..." {
		t.Errorf("expected ellipsis suffix, got %q", result)
	}
}

func TestGroupByCategory_Empty(t *testing.T) {
	groups := GroupByCategory(nil)
	if len(groups) != 0 {
		t.Errorf("expected empty map, got %d entries", len(groups))
	}
}

func TestGroupByCategory_MultipleCategories(t *testing.T) {
	summaries := []ClassifySummary{
		{Key: "TOKEN", Category: CategorySecret},
		{Key: "HOST", Category: CategoryNetwork},
		{Key: "HOST2", Category: CategoryNetwork},
	}
	groups := GroupByCategory(summaries)
	if len(groups[CategoryNetwork]) != 2 {
		t.Errorf("expected 2 network entries, got %d", len(groups[CategoryNetwork]))
	}
	if len(groups[CategorySecret]) != 1 {
		t.Errorf("expected 1 secret entry, got %d", len(groups[CategorySecret]))
	}
}
