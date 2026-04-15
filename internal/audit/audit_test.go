package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/envoy-cli/internal/audit"
	"github.com/envoy-cli/internal/parser"
)

func makeEnvFile(entries map[string]string) *parser.EnvFile {
	f := &parser.EnvFile{}
	for k, v := range entries {
		f.Entries = append(f.Entries, parser.Entry{Key: k, Value: v})
	}
	return f
}

func TestRecord_CreatesLogFile(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.json")

	env := makeEnvFile(map[string]string{"APP_ENV": "production"})
	err := audit.Record(logPath, audit.EventLoad, ".env", env, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Fatal("expected log file to be created")
	}
}

func TestRecord_AppendsEntries(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.json")

	env := makeEnvFile(map[string]string{"KEY": "val"})
	_ = audit.Record(logPath, audit.EventLoad, ".env", env, "first")
	_ = audit.Record(logPath, audit.EventValidate, ".env", env, "second")

	log, err := audit.ReadLog(logPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(log.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(log.Entries))
	}
	if log.Entries[0].Event != audit.EventLoad {
		t.Errorf("expected first event to be %q", audit.EventLoad)
	}
	if log.Entries[1].Event != audit.EventValidate {
		t.Errorf("expected second event to be %q", audit.EventValidate)
	}
}

func TestRecord_StoresKeys(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.json")

	env := makeEnvFile(map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"})
	_ = audit.Record(logPath, audit.EventExport, "prod.env", env, "")

	log, _ := audit.ReadLog(logPath)
	if len(log.Entries[0].Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(log.Entries[0].Keys))
	}
}

func TestRecord_TimestampIsRecent(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.json")
	before := time.Now().UTC().Add(-time.Second)

	_ = audit.Record(logPath, audit.EventDiff, ".env", nil, "diff run")

	log, _ := audit.ReadLog(logPath)
	if log.Entries[0].Timestamp.Before(before) {
		t.Error("expected timestamp to be recent")
	}
}

func TestReadLog_NonExistentFile_ReturnsEmpty(t *testing.T) {
	log, err := audit.ReadLog("/nonexistent/path/audit.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(log.Entries) != 0 {
		t.Errorf("expected empty log, got %d entries", len(log.Entries))
	}
}
