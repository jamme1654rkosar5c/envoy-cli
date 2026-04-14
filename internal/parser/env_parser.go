package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvEntry represents a single key-value pair from a .env file.
type EnvEntry struct {
	Key     string
	Value   string
	Comment string
	Line    int
}

// EnvFile represents a parsed .env file.
type EnvFile struct {
	Path    string
	Entries []EnvEntry
}

// ParseFile reads and parses a .env file from the given path.
func ParseFile(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file %q: %w", path, err)
	}
	defer f.Close()

	envFile := &EnvFile{Path: path}
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		entry, err := parseLine(trimmed, lineNum)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}
		envFile.Entries = append(envFile.Entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning env file: %w", err)
	}

	return envFile, nil
}

// parseLine parses a single KEY=VALUE line, optionally with an inline comment.
func parseLine(line string, lineNum int) (EnvEntry, error) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return EnvEntry{}, fmt.Errorf("invalid format, expected KEY=VALUE")
	}

	key := strings.TrimSpace(parts[0])
	if key == "" {
		return EnvEntry{}, fmt.Errorf("empty key")
	}

	raw := parts[1]
	value, comment := splitValueComment(raw)

	return EnvEntry{
		Key:     key,
		Value:   strings.Trim(strings.TrimSpace(value), `"`),
		Comment: comment,
		Line:    lineNum,
	}, nil
}

// splitValueComment separates inline comments from the value.
func splitValueComment(raw string) (value, comment string) {
	if idx := strings.Index(raw, " #"); idx != -1 {
		return raw[:idx], strings.TrimSpace(raw[idx+2:])
	}
	return raw, ""
}

// ToMap converts the EnvFile entries into a key-value map.
func (e *EnvFile) ToMap() map[string]string {
	m := make(map[string]string, len(e.Entries))
	for _, entry := range e.Entries {
		m[entry.Key] = entry.Value
	}
	return m
}
