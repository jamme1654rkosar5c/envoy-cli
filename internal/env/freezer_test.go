package env

import (
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeFreezEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "APP_ENV", Value: "production"},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
}

func TestFreeze_SetsComment(t *testing.T) {
	entries := makeFreezEntries()
	out, err := Freeze(entries, "APP_ENV", DefaultFreezeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !IsFrozen(out, "APP_ENV") {
		t.Error("expected APP_ENV to be frozen")
	}
}

func TestFreeze_MissingKey_Error(t *testing.T) {
	entries := makeFreezEntries()
	_, err := Freeze(entries, "MISSING", DefaultFreezeOptions())
	if err == nil {
		t.Error("expected error for missing key")
	}
}

func TestFreeze_MissingKey_NoError(t *testing.T) {
	entries := makeFreezEntries()
	opts := DefaultFreezeOptions()
	opts.AllowMissing = true
	_, err := Freeze(entries, "MISSING", opts)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUnfreeze_ClearsComment(t *testing.T) {
	entries := makeFreezEntries()
	out, _ := Freeze(entries, "APP_NAME", DefaultFreezeOptions())
	out, err := Unfreeze(out, "APP_NAME", DefaultFreezeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if IsFrozen(out, "APP_NAME") {
		t.Error("expected APP_NAME to be unfrozen")
	}
}

func TestFreeze_DoesNotMutateOriginal(t *testing.T) {
	entries := makeFreezEntries()
	_, _ = Freeze(entries, "APP_NAME", DefaultFreezeOptions())
	if IsFrozen(entries, "APP_NAME") {
		t.Error("original entries should not be mutated")
	}
}

func TestIsFrozen_ReturnsFalseForUnknownKey(t *testing.T) {
	entries := makeFreezEntries()
	if IsFrozen(entries, "NONEXISTENT") {
		t.Error("expected false for nonexistent key")
	}
}
