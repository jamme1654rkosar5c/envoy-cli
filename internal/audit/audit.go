// Package audit provides functionality for tracking and recording
// changes to .env files over time, enabling audit trail support.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/envoy-cli/internal/parser"
)

// EventType represents the kind of audit event.
type EventType string

const (
	EventLoad    EventType = "load"
	EventValidate EventType = "validate"
	EventMerge   EventType = "merge"
	EventExport  EventType = "export"
	EventDiff    EventType = "diff"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time  `json:"timestamp"`
	Event     EventType  `json:"event"`
	File      string     `json:"file"`
	Keys      []string   `json:"keys,omitempty"`
	Message   string     `json:"message,omitempty"`
}

// Log holds a sequence of audit entries.
type Log struct {
	Entries []Entry `json:"entries"`
}

// Record appends a new entry to the audit log file at logPath.
func Record(logPath string, event EventType, file string, env *parser.EnvFile, message string) error {
	log, err := loadLog(logPath)
	if err != nil {
		log = &Log{}
	}

	entry := Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		File:      file,
		Message:   message,
	}

	if env != nil {
		for _, e := range env.Entries {
			entry.Keys = append(entry.Keys, e.Key)
		}
	}

	log.Entries = append(log.Entries, entry)
	return saveLog(logPath, log)
}

// ReadLog loads and returns the audit log from logPath.
func ReadLog(logPath string) (*Log, error) {
	return loadLog(logPath)
}

func loadLog(logPath string) (*Log, error) {
	data, err := os.ReadFile(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Log{}, nil
		}
		return nil, fmt.Errorf("audit: read log: %w", err)
	}
	var log Log
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, fmt.Errorf("audit: parse log: %w", err)
	}
	return &log, nil
}

func saveLog(logPath string, log *Log) error {
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return fmt.Errorf("audit: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return fmt.Errorf("audit: marshal log: %w", err)
	}
	return os.WriteFile(logPath, data, 0o644)
}
