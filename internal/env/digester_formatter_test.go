package env

import (
	"strings"
	"testing"
)

func makeDigestSummaries() []DigestSummary {
	return []DigestSummary{
		{Label: "production", Algorithm: DigestSHA256, Hash: strings.Repeat("a", 64), KeyCount: 12},
		{Label: "staging", Algorithm: DigestMD5, Hash: strings.Repeat("b", 32), KeyCount: 8},
	}
}

func TestFormatDigest_ContainsHeaders(t *testing.T) {
	summaries := makeDigestSummaries()
	out := FormatDigest(summaries)

	for _, hdr := range []string{"LABEL", "ALGO", "HASH", "KEYS"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("expected header %q in output", hdr)
		}
	}
}

func TestFormatDigest_ShowsLabel(t *testing.T) {
	summaries := makeDigestSummaries()
	out := FormatDigest(summaries)

	if !strings.Contains(out, "production") {
		t.Error("expected label 'production' in output")
	}
	if !strings.Contains(out, "staging") {
		t.Error("expected label 'staging' in output")
	}
}

func TestFormatDigest_ShowsAlgorithm(t *testing.T) {
	summaries := makeDigestSummaries()
	out := FormatDigest(summaries)

	if !strings.Contains(out, "sha256") {
		t.Error("expected algorithm 'sha256' in output")
	}
	if !strings.Contains(out, "md5") {
		t.Error("expected algorithm 'md5' in output")
	}
}

func TestFormatDigest_EmptySummaries(t *testing.T) {
	out := FormatDigest(nil)
	if !strings.Contains(out, "No digests") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestBuildDigestSummary_Fields(t *testing.T) {
	s := BuildDigestSummary("dev", "abc123", DigestMD5, 5)
	if s.Label != "dev" {
		t.Errorf("expected label 'dev', got %q", s.Label)
	}
	if s.Hash != "abc123" {
		t.Errorf("expected hash 'abc123', got %q", s.Hash)
	}
	if s.KeyCount != 5 {
		t.Errorf("expected key count 5, got %d", s.KeyCount)
	}
}
