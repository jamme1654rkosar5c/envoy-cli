package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/envoy-cli/internal/parser"
	"github.com/envoy-cli/internal/snapshot"
)

func makeEntries() []parser.EnvEntry {
	return []parser.EnvEntry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "DEBUG", Value: "true"},
		{Key: "SECRET_KEY", Value: "s3cr3t"},
	}
}

func TestSave_CreatesSnapshotFile(t *testing.T) {
	dir := t.TempDir()
	err := snapshot.Save(dir, "staging", makeEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	paths, err := snapshot.List(dir)
	if err != nil {
		t.Fatalf("unexpected error listing: %v", err)
	}
	if len(paths) != 1 {
		t.Errorf("expected 1 snapshot, got %d", len(paths))
	}
}

func TestLoad_ReturnsCorrectSnapshot(t *testing.T) {
	dir := t.TempDir()
	entries := makeEntries()

	if err := snapshot.Save(dir, "prod", entries); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	paths, _ := snapshot.List(dir)
	snap, err := snapshot.Load(paths[0])
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if snap.Label != "prod" {
		t.Errorf("expected label 'prod', got %q", snap.Label)
	}
	if len(snap.Entries) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(snap.Entries))
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestList_EmptyDir_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	paths, err := snapshot.List(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 0 {
		t.Errorf("expected 0 paths, got %d", len(paths))
	}
}

func TestList_NonExistentDir_ReturnsEmpty(t *testing.T) {
	paths, err := snapshot.List("/tmp/envoy-cli-nonexistent-dir-xyz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 0 {
		t.Errorf("expected 0 paths, got %d", len(paths))
	}
}

func TestLoad_InvalidFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	badFile := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(badFile, []byte("not valid json{"), 0644)

	_, err := snapshot.Load(badFile)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestSave_TimestampIsRecent(t *testing.T) {
	dir := t.TempDir()
	before := time.Now().UTC().Add(-time.Second)
	_ = snapshot.Save(dir, "dev", makeEntries())
	after := time.Now().UTC().Add(time.Second)

	paths, _ := snapshot.List(dir)
	snap, _ := snapshot.Load(paths[0])

	if snap.Timestamp.Before(before) || snap.Timestamp.After(after) {
		t.Errorf("timestamp %v out of expected range", snap.Timestamp)
	}
}
