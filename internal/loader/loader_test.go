package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func TestLoadEnv_ValidFile(t *testing.T) {
	path := writeTempEnv(t, "APP_ENV=production\nDEBUG=false\n")

	ef, err := LoadEnv(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if ef.File == nil {
		t.Fatal("expected parsed EnvFile, got nil")
	}
	if ef.Name != ".env" {
		t.Errorf("expected name '.env', got %q", ef.Name)
	}
}

func TestLoadEnv_FileNotFound(t *testing.T) {
	_, err := LoadEnv("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadAll_AllValid(t *testing.T) {
	p1 := writeTempEnv(t, "KEY1=value1\n")
	p2 := writeTempEnv(t, "KEY2=value2\n")

	results, err := LoadAll([]string{p1, p2})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestLoadAll_PartialFailure(t *testing.T) {
	valid := writeTempEnv(t, "KEY=value\n")

	results, err := LoadAll([]string{valid, "/bad/path/.env"})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
	if len(results) != 1 {
		t.Errorf("expected 1 valid result, got %d", len(results))
	}
}

func TestLoadAll_Empty(t *testing.T) {
	results, err := LoadAll([]string{})
	if err != nil {
		t.Fatalf("expected no error for empty input, got: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
