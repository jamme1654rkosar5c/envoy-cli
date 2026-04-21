package env

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeMaskEntries(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestMask_SensitiveKeyIsmasked(t *testing.T) {
	entries := makeMaskEntries("DB_PASSWORD", "secret123", "APP_NAME", "myapp")
	opts := DefaultMaskOptions()
	out := Mask(entries, opts)

	if out[0].Value == "secret123" {
		t.Error("expected DB_PASSWORD value to be masked")
	}
	if out[1].Value != "myapp" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", out[1].Value)
	}
}

func TestMask_DoesNotMutateOriginal(t *testing.T) {
	entries := makeMaskEntries("API_TOKEN", "tok_abc123")
	opts := DefaultMaskOptions()
	_ = Mask(entries, opts)

	if entries[0].Value != "tok_abc123" {
		t.Error("original entry was mutated")
	}
}

func TestMask_VisibleSuffix(t *testing.T) {
	entries := makeMaskEntries("SECRET_KEY", "abcdef")
	opts := DefaultMaskOptions()
	opts.VisibleSuffix = 2
	out := Mask(entries, opts)

	if !strings.HasSuffix(out[0].Value, "ef") {
		t.Errorf("expected suffix 'ef', got %q", out[0].Value)
	}
	if strings.Contains(out[0].Value[:len(out[0].Value)-2], "a") {
		t.Error("expected leading chars to be masked")
	}
}

func TestMask_UsePlaceholder(t *testing.T) {
	entries := makeMaskEntries("DB_PASS", "hunter2")
	opts := DefaultMaskOptions()
	opts.UsePlaceholder = true
	opts.Placeholder = "[REDACTED]"
	out := Mask(entries, opts)

	if out[0].Value != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", out[0].Value)
	}
}

func TestMask_MaskAllValues(t *testing.T) {
	entries := makeMaskEntries("SAFE_KEY", "visible", "OTHER", "data")
	opts := DefaultMaskOptions()
	opts.MaskAllValues = true
	out := Mask(entries, opts)

	for _, e := range out {
		if e.Value == entries[0].Value || e.Value == entries[1].Value {
			t.Errorf("expected value to be masked for key %s", e.Key)
		}
	}
}

func TestMask_EmptyValueSkipped(t *testing.T) {
	entries := makeMaskEntries("SECRET", "")
	opts := DefaultMaskOptions()
	out := Mask(entries, opts)

	if out[0].Value != "" {
		t.Errorf("expected empty value to stay empty, got %q", out[0].Value)
	}
}
