package env

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/parser"
)

func TestBuildMaskSummaries_DetectsMasked(t *testing.T) {
	orig := []parser.Entry{{Key: "API_SECRET", Value: "abc123"}, {Key: "APP_NAME", Value: "myapp"}}
	masked := []parser.Entry{{Key: "API_SECRET", Value: "******"}, {Key: "APP_NAME", Value: "myapp"}}

	summaries := BuildMaskSummaries(orig, masked)

	if len(summaries) != 2 {
		t.Fatalf("expected 2 summaries, got %d", len(summaries))
	}
	if !summaries[0].WasMasked {
		t.Error("expected API_SECRET to be marked as masked")
	}
	if summaries[1].WasMasked {
		t.Error("expected APP_NAME to not be marked as masked")
	}
}

func TestBuildMaskSummaries_OrigLen(t *testing.T) {
	orig := []parser.Entry{{Key: "DB_PASSWORD", Value: "hunter2"}}
	masked := []parser.Entry{{Key: "DB_PASSWORD", Value: "*******"}}

	summaries := BuildMaskSummaries(orig, masked)

	if summaries[0].OrigLen != 7 {
		t.Errorf("expected OrigLen 7, got %d", summaries[0].OrigLen)
	}
}

func TestFormatMask_ContainsHeaders(t *testing.T) {
	summaries := []MaskSummary{
		{Key: "SECRET", OrigLen: 6, Masked: "******", WasMasked: true},
	}
	out := FormatMask(summaries)

	for _, header := range []string{"KEY", "MASKED", "ORIG LEN", "DISPLAY"} {
		if !strings.Contains(out, header) {
			t.Errorf("expected header %q in output", header)
		}
	}
}

func TestFormatMask_EmptySummaries(t *testing.T) {
	out := FormatMask(nil)
	if !strings.Contains(out, "No entries") {
		t.Errorf("expected empty message, got %q", out)
	}
}

func TestFormatMask_ShowsYesForMasked(t *testing.T) {
	summaries := []MaskSummary{
		{Key: "TOKEN", OrigLen: 10, Masked: "**********", WasMasked: true},
	}
	out := FormatMask(summaries)
	if !strings.Contains(out, "yes") {
		t.Errorf("expected 'yes' in output for masked entry, got:\n%s", out)
	}
}
