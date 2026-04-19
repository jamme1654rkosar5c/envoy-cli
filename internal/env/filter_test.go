package env

import (
	"testing"

	"github.com/envoy-cli/internal/parser"
)

func makeFilterEntries() []parser.Entry {
	return []parser.Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PORT", Value: ""},
		{Key: "APP_NAME", Value: "envoy"},
		{Key: "APP_DEBUG", Value: "true"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestFilter_ByPrefix(t *testing.T) {
	result := Filter(makeFilterEntries(), FilterOptions{Prefix: "DB_"})
	if len(result) != 2 {
		t.Errorf("expected 2, got %d", len(result))
	}
}

func TestFilter_ByContains(t *testing.T) {
	result := Filter(makeFilterEntries(), FilterOptions{Contains: "APP"})
	if len(result) != 2 {
		t.Errorf("expected 2, got %d", len(result))
	}
}

func TestFilter_NoEmpty(t *testing.T) {
	result := Filter(makeFilterEntries(), FilterOptions{NoEmpty: true})
	for _, e := range result {
		if e.Value == "" {
			t.Errorf("expected no empty values, found key %s", e.Key)
		}
	}
	if len(result) != 4 {
		t.Errorf("expected 4 non-empty entries, got %d", len(result))
	}
}

func TestFilter_Combined(t *testing.T) {
	result := Filter(makeFilterEntries(), FilterOptions{Prefix: "APP_", NoEmpty: true})
	if len(result) != 2 {
		t.Errorf("expected 2, got %d", len(result))
	}
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	result := Filter(makeFilterEntries(), FilterOptions{})
	if len(result) != 5 {
		t.Errorf("expected 5, got %d", len(result))
	}
}

func TestFilter_EmptyInput(t *testing.T) {
	result := Filter([]parser.Entry{}, FilterOptions{Prefix: "DB_", NoEmpty: true})
	if len(result) != 0 {
		t.Errorf("expected 0, got %d", len(result))
	}
}
