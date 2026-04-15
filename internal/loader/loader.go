package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/envoy-cli/internal/parser"
)

// EnvFile represents a named, parsed environment file.
type EnvFile struct {
	Name string
	Path string
	File *parser.EnvFile
}

// LoadEnv loads and parses a single .env file from the given path.
func LoadEnv(path string) (*EnvFile, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	envFile, err := parser.ParseFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	return &EnvFile{
		Name: filepath.Base(path),
		Path: path,
		File: envFile,
	}, nil
}

// LoadAll loads and parses multiple .env files by path.
// It returns all successfully loaded files and a combined error for any failures.
func LoadAll(paths []string) ([]*EnvFile, error) {
	var results []*EnvFile
	var errs []error

	for _, p := range paths {
		ef, err := LoadEnv(p)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		results = append(results, ef)
	}

	if len(errs) > 0 {
		return results, joinErrors(errs)
	}
	return results, nil
}

// joinErrors combines multiple errors into a single error message.
func joinErrors(errs []error) error {
	msg := ""
	for i, e := range errs {
		if i > 0 {
			msg += "; "
		}
		msg += e.Error()
	}
	return fmt.Errorf("%s", msg)
}
