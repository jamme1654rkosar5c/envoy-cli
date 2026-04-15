package compare_test

import (
	"testing"

	"github.com/yourorg/envoy-cli/internal/compare"
	"github.com/yourorg/envoy-cli/internal/parser"
)

func makeEntries(pairs ...string) []parser.Entry {
	entries := make([]parser.Entry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestCrossCompare_AllPresent(t *testing.T) {
	envs := compare.EnvMap{
		"dev":  makeEntries("DB_HOST", "localhost", "PORT", "5432"),
		"prod": makeEntries("DB_HOST", "prod.db", "PORT", "5432"),
	}
	report := compare.CrossCompare(envs)
	if len(report.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(report.Keys))
	}
	for _, ks := range report.Keys {
		if !ks.Present["dev"] || !ks.Present["prod"] {
			t.Errorf("key %s should be present in both envs", ks.Key)
		}
	}
}

func TestCrossCompare_MissingKey(t *testing.T) {
	envs := compare.EnvMap{
		"dev":  makeEntries("DB_HOST", "localhost", "DEBUG", "true"),
		"prod": makeEntries("DB_HOST", "prod.db"),
	}
	report := compare.CrossCompare(envs)
	missing := report.MissingIn["prod"]
	if len(missing) != 1 || missing[0] != "DEBUG" {
		t.Errorf("expected DEBUG missing in prod, got %v", missing)
	}
}

func TestCrossCompare_DifferentValues(t *testing.T) {
	envs := compare.EnvMap{
		"dev":  makeEntries("API_URL", "http://dev.api"),
		"prod": makeEntries("API_URL", "https://prod.api"),
	}
	report := compare.CrossCompare(envs)
	if len(report.Keys) != 1 {
		t.Fatalf("expected 1 key")
	}
	ks := report.Keys[0]
	if ks.Values["dev"] == ks.Values["prod"] {
		t.Error("expected different values across envs")
	}
}

func TestCrossCompare_EmptyEnvs(t *testing.T) {
	report := compare.CrossCompare(compare.EnvMap{})
	if len(report.Keys) != 0 {
		t.Errorf("expected no keys for empty input")
	}
}

func TestSummary_Format(t *testing.T) {
	envs := compare.EnvMap{
		"dev":  makeEntries("KEY", "val"),
		"prod": makeEntries("KEY", "val", "EXTRA", "x"),
	}
	report := compare.CrossCompare(envs)
	s := compare.Summary(report)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
