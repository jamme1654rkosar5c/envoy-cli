package env

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makePatchEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "APP_ENV", Value: "development"},
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestPatch_UpdateExistingKey(t *testing.T) {
	entries := makePatchEntries()
	ops := []PatchOp{{Key: "PORT", Value: "9090"}}
	out, err := Patch(entries, ops, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[2].Value != "9090" {
		t.Errorf("expected 9090, got %s", out[2].Value)
	}
}

func TestPatch_InsertNewKey(t *testing.T) {
	entries := makePatchEntries()
	ops := []PatchOp{{Key: "LOG_LEVEL", Value: "debug"}}
	out, err := Patch(entries, ops, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(out))
	}
	if out[3].Key != "LOG_LEVEL" || out[3].Value != "debug" {
		t.Errorf("unexpected entry: %+v", out[3])
	}
}

func TestPatch_DeleteKey(t *testing.T) {
	entries := makePatchEntries()
	ops := []PatchOp{{Key: "DB_HOST", Delete: true}}
	out, err := Patch(entries, ops, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	for _, e := range out {
		if e.Key == "DB_HOST" {
			t.Error("DB_HOST should have been deleted")
		}
	}
}

func TestPatch_ErrorOnMissing(t *testing.T) {
	entries := makePatchEntries()
	opts := DefaultPatchOptions()
	opts.ErrorOnMissing = true
	ops := []PatchOp{{Key: "MISSING_KEY", Value: "val"}}
	_, err := Patch(entries, ops, opts)
	if err == nil {
		t.Error("expected error for missing key, got nil")
	}
}

func TestPatch_DoesNotMutateOriginal(t *testing.T) {
	entries := makePatchEntries()
	ops := []PatchOp{{Key: "PORT", Value: "1111"}}
	_, err := Patch(entries, ops, DefaultPatchOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[2].Value != "8080" {
		t.Error("original entries were mutated")
	}
}
