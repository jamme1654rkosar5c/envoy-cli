package compare_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envoy-cli/internal/compare"
)

func TestFormatTable_ContainsHeaders(t *testing.T) {
	envs := compare.EnvMap{
		"dev":  makeEntries("DB_HOST", "localhost"),
		"prod": makeEntries("DB_HOST", "prod.db"),
	}
	report := compare.CrossCompare(envs)
	table := compare.FormatTable(report)

	if !strings.Contains(table, "dev") {
		t.Error("table should contain env name 'dev'")
	}
	if !strings.Contains(table, "prod") {
		t.Error("table should contain env name 'prod'")
	}
	if !strings.Contains(table, "DB_HOST") {
		t.Error("table should contain key 'DB_HOST'")
	}
}

func TestFormatTable_ShowsMissing(t *testing.T) {
	envs := compare.EnvMap{
		"dev":  makeEntries("ONLY_DEV", "yes"),
		"prod": makeEntries("OTHER", "no"),
	}
	report := compare.CrossCompare(envs)
	table := compare.FormatTable(report)

	if !strings.Contains(table, "<missing>") {
		t.Error("expected <missing> marker in table")
	}
}

func TestFormatMissing_NoMissing(t *testing.T) {
	envs := compare.EnvMap{
		"dev":  makeEntries("KEY", "val"),
		"prod": makeEntries("KEY", "val"),
	}
	report := compare.CrossCompare(envs)
	out := compare.FormatMissing(report)
	if !strings.Contains(out, "No missing") {
		t.Errorf("expected no-missing message, got: %s", out)
	}
}

func TestFormatMissing_ListsMissingKeys(t *testing.T) {
	envs := compare.EnvMap{
		"dev":  makeEntries("SECRET_KEY", "abc", "PORT", "8080"),
		"prod": makeEntries("PORT", "80"),
	}
	report := compare.CrossCompare(envs)
	out := compare.FormatMissing(report)
	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected SECRET_KEY in missing output, got: %s", out)
	}
}
