package env

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func makeDigestEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "DEBUG", Value: "false"},
		{Key: "SECRET_KEY", Value: "abc123"},
	}
}

func TestDigest_SHA256_ReturnsSameHashForSameInput(t *testing.T) {
	entries := makeDigestEntries()
	opts := DefaultDigestOptions()

	h1, err := Digest(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	h2, err := Digest(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h1 != h2 {
		t.Errorf("expected same hash, got %s vs %s", h1, h2)
	}
}

func TestDigest_MD5_ProducesHash(t *testing.T) {
	entries := makeDigestEntries()
	opts := DefaultDigestOptions()
	opts.Algorithm = DigestMD5

	h, err := Digest(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(h) != 32 {
		t.Errorf("expected 32-char md5 hex, got %d chars", len(h))
	}
}

func TestDigest_UnsupportedAlgorithm_ReturnsError(t *testing.T) {
	entries := makeDigestEntries()
	opts := DefaultDigestOptions()
	opts.Algorithm = DigestAlgorithm("blake2")

	_, err := Digest(entries, opts)
	if err == nil {
		t.Fatal("expected error for unsupported algorithm")
	}
}

func TestDigest_DifferentValues_ProducesDifferentHash(t *testing.T) {
	a := []parser.Entry{{Key: "FOO", Value: "bar"}}
	b := []parser.Entry{{Key: "FOO", Value: "baz"}}
	opts := DefaultDigestOptions()

	h1, _ := Digest(a, opts)
	h2, _ := Digest(b, opts)
	if h1 == h2 {
		t.Error("expected different hashes for different values")
	}
}

func TestDigest_ExcludeKeys_OmitsFromHash(t *testing.T) {
	entries := makeDigestEntries()
	opts := DefaultDigestOptions()

	hFull, _ := Digest(entries, opts)

	opts.ExcludeKeys = []string{"SECRET_KEY"}
	hPartial, _ := Digest(entries, opts)

	if hFull == hPartial {
		t.Error("expected different hashes when excluding a key")
	}
}

func TestDigest_IncludeKeys_LimitsScope(t *testing.T) {
	entries := makeDigestEntries()
	opts := DefaultDigestOptions()
	opts.IncludeKeys = []string{"APP_NAME"}

	h, err := Digest(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(h) == "" {
		t.Error("expected non-empty hash")
	}
}

func TestDigest_SortKeys_IsDeterministic(t *testing.T) {
	a := []parser.Entry{
		{Key: "Z_KEY", Value: "1"},
		{Key: "A_KEY", Value: "2"},
	}
	b := []parser.Entry{
		{Key: "A_KEY", Value: "2"},
		{Key: "Z_KEY", Value: "1"},
	}
	opts := DefaultDigestOptions()

	h1, _ := Digest(a, opts)
	h2, _ := Digest(b, opts)
	if h1 != h2 {
		t.Error("expected same hash regardless of input order when SortKeys=true")
	}
}
