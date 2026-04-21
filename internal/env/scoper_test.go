package env

import (
	"testing"
)

func makeScopeEntries() []Entry {
	return []Entry{
		{Key: "APP_HOST", Value: "localhost"},
		{Key: "APP_PORT", Value: "8080"},
		{Key: "DB_HOST", Value: "db.local"},
		{Key: "DB_PASS", Value: "secret"},
		{Key: "DEBUG", Value: "true"},
	}
}

func TestScope_FiltersByPrefix(t *testing.T) {
	entries := makeScopeEntries()
	opts := DefaultScopeOptions()
	opts.Scope = "APP_"

	result := Scope(entries, opts)

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Key != "APP_HOST" || result[1].Key != "APP_PORT" {
		t.Errorf("unexpected keys: %v, %v", result[0].Key, result[1].Key)
	}
}

func TestScope_StripPrefix_RemovesScopeFromKey(t *testing.T) {
	entries := makeScopeEntries()
	opts := DefaultScopeOptions()
	opts.Scope = "APP_"
	opts.StripPrefix = true

	result := Scope(entries, opts)

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Key != "HOST" {
		t.Errorf("expected HOST, got %s", result[0].Key)
	}
	if result[1].Key != "PORT" {
		t.Errorf("expected PORT, got %s", result[1].Key)
	}
}

func TestScope_AddPrefix_PrependsToAllKeys(t *testing.T) {
	entries := []Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "PORT", Value: "9090"},
	}
	opts := DefaultScopeOptions()
	opts.Scope = "SVC_"
	opts.AddPrefix = true

	result := Scope(entries, opts)

	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Key != "SVC_HOST" {
		t.Errorf("expected SVC_HOST, got %s", result[0].Key)
	}
	if result[1].Key != "SVC_PORT" {
		t.Errorf("expected SVC_PORT, got %s", result[1].Key)
	}
}

func TestScope_EmptyScope_ReturnsAll(t *testing.T) {
	entries := makeScopeEntries()
	opts := DefaultScopeOptions()

	result := Scope(entries, opts)

	if len(result) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(result))
	}
}

func TestScope_DoesNotMutateOriginal(t *testing.T) {
	entries := makeScopeEntries()
	opts := DefaultScopeOptions()
	opts.Scope = "APP_"
	opts.StripPrefix = true

	_ = Scope(entries, opts)

	if entries[0].Key != "APP_HOST" {
		t.Errorf("original entry mutated, expected APP_HOST got %s", entries[0].Key)
	}
}
