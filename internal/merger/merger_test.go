package merger

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makeFile(pairs ...string) parser.EnvFile {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return parser.EnvFile{Entries: entries}
}

func TestMerge_NoConflicts(t *testing.T) {
	base := makeFile("APP_ENV", "production", "DB_HOST", "localhost")
	override := makeFile("LOG_LEVEL", "debug")

	res, err := Merge(base, override, StrategyBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.File.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(res.File.Entries))
	}
	if len(res.Added) != 1 || res.Added[0] != "LOG_LEVEL" {
		t.Errorf("expected LOG_LEVEL in Added, got %v", res.Added)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
}

func TestMerge_StrategyBase_KeepsBaseValue(t *testing.T) {
	base := makeFile("APP_ENV", "production")
	override := makeFile("APP_ENV", "staging")

	res, err := Merge(base, override, StrategyBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.File.Entries[0].Value != "production" {
		t.Errorf("expected base value 'production', got %q", res.File.Entries[0].Value)
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
}

func TestMerge_StrategyOverride_ReplacesValue(t *testing.T) {
	base := makeFile("APP_ENV", "production")
	override := makeFile("APP_ENV", "staging")

	res, err := Merge(base, override, StrategyOverride)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.File.Entries[0].Value != "staging" {
		t.Errorf("expected override value 'staging', got %q", res.File.Entries[0].Value)
	}
}

func TestMerge_StrategyError_ReturnsError(t *testing.T) {
	base := makeFile("APP_ENV", "production")
	override := makeFile("APP_ENV", "staging")

	_, err := Merge(base, override, StrategyError)
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMerge_EmptyOverride(t *testing.T) {
	base := makeFile("APP_ENV", "production", "DB_HOST", "localhost")
	override := makeFile()

	res, err := Merge(base, override, StrategyBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.File.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.File.Entries))
	}
}

func TestMerge_MultipleConflicts(t *testing.T) {
	base := makeFile("APP_ENV", "production", "DB_HOST", "localhost", "LOG_LEVEL", "info")
	override := makeFile("APP_ENV", "staging", "DB_HOST", "remotehost", "LOG_LEVEL", "debug")

	res, err := Merge(base, override, StrategyBase)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 3 {
		t.Errorf("expected 3 conflicts, got %d", len(res.Conflicts))
	}
	// Base strategy: all original values should be preserved
	for _, e := range res.File.Entries {
		switch e.Key {
		case "APP_ENV":
			if e.Value != "production" {
				t.Errorf("expected 'production', got %q", e.Value)
			}
		case "DB_HOST":
			if e.Value != "localhost" {
				t.Errorf("expected 'localhost', got %q", e.Value)
			}
		case "LOG_LEVEL":
			if e.Value != "info" {
				t.Errorf("expected 'info', got %q", e.Value)
			}
		}
	}
}
