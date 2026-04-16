package env

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makeTransformEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "APP_ENV", Value: "  production  "},
		{Key: "SECRET_KEY", Value: "abc123"},
		{Key: "DB_HOST", Value: "localhost"},
	}
}

func TestTransform_UppercaseAll(t *testing.T) {
	entries := makeTransformEntries()
	out := Transform(entries, BuiltinUppercase, DefaultTransformOptions())
	if out[0].Value != "MYAPP" {
		t.Errorf("expected MYAPP, got %s", out[0].Value)
	}
	if out[2].Value != "ABC123" {
		t.Errorf("expected ABC123, got %s", out[2].Value)
	}
}

func TestTransform_OnlyKeys(t *testing.T) {
	entries := makeTransformEntries()
	opts := DefaultTransformOptions()
	opts.OnlyKeys = []string{"APP_NAME"}
	out := Transform(entries, BuiltinUppercase, opts)
	if out[0].Value != "MYAPP" {
		t.Errorf("expected MYAPP, got %s", out[0].Value)
	}
	// APP_ENV should be unchanged
	if out[1].Value != "  production  " {
		t.Errorf("expected unchanged value, got %s", out[1].Value)
	}
}

func TestTransform_SkipPrefixes(t *testing.T) {
	entries := makeTransformEntries()
	opts := DefaultTransformOptions()
	opts.SkipPrefixes = []string{"SECRET_"}
	out := Transform(entries, BuiltinUppercase, opts)
	// SECRET_KEY should be unchanged
	if out[2].Value != "abc123" {
		t.Errorf("expected abc123 unchanged, got %s", out[2].Value)
	}
	if out[0].Value != "MYAPP" {
		t.Errorf("expected MYAPP, got %s", out[0].Value)
	}
}

func TestTransform_TrimSpace(t *testing.T) {
	entries := makeTransformEntries()
	out := Transform(entries, BuiltinTrimSpace, DefaultTransformOptions())
	if out[1].Value != "production" {
		t.Errorf("expected 'production', got %q", out[1].Value)
	}
}

func TestTransform_PreservesComment(t *testing.T) {
	entries := []parser.Entry{
		{Key: "FOO", Value: "bar", Comment: "important"},
	}
	out := Transform(entries, BuiltinUppercase, DefaultTransformOptions())
	if out[0].Comment != "important" {
		t.Errorf("expected comment preserved, got %s", out[0].Comment)
	}
}
