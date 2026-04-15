package diff_test

import (
	"testing"

	"github.com/envoy-cli/internal/diff"
	"github.com/envoy-cli/internal/parser"
)

func makeFile(entries []parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func TestCompare_NoChanges(t *testing.T) {
	base := makeFile([]parser.Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "PORT", Value: "8080"},
	})
	target := makeFile([]parser.Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "PORT", Value: "8080"},
	})
	result := diff.Compare(base, target)
	if result.HasChanges() {
		t.Errorf("expected no changes, got: %s", result.Summary())
	}
}

func TestCompare_AddedKey(t *testing.T) {
	base := makeFile([]parser.Entry{{Key: "PORT", Value: "8080"}})
	target := makeFile([]parser.Entry{
		{Key: "PORT", Value: "8080"},
		{Key: "NEW_KEY", Value: "value"},
	})
	result := diff.Compare(base, target)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != diff.Added || result.Changes[0].Key != "NEW_KEY" {
		t.Errorf("unexpected change: %+v", result.Changes[0])
	}
}

func TestCompare_RemovedKey(t *testing.T) {
	base := makeFile([]parser.Entry{
		{Key: "PORT", Value: "8080"},
		{Key: "OLD_KEY", Value: "old"},
	})
	target := makeFile([]parser.Entry{{Key: "PORT", Value: "8080"}})
	result := diff.Compare(base, target)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	if result.Changes[0].Type != diff.Removed || result.Changes[0].Key != "OLD_KEY" {
		t.Errorf("unexpected change: %+v", result.Changes[0])
	}
}

func TestCompare_ModifiedValue(t *testing.T) {
	base := makeFile([]parser.Entry{{Key: "DB_URL", Value: "localhost:5432"}})
	target := makeFile([]parser.Entry{{Key: "DB_URL", Value: "prod-db:5432"}})
	result := diff.Compare(base, target)
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result.Changes))
	}
	c := result.Changes[0]
	if c.Type != diff.Modified {
		t.Errorf("expected Modified, got %s", c.Type)
	}
	if c.OldValue != "localhost:5432" || c.NewValue != "prod-db:5432" {
		t.Errorf("unexpected values: %+v", c)
	}
}

func TestCompare_MixedChanges(t *testing.T) {
	base := makeFile([]parser.Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	})
	target := makeFile([]parser.Entry{
		{Key: "A", Value: "99"},
		{Key: "C", Value: "3"},
	})
	result := diff.Compare(base, target)
	if len(result.Changes) != 3 {
		t.Fatalf("expected 3 changes, got %d: %s", len(result.Changes), result.Summary())
	}
}
