package rename_test

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
	"github.com/envoy-cli/internal/rename"
)

func makeEnvFile(path string, pairs [][2]string) *parser.EnvFile {
	entries := make([]parser.Entry, 0, len(pairs))
	for _, p := range pairs {
		entries = append(entries, parser.Entry{Key: p[0], Value: p[1]})
	}
	return &parser.EnvFile{Path: path, Entries: entries}
}

func TestRenameKey_Success(t *testing.T) {
	f := makeEnvFile(".env", [][2]string{{"OLD_KEY", "value"}})
	r, err := rename.RenameKey(f, "OLD_KEY", "NEW_KEY", rename.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Renamed {
		t.Error("expected Renamed=true")
	}
	if f.Entries[0].Key != "NEW_KEY" {
		t.Errorf("expected NEW_KEY, got %s", f.Entries[0].Key)
	}
}

func TestRenameKey_NotFound_Error(t *testing.T) {
	f := makeEnvFile(".env", [][2]string{{"FOO", "bar"}})
	_, err := rename.RenameKey(f, "MISSING", "NEW_KEY", rename.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRenameKey_NotFound_NoError(t *testing.T) {
	f := makeEnvFile(".env", [][2]string{{"FOO", "bar"}})
	opts := rename.DefaultOptions()
	opts.ErrorIfNotFound = false
	r, err := rename.RenameKey(f, "MISSING", "NEW_KEY", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Renamed {
		t.Error("expected Renamed=false")
	}
}

func TestRenameKey_DryRun_DoesNotMutate(t *testing.T) {
	f := makeEnvFile(".env", [][2]string{{"OLD_KEY", "val"}})
	opts := rename.DefaultOptions()
	opts.DryRun = true
	r, err := rename.RenameKey(f, "OLD_KEY", "NEW_KEY", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Renamed {
		t.Error("expected Renamed=true even in dry-run")
	}
	if f.Entries[0].Key != "OLD_KEY" {
		t.Errorf("dry-run must not mutate: got %s", f.Entries[0].Key)
	}
}

func TestRenameKey_DuplicateTarget_Error(t *testing.T) {
	f := makeEnvFile(".env", [][2]string{{"OLD_KEY", "v1"}, {"NEW_KEY", "v2"}})
	_, err := rename.RenameKey(f, "OLD_KEY", "NEW_KEY", rename.DefaultOptions())
	if err == nil {
		t.Fatal("expected error when target key already exists")
	}
}

func TestRenameKeyInAll_AllRenamed(t *testing.T) {
	f1 := makeEnvFile(".env.dev", [][2]string{{"OLD", "a"}})
	f2 := makeEnvFile(".env.prod", [][2]string{{"OLD", "b"}})
	results, err := rename.RenameKeyInAll([]*parser.EnvFile{f1, f2}, "OLD", "RENAMED", rename.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Renamed {
			t.Errorf("expected Renamed=true for %s", r.File)
		}
	}
}
