package env

import (
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeClassifyEntries(pairs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestClassify_SecretKeys(t *testing.T) {
	entries := makeClassifyEntries("API_KEY", "abc", "DB_PASSWORD", "secret")
	results := Classify(entries, DefaultClassifyOptions())
	if results[0].Category != CategorySecret {
		t.Errorf("expected secret, got %s", results[0].Category)
	}
	if results[1].Category != CategorySecret {
		t.Errorf("expected secret, got %s", results[1].Category)
	}
}

func TestClassify_DatabaseKeys(t *testing.T) {
	entries := makeClassifyEntries("DB_HOST", "localhost", "DATABASE_NAME", "mydb")
	results := Classify(entries, DefaultClassifyOptions())
	for _, r := range results {
		if r.Category != CategoryDatabase {
			t.Errorf("key %s: expected database, got %s", r.Entry.Key, r.Category)
		}
	}
}

func TestClassify_NetworkKeys(t *testing.T) {
	entries := makeClassifyEntries("APP_HOST", "0.0.0.0", "HTTP_PORT", "8080")
	results := Classify(entries, DefaultClassifyOptions())
	for _, r := range results {
		if r.Category != CategoryNetwork {
			t.Errorf("key %s: expected network, got %s", r.Entry.Key, r.Category)
		}
	}
}

func TestClassify_FeatureKeys(t *testing.T) {
	entries := makeClassifyEntries("FEATURE_DARK_MODE", "true", "ENABLE_BETA", "false")
	results := Classify(entries, DefaultClassifyOptions())
	for _, r := range results {
		if r.Category != CategoryFeature {
			t.Errorf("key %s: expected feature, got %s", r.Entry.Key, r.Category)
		}
	}
}

func TestClassify_GeneralKeys(t *testing.T) {
	entries := makeClassifyEntries("APP_ENV", "production", "LOG_LEVEL", "info")
	results := Classify(entries, DefaultClassifyOptions())
	for _, r := range results {
		if r.Category != CategoryGeneral {
			t.Errorf("key %s: expected general, got %s", r.Entry.Key, r.Category)
		}
	}
}

func TestClassify_CustomRules_TakePrecedence(t *testing.T) {
	opts := DefaultClassifyOptions()
	opts.CustomRules["billing"] = []string{"STRIPE_"}
	entries := makeClassifyEntries("STRIPE_SECRET_KEY", "sk_live_abc")
	results := Classify(entries, opts)
	if results[0].Category != Category("billing") {
		t.Errorf("expected billing, got %s", results[0].Category)
	}
}

func TestBuildClassifySummaries_Length(t *testing.T) {
	entries := makeClassifyEntries("FOO", "bar", "BAZ", "qux")
	classified := Classify(entries, DefaultClassifyOptions())
	summaries := BuildClassifySummaries(classified)
	if len(summaries) != 2 {
		t.Errorf("expected 2 summaries, got %d", len(summaries))
	}
}

func TestFormatClassify_ContainsHeaders(t *testing.T) {
	entries := makeClassifyEntries("APP_ENV", "prod")
	classified := Classify(entries, DefaultClassifyOptions())
	summaries := BuildClassifySummaries(classified)
	out := FormatClassify(summaries)
	for _, header := range []string{"KEY", "CATEGORY", "VALUE"} {
		if !containsStr(out, header) {
			t.Errorf("expected header %q in output", header)
		}
	}
}

func TestGroupByCategory_Buckets(t *testing.T) {
	entries := makeClassifyEntries("API_KEY", "x", "APP_ENV", "prod")
	classified := Classify(entries, DefaultClassifyOptions())
	summaries := BuildClassifySummaries(classified)
	groups := GroupByCategory(summaries)
	if _, ok := groups[CategorySecret]; !ok {
		t.Error("expected secret bucket")
	}
	if _, ok := groups[CategoryGeneral]; !ok {
		t.Error("expected general bucket")
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsSubstr(s, sub))
}

func containsSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
