package exporter_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/internal/exporter"
	"github.com/envoy-cli/internal/parser"
)

func makeEnvFile(entries []parser.Entry) *parser.EnvFile {
	return &parser.EnvFile{Entries: entries}
}

func TestExport_DotEnvFormat(t *testing.T) {
	env := makeEnvFile([]parser.Entry{
		{Key: "APP_ENV", Value: "production"},
		{Key: "PORT", Value: "8080"},
	})
	dest := filepath.Join(t.TempDir(), "out.env")
	if err := exporter.Export(env, dest, exporter.FormatDotEnv); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(dest)
	got := string(data)
	if got != "APP_ENV=production\nPORT=8080\n" {
		t.Errorf("unexpected dotenv output:\n%s", got)
	}
}

func TestExport_JSONFormat(t *testing.T) {
	env := makeEnvFile([]parser.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: "5432"},
	})
	dest := filepath.Join(t.TempDir(), "out.json")
	if err := exporter.Export(env, dest, exporter.FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(dest)
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("invalid json output: %v", err)
	}
	if m["DB_HOST"] != "localhost" || m["DB_PORT"] != "5432" {
		t.Errorf("unexpected json values: %v", m)
	}
}

func TestExport_ExportScriptFormat(t *testing.T) {
	env := makeEnvFile([]parser.Entry{
		{Key: "SECRET", Value: "abc123"},
	})
	dest := filepath.Join(t.TempDir(), "out.sh")
	if err := exporter.Export(env, dest, exporter.FormatExport); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(dest)
	got := string(data)
	if got != "export SECRET=\"abc123\"\n" {
		t.Errorf("unexpected export script output:\n%s", got)
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	env := makeEnvFile([]parser.Entry{})
	dest := filepath.Join(t.TempDir(), "out.txt")
	err := exporter.Export(env, dest, exporter.Format("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}
