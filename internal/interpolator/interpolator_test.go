package interpolator_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/interpolator"
	"github.com/user/envoy-cli/internal/parser"
)

func makeEnvFile(entries []parser.Entry) parser.EnvFile {
	return parser.EnvFile{Path: "test.env", Entries: entries}
}

func TestInterpolate_NoBraces(t *testing.T) {
	file := makeEnvFile([]parser.Entry{
		{Key: "BASE", Value: "hello"},
		{Key: "GREETING", Value: "$BASE world"},
	})
	out, err := interpolator.Interpolate(file, interpolator.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out.Entries[1].Value; got != "hello world" {
		t.Errorf("expected %q, got %q", "hello world", got)
	}
}

func TestInterpolate_WithBraces(t *testing.T) {
	file := makeEnvFile([]parser.Entry{
		{Key: "HOST", Value: "localhost"},
		{Key: "DSN", Value: "postgres://${HOST}/db"},
	})
	out, err := interpolator.Interpolate(file, interpolator.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out.Entries[1].Value; got != "postgres://localhost/db" {
		t.Errorf("expected %q, got %q", "postgres://localhost/db", got)
	}
}

func TestInterpolate_UnresolvedError(t *testing.T) {
	file := makeEnvFile([]parser.Entry{
		{Key: "URL", Value: "http://${MISSING_HOST}/path"},
	})
	_, err := interpolator.Interpolate(file, interpolator.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for unresolved variable, got nil")
	}
}

func TestInterpolate_AllowMissing(t *testing.T) {
	file := makeEnvFile([]parser.Entry{
		{Key: "URL", Value: "http://${MISSING_HOST}/path"},
	})
	opts := interpolator.Options{AllowMissing: true}
	out, err := interpolator.Interpolate(file, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out.Entries[0].Value; got != "http://${MISSING_HOST}/path" {
		t.Errorf("expected original value preserved, got %q", got)
	}
}

func TestInterpolate_DoesNotMutateOriginal(t *testing.T) {
	original := makeEnvFile([]parser.Entry{
		{Key: "A", Value: "foo"},
		{Key: "B", Value: "$A bar"},
	})
	_, _ = interpolator.Interpolate(original, interpolator.DefaultOptions())
	if original.Entries[1].Value != "$A bar" {
		t.Error("original EnvFile was mutated")
	}
}

func TestInterpolate_PreservesComment(t *testing.T) {
	file := makeEnvFile([]parser.Entry{
		{Key: "X", Value: "val"},
		{Key: "Y", Value: "$X", Comment: "uses X"},
	})
	out, err := interpolator.Interpolate(file, interpolator.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[1].Comment != "uses X" {
		t.Errorf("comment was lost, got %q", out.Entries[1].Comment)
	}
}
