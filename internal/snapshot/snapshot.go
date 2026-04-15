package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/envoy-cli/internal/parser"
)

// Snapshot represents a saved state of an env file at a point in time.
type Snapshot struct {
	Timestamp time.Time          `json:"timestamp"`
	Label     string             `json:"label"`
	Entries   []parser.EnvEntry  `json:"entries"`
}

// Save writes a snapshot of the given env file to the snapshot directory.
func Save(dir, label string, entries []parser.EnvEntry) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create snapshot directory: %w", err)
	}

	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Entries:   entries,
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	filename := fmt.Sprintf("%s_%s.json", label, snap.Timestamp.Format("20060102T150405Z"))
	dest := filepath.Join(dir, filename)

	if err := os.WriteFile(dest, data, 0644); err != nil {
		return fmt.Errorf("failed to write snapshot file: %w", err)
	}

	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot file: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("failed to parse snapshot file: %w", err)
	}

	return &snap, nil
}

// List returns all snapshot file paths found in the given directory.
func List(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read snapshot directory: %w", err)
	}

	var paths []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			paths = append(paths, filepath.Join(dir, e.Name()))
		}
	}
	return paths, nil
}
