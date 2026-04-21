package env

import (
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeExpandEntries(pairs ...string) []parser.Entry {
	var entries []parser.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, parser.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return entries
}

func TestExpand_NothingToExpand(t *testing.T) {
	entries := makeExpandEntries("HOST", "localhost", "PORT", "5432")
	out, err := Expand(entries, DefaultExpandOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "localhost" || out[1].Value != "5432" {
		t.Errorf("values should be unchanged")
	}
}

func TestExpand_BraceStyle(t *testing.T) {
	entries := makeExpandEntries(
		"HOST", "db.local",
		"DSN", "postgres://${HOST}/mydb",
	)
	out, err := Expand(entries, DefaultExpandOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "postgres://db.local/mydb"
	if out[1].Value != want {
		t.Errorf("got %q, want %q", out[1].Value, want)
	}
}

func TestExpand_BareStyle(t *testing.T) {
	entries := makeExpandEntries(
		"USER", "admin",
		"GREETING", "hello $USER",
	)
	out, err := Expand(entries, DefaultExpandOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[1].Value != "hello admin" {
		t.Errorf("got %q", out[1].Value)
	}
}

func TestExpand_UnresolvedError(t *testing.T) {
	entries := makeExpandEntries("DSN", "postgres://${HOST}/mydb")
	_, err := Expand(entries, DefaultExpandOptions())
	if err == nil {
		t.Fatal("expected error for unresolved variable")
	}
}

func TestExpand_AllowMissing_NoError(t *testing.T) {
	opts := DefaultExpandOptions()
	opts.AllowMissing = true
	entries := makeExpandEntries("DSN", "postgres://${HOST}/mydb")
	out, err := Expand(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// reference left in place
	if out[0].Value != "postgres://${HOST}/mydb" {
		t.Errorf("got %q", out[0].Value)
	}
}

func TestExpand_EscapedDollar(t *testing.T) {
	entries := makeExpandEntries("PRICE", "$$10.00")
	out, err := Expand(entries, DefaultExpandOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "$10.00" {
		t.Errorf("got %q", out[0].Value)
	}
}

func TestExpand_MaxDepthExceeded(t *testing.T) {
	// A refers to B which refers back — depth check fires before true cycle.
	opts := DefaultExpandOptions()
	opts.MaxDepth = 2
	// Build a deep chain: A -> B -> C -> D (3 levels deep)
	entries := makeExpandEntries(
		"A", "1",
		"B", "${A}",
		"C", "${B}",
		"D", "${C}",
	)
	// With MaxDepth=2 the chain A->B->C->D should still resolve (depth never
	// exceeds 2 in this linear case). Use a circular-ish value to force it.
	opts.MaxDepth = 0
	_, err := Expand(entries, opts)
	if err == nil {
		t.Fatal("expected depth error")
	}
}
