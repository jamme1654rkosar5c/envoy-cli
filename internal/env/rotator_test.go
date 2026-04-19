package env

import (
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeRotateEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASS", Value: "secret"},
		{Key: "APP_ENV", Value: "production"},
	}
}

func TestRotate_RenamesKey(t *testing.T) {
	entries := makeRotateEntries()
	rotations := []RotateEntry{{OldKey: "DB_PASS", NewKey: "DB_PASSWORD", NewValue: ""}}
	out, err := Rotate(entries, rotations, DefaultRotateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Key != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD, got %s", out[1].Key)
	}
	if out[1].Value != "secret" {
		t.Errorf("expected original value preserved, got %s", out[1].Value)
	}
}

func TestRotate_ReplacesValue(t *testing.T) {
	entries := makeRotateEntries()
	rotations := []RotateEntry{{OldKey: "DB_PASS", NewKey: "DB_PASSWORD", NewValue: "newsecret"}}
	out, err := Rotate(entries, rotations, DefaultRotateOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "newsecret" {
		t.Errorf("expected newsecret, got %s", out[1].Value)
	}
}

func TestRotate_MissingKey_Error(t *testing.T) {
	entries := makeRotateEntries()
	rotations := []RotateEntry{{OldKey: "MISSING", NewKey: "ALSO_MISSING", NewValue: ""}}
	_, err := Rotate(entries, rotations, DefaultRotateOptions())
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRotate_MissingKey_NoError(t *testing.T) {
	entries := makeRotateEntries()
	rotations := []RotateEntry{{OldKey: "MISSING", NewKey: "ALSO_MISSING", NewValue: ""}}
	opts := DefaultRotateOptions()
	opts.ErrorOnMissing = false
	out, err := Rotate(entries, rotations, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries, got %d", len(out))
	}
}

func TestRotate_DryRun_DoesNotMutate(t *testing.T) {
	entries := makeRotateEntries()
	rotations := []RotateEntry{{OldKey: "DB_PASS", NewKey: "DB_PASSWORD", NewValue: "new"}}
	opts := DefaultRotateOptions()
	opts.DryRun = true
	_, err := Rotate(entries, rotations, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries[1].Key != "DB_PASS" {
		t.Errorf("original entries mutated during dry run")
	}
}
