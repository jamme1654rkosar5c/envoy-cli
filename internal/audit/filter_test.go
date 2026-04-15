package audit_test

import (
	"testing"
	"time"

	"github.com/envoy-cli/internal/audit"
)

func makeLog(entries []audit.Entry) *audit.Log {
	return &audit.Log{Entries: entries}
}

func TestFilter_ByEventType(t *testing.T) {
	log := makeLog([]audit.Entry{
		{Event: audit.EventLoad, File: "a.env"},
		{Event: audit.EventDiff, File: "b.env"},
		{Event: audit.EventLoad, File: "c.env"},
	})

	result := audit.Filter(log, audit.FilterOptions{Event: audit.EventLoad})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestFilter_ByFile(t *testing.T) {
	log := makeLog([]audit.Entry{
		{Event: audit.EventLoad, File: "prod.env"},
		{Event: audit.EventLoad, File: "dev.env"},
	})

	result := audit.Filter(log, audit.FilterOptions{File: "prod.env"})
	if len(result) != 1 || result[0].File != "prod.env" {
		t.Errorf("expected 1 prod.env entry, got %d", len(result))
	}
}

func TestFilter_BySince(t *testing.T) {
	now := time.Now().UTC()
	log := makeLog([]audit.Entry{
		{Event: audit.EventLoad, Timestamp: now.Add(-2 * time.Hour)},
		{Event: audit.EventLoad, Timestamp: now.Add(-30 * time.Minute)},
		{Event: audit.EventLoad, Timestamp: now},
	})

	cutoff := now.Add(-1 * time.Hour)
	result := audit.Filter(log, audit.FilterOptions{Since: cutoff})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries after cutoff, got %d", len(result))
	}
}

func TestFilter_MaxRows(t *testing.T) {
	log := makeLog([]audit.Entry{
		{Event: audit.EventLoad},
		{Event: audit.EventDiff},
		{Event: audit.EventExport},
		{Event: audit.EventMerge},
	})

	result := audit.Filter(log, audit.FilterOptions{MaxRows: 2})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[1].Event != audit.EventMerge {
		t.Errorf("expected last entry to be merge, got %q", result[1].Event)
	}
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	log := makeLog([]audit.Entry{
		{Event: audit.EventLoad},
		{Event: audit.EventDiff},
	})

	result := audit.Filter(log, audit.FilterOptions{})
	if len(result) != 2 {
		t.Fatalf("expected all 2 entries, got %d", len(result))
	}
}
